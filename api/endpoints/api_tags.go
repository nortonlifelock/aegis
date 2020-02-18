package endpoints

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nortonlifelock/aegis/internal/implementations"
	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/domain"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

type tagNames struct {
	Names []string `json:"name"`
}

func getTagsFromDb(w http.ResponseWriter, r *http.Request) {

	executeTransaction(w, r, getTagsFromDBEndpoint, admin|manager, func(trans *transaction) {
		var params = mux.Vars(r)
		cloudService := params[sourceParam]
		if len(cloudService) > 0 {

			var cloudSource domain.Source
			cloudSource, trans.err = Ms.GetSourceByName(cloudService)
			if trans.err == nil {
				if cloudSource != nil {

					var tags []domain.TagMap
					tags, trans.err = Ms.GetTagMapsByOrgCloudSourceID(cloudSource.ID(), trans.permission.OrgID())
					if trans.err == nil {
						trans.obj, trans.status = toTagDtoSlice(tags), http.StatusOK
					} else {
						(&trans.wrapper).addError(trans.err, databaseError)
					}

				} else {
					trans.err = fmt.Errorf("could not find a source for [%s]", cloudService)
					(&trans.wrapper).addError(trans.err, requestFormatError)
				}

			} else {
				(&trans.wrapper).addError(trans.err, databaseError)
			}

		} else {
			trans.err = fmt.Errorf("cloud service not provided to getTagsFromDb")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}

	})
}

func getTagsFromAzure(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getAzureTagsEndpoint, admin|manager, func(trans *transaction) {
		var connection integrations.CloudServiceConnection
		connection, trans.err = getCloudConnection(integrations.Azure, trans.permission.OrgID())
		if trans.err == nil {

			var tags = &tagNames{}
			tags.Names, trans.err = connection.GetAllTagNames()
			if trans.err == nil {
				sort.Strings(tags.Names)
				trans.obj, trans.status = tags, http.StatusOK
			} else {
				(&trans.wrapper).addError(trans.err, processError)
			}
		} else {
			(&trans.wrapper).addError(trans.err, processError)
		}
	})
}

func getTagsFromAws(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getAWSTagsEndpoint, admin|manager, func(trans *transaction) {
		var connection integrations.CloudServiceConnection
		connection, trans.err = getCloudConnection(integrations.AWS, trans.permission.OrgID())
		if trans.err == nil {

			var tags = &tagNames{}
			tags.Names, trans.err = connection.GetAllTagNames()
			if trans.err == nil {
				sort.Strings(tags.Names)
				trans.obj, trans.status = tags, http.StatusOK
			} else {
				(&trans.wrapper).addError(trans.err, processError)
			}
		} else {
			(&trans.wrapper).addError(trans.err, processError)
		}
	})
}

// TODO TagKey field in db is a thing now
func createAwsTags(w http.ResponseWriter, r *http.Request) {

	executeTransaction(w, r, createAWSTagsEndpoint, admin|manager, func(trans *transaction) {
		if trans.endpoint.verify() {
			if tag, isTag := trans.endpoint.(*TagMap); isTag {
				trans.obj, trans.status, trans.err = tag.create(trans.user, trans.permission)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = fmt.Errorf("tried to pass a non-tag as a tag")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = fmt.Errorf("input did not pass validation")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func updateAwsTags(w http.ResponseWriter, r *http.Request) {

	executeTransaction(w, r, updateAWSTagsEndpoint, admin|manager, func(trans *transaction) {
		if trans.endpoint.verify() {
			if tag, isTag := trans.endpoint.(*TagMap); isTag {
				trans.obj, trans.status, trans.err = tag.update(trans.user, trans.permission, trans.originalBody)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = fmt.Errorf("tried to pass a non-tag as a tag")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = fmt.Errorf("input did not pass validation")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func deleteAwsTags(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, deleteAWSTagsEndpoint, admin, func(trans *transaction) {
		if trans.endpoint.verify() {
			if tag, isTag := trans.endpoint.(*TagMap); isTag {
				trans.obj, trans.status, trans.err = tag.delete(trans.user, trans.permission)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = fmt.Errorf("tried to pass a non-tag as a tag")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = fmt.Errorf("input did not pass validation")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func (tag *TagMap) create(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}
	status = http.StatusBadRequest

	var ticketSource, cloudSource domain.Source
	ticketSource, err = Ms.GetSourceByName(tag.TicketSource)
	if err == nil {
		cloudSource, err = Ms.GetSourceByName(tag.CloudSource)
		if err == nil {

			var existingMaps []domain.TagMap
			existingMaps, err = Ms.GetTagMapsByOrg(permission.OrgID())
			if err == nil {

				// TODO need to check if there is already a tag -> ask if they want to update if there is
				if ticketSource != nil && cloudSource != nil {

					var existingMapErrorMessage = ""

					var exactMapExists = false
					for _, existingMap := range existingMaps {
						var existMapCheck bool
						existMapCheck, existingMapErrorMessage, err = tag.findAndDeleteOverlappingTags(existingMap, ticketSource, cloudSource, user, permission)
						if err == nil && len(existingMapErrorMessage) == 0 {
							if existMapCheck {
								exactMapExists = true
							}
						} else {
							break
						}
					}

					if err == nil {
						if !exactMapExists {
							if len(existingMapErrorMessage) == 0 {
								_, _, err = Ms.CreateTagMap(ticketSource.ID(), tag.TicketTag, cloudSource.ID(), tag.CloudTag, tag.Option, permission.OrgID())
								if err == nil {
									status = http.StatusOK
									generalResp.Response = "tag created"
								} else {
									err = fmt.Errorf("error while creating tag - [%s]", err.Error())
								}
							} else {
								status = http.StatusOK
								type retry struct {
									Retry string `json:"retry"`
								}
								rty := &retry{}
								rty.Retry = existingMapErrorMessage
								generalResp.Response = rty
							}
						} else {
							// this tag mapping exists - no need to create it again
							// no message needs to be sent to the client either
							status = http.StatusOK
						}
					}

				} else {
					err = fmt.Errorf("either the ticketing or the cloud source could not be found in the database")
				}

			} else {
				err = fmt.Errorf("error while grabbing tag maps - %s", err.Error())
			}

		} else {
			err = fmt.Errorf("error while gathering cloud source - %s", err.Error())
		}
	} else {
		err = fmt.Errorf("error while gathering ticketing source - %s", err.Error())
	}

	return generalResp, status, err
}

func (tag *TagMap) findAndDeleteOverlappingTags(existingMap domain.TagMap, ticketSource domain.Source, cloudSource domain.Source, user domain.User, permission domain.Permission) (exactMapExists bool, existingMapErrorMessage string, err error) {
	if existingMap.TicketingSourceID() == ticketSource.ID() {
		if existingMap.CloudSourceID() == cloudSource.ID() {
			if strings.ToLower(existingMap.TicketingTag()) == strings.ToLower(tag.TicketTag) {

				if strings.ToLower(existingMap.CloudTag()) == strings.ToLower(tag.CloudTag) {
					exactMapExists = true
				} else if strings.ToLower(existingMap.Options()) == strings.ToLower(implementations.Overwrite) {
					// there is an existing tag that overwrites this field

					// delete the old tag
					if tag.Overwrite {
						deleteExistingTag := &TagMap{
							CloudSource:  cloudSource.Source(),
							CloudTag:     existingMap.CloudTag(),
							TicketSource: ticketSource.Source(),
							TicketTag:    existingMap.TicketingTag(),
							Option:       existingMap.Options(),
						}

						_, _, err = deleteExistingTag.delete(user, permission)
					} else {
						existingMapErrorMessage = fmt.Sprintf("There is an existing mapping [%s->%s] that overwrites the ticketing field. Do you want to replace this existing mapping?", existingMap.CloudTag(), existingMap.TicketingTag())
					}

				}

			}
		}
	}
	return exactMapExists, existingMapErrorMessage, err
}

func (tag *TagMap) update(user domain.User, permission domain.Permission, originalBody string) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}
	status = http.StatusBadRequest

	var ticketSource, cloudSource domain.Source
	ticketSource, err = Ms.GetSourceByName(tag.TicketSource)
	if err == nil {
		cloudSource, err = Ms.GetSourceByName(tag.CloudSource)
		if err == nil {

			if ticketSource != nil && cloudSource != nil {

				_, _, err = Ms.UpdateTagMap(ticketSource.ID(), tag.TicketTag, cloudSource.ID(), tag.CloudTag, tag.Option, permission.OrgID())
				if err == nil {
					status = http.StatusOK
				} else {
					err = fmt.Errorf("error while updating tag - [%s]", err.Error())
				}

			} else {
				err = fmt.Errorf("either the ticketing or the cloud source could not be found in the database")
			}

		} else {
			err = fmt.Errorf("error while gathering cloud source")
		}
	} else {
		err = fmt.Errorf("error while gathering ticketing source")
	}

	return generalResp, status, err
}

func (tag *TagMap) delete(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}
	status = http.StatusBadRequest

	var ticketSource, cloudSource domain.Source
	ticketSource, err = Ms.GetSourceByName(tag.TicketSource)
	if err == nil {
		cloudSource, err = Ms.GetSourceByName(tag.CloudSource)
		if err == nil {

			if ticketSource != nil && cloudSource != nil {

				_, _, err = Ms.DeleteTagMap(ticketSource.ID(), tag.TicketTag, cloudSource.ID(), tag.CloudTag, permission.OrgID())
				if err == nil {
					status = http.StatusOK
					generalResp.Response = "tag deleted"
				} else {
					err = fmt.Errorf("error while deleting tag - [%s]", err.Error())
				}

			} else {
				err = fmt.Errorf("either the ticketing or the cloud source could not be found in the database")
			}

		} else {
			err = fmt.Errorf("error while gathering cloud source")
		}
	} else {
		err = fmt.Errorf("error while gathering ticketing source")
	}

	return generalResp, status, err
}

func (tag *TagMap) verify() bool {
	valid := false

	if len(tag.CloudSource) > 0 && len(tag.CloudTag) > 0 && len(tag.TicketTag) > 0 && len(tag.TicketSource) > 0 {

		if check(alphanumericRegex, tag.CloudSource) {
			if check(alphanumericRegex, tag.CloudTag) {
				if check(alphanumericRegex, tag.TicketSource) {
					if check(alphanumericRegex, tag.TicketTag) {
						if check(alphanumericSpaceRegex, tag.Option) {
							valid = true
						}
					}
				}
			}
		}
	}

	return valid
}

func check(regex string, entry string) bool {
	match, err := regexp.MatchString(regex, entry)
	return err == nil && match
}
