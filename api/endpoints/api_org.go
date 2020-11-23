package endpoints

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/nortonlifelock/aegis/pkg/domain"
	"github.com/pkg/errors"
)

func getMyOrg(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getMyOrgEndpoint, allAllowed, func(trans *transaction) {
		var organization domain.Organization
		organization, trans.err = Ms.GetOrganizationByID(trans.permission.OrgID())
		if trans.err == nil {
			trans.obj = toOrgDto(organization)
			trans.status = http.StatusOK
		} else {
			(&trans.wrapper).addError(trans.err, databaseError)
		}
	})
}

func getOrgForUser(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getOrgForUserEndpoint, allAllowed, func(trans *transaction) {
		var organizations []domain.Organization
		organizations, trans.err = Ms.GetLeafOrganizationsForUser(trans.user.ID())
		if trans.err == nil {
			trans.obj = toOrgDtoSlice(organizations)
			trans.status = http.StatusOK
		} else {
			(&trans.wrapper).addError(trans.err, databaseError)
		}
	})
}

func getAllOrganizations(w http.ResponseWriter, r *http.Request) {
	executeTransaction(w, r, getAllOrgsEndpoint, allAllowed, func(trans *transaction) {
		// TODO are permissions required for this? I don't think they should be...
		var organizations []domain.Organization
		organizations, trans.err = Ms.GetOrganizations()
		if trans.err == nil {
			trans.status = http.StatusOK
			trans.obj = toOrgDtoSlice(organizations)
		} else {
			(&trans.wrapper).addError(trans.err, processError)
		}
	})
}

func updateOrg(w http.ResponseWriter, r *http.Request) {

	executeTransaction(w, r, updateOrgEndpoint, admin|manager, func(trans *transaction) {
		if trans.endpoint.verify() {
			if organization, isOrg := trans.endpoint.(*Organization); isOrg {
				trans.obj, trans.status, trans.err = organization.update(trans.user, trans.permission, trans.originalBody)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = fmt.Errorf("tried to pass a non-organization as an organization")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = fmt.Errorf("invalid length of format of apiRequest fields")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func deleteOrg(w http.ResponseWriter, r *http.Request) {

	executeTransaction(w, r, deleteOrgEndpoint, admin, func(trans *transaction) {
		if trans.endpoint.verify() {
			if organization, isOrg := trans.endpoint.(*Organization); isOrg {
				trans.obj, trans.status, trans.err = organization.delete(trans.user, trans.permission)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = fmt.Errorf("tried to pass a non-organization as an organization")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = fmt.Errorf("invalid length of format of apiRequest fields")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func createOrg(w http.ResponseWriter, r *http.Request) {

	executeTransaction(w, r, createOrgEndpoint, admin|manager, func(trans *transaction) {
		if trans.endpoint.verify() {
			if organization, isOrg := trans.endpoint.(*Organization); isOrg {
				trans.obj, trans.status, trans.err = organization.create(trans.user, trans.permission)
				(&trans.wrapper).addError(trans.err, processError)
			} else {
				trans.err = fmt.Errorf("tried to pass a non-organization as an organization")
				(&trans.wrapper).addError(trans.err, requestFormatError)
			}
		} else {
			trans.err = fmt.Errorf("invalid length of format of apiRequest fields")
			(&trans.wrapper).addError(trans.err, requestFormatError)
		}
	})
}

func (organization *Organization) create(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}

	status = http.StatusBadRequest

	_, _, err = Ms.CreateOrganization(organization.Code, organization.Description,
		organization.TimeZoneOffset, sord(user.Username()))
	if err == nil {
		status = http.StatusOK
		generalResp.Response = fmt.Sprintf("organization created")
	} else {
		err = errors.New("error while creating org")
	}

	return generalResp, status, err
}

func (organization *Organization) update(user domain.User, permission domain.Permission, originalBody string) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}

	status = http.StatusBadRequest

	var existingOrg domain.Organization
	existingOrg, err = Ms.GetOrganizationByID(organization.ID)
	if err == nil {
		if existingOrg != nil {
			var bodyString = originalBody
			var description = organization.Description
			var offset = organization.TimeZoneOffset

			if strings.Index(bodyString, "description") < 0 {
				description = sord(existingOrg.Description())
			}
			if strings.Index(bodyString, "timezone_offset") < 0 {
				offset = existingOrg.TimeZoneOffset()
			}

			_, _, err = Ms.UpdateOrganization(organization.ID, description, offset, sord(user.Username()))

			if err == nil {
				status = http.StatusOK
				generalResp.Response = fmt.Sprintf("organization updated")
			} else {
				err = errors.New("error while updating org")
			}
		} else {
			err = errors.Errorf("could not find organization with id [%s]", organization.ID)
		}
	} else {
		err = errors.New("error while grabbing organization from database")
	}

	return generalResp, status, err
}

func (organization *Organization) delete(user domain.User, permission domain.Permission) (generalResp *GeneralResp, status int, err error) {
	generalResp = &GeneralResp{}

	status = http.StatusBadRequest

	_, _, err = Ms.DisableOrganization(organization.ID, sord(user.Username()))

	if err == nil {
		status = http.StatusOK
		generalResp.Response = fmt.Sprintf("organization deleted")
	} else {
		err = errors.New("error while deleting org")
	}

	return generalResp, status, err
}

func (organization *Organization) verify() (valid bool) {
	valid = false
	if organization.TimeZoneOffset >= -12 && organization.TimeZoneOffset <= 12 {
		if len(organization.Code) > 0 && len(organization.Code) <= upperBoundNameLen {
			if len(organization.Description) > 0 && len(organization.Description) < upperBoundNameLen {
				if len(organization.ID) > 0 {
					if match, err := regexp.MatchString(lettersRegex, organization.Code); err == nil && match {
						if match, err = regexp.MatchString(lettersRegex, organization.Description); err == nil && match {
							valid = true
						}
					}
				}
			}
		}
	}

	return valid
}
