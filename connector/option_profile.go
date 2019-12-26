package connector

import (
	"fmt"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/qualys"
	"strconv"
	"strings"
	"time"
)

func (session *QsSession) deleteSearchListAndOptionProfile(scan domain.ScanSummary) (err error) {
	if len(sord(scan.TemplateID())) > 0 {
		var optionProfileID, searchListID string
		var templateIDs = strings.Split(sord(scan.TemplateID()), templateDelimiter)
		if len(templateIDs) == 2 {
			optionProfileID, searchListID = templateIDs[0], templateIDs[1]
			if len(optionProfileID) > 0 {
				err = session.apiSession.DeleteOptionProfile(optionProfileID)
				if err == nil {
					if len(searchListID) > 0 {
						err = session.apiSession.DeleteSearchList(searchListID)
						if err != nil {
							err = fmt.Errorf("error while deleting search list - %s", err.Error())
						}
					} else {
						// intentionally left blank - there is no search list on a discovery scan
					}
				} else {
					err = fmt.Errorf("error while deleting option profile - %s", err.Error())
				}
			} else {
				err = fmt.Errorf("option profile ID not found for scan")
			}
		} else {
			err = fmt.Errorf("should have only had 2 template fields but found %d", len(templateIDs))
		}
	} else {
		err = fmt.Errorf("could not find option profile id or search list id in the scan summary")
	}

	return err
}

func (session *QsSession) createCopyOfOptionProfile(optionProfileToCopy int) (optionProfileID string, err error) {
	var optionProfileTemplate *qualys.OptionProfiles
	if optionProfileTemplate, err = session.apiSession.GetOptionProfile(optionProfileToCopy); err == nil {
		if optionProfileTemplate.OptionProfile.BasicInfo.ID == strconv.Itoa(optionProfileToCopy) {
			var title = fmt.Sprintf(session.payload.OptionProfileFormatString, strconv.Itoa(time.Now().Nanosecond()))
			optionProfileTemplate.OptionProfile.BasicInfo.GroupName = title
			optionProfileID, err = session.apiSession.CreateOptionProfile(optionProfileTemplate)
		} else {
			err = fmt.Errorf("could not find option profile [%v]", optionProfileToCopy)
		}
	}

	return optionProfileID, err
}

func (session *QsSession) createOptionProfileWithSearchList(QIDs []string, optionProfileToCopy int) (optionProfileID string, searchListID string, err error) {
	if searchListID, err = session.apiSession.CreateSearchList(QIDs, session.payload.SearchListFormatString); err == nil {
		var optionProfileTemplate *qualys.OptionProfiles
		if optionProfileTemplate, err = session.apiSession.GetOptionProfile(optionProfileToCopy); err == nil {

			var title = fmt.Sprintf(session.payload.OptionProfileFormatString, strconv.Itoa(time.Now().Nanosecond()))
			optionProfileTemplate.OptionProfile.BasicInfo.GroupName = title

			var foundSearchListToOverwrite bool
			// need to find the search list template to overwrite
			for index, searchList := range optionProfileTemplate.OptionProfile.Scan.VulnerabilityDetection.CustomList.Custom {
				if searchList.ID == strconv.Itoa(session.payload.SearchListID) {
					foundSearchListToOverwrite = true
					var overwriteSearchList = &optionProfileTemplate.OptionProfile.Scan.VulnerabilityDetection.CustomList.Custom[index]
					overwriteSearchList.ID = searchListID
					overwriteSearchList.Title = title
					break
				}
			}

			if foundSearchListToOverwrite {
				optionProfileID, err = session.apiSession.CreateOptionProfile(optionProfileTemplate)
			} else {
				// failed while setting up the option profile, clean up by deleting the search list that will not be used
				err = session.apiSession.DeleteSearchList(searchListID)
				if err == nil {
					err = fmt.Errorf("failed to find the search list ID to overwrite [%d] in the option profile", session.payload.SearchListID)
				} else {
					err = fmt.Errorf("failed to find the search list ID to overwrite [%d] in the option profile AND failure while deleting search list [%s]", session.payload.SearchListID, searchListID)
				}
			}

		} else {
			err = fmt.Errorf("error while gathering the option profile template - %s", err.Error())
		}
	} else {
		err = fmt.Errorf("error while creating the search list for the scan - %s", err.Error())
	}

	return optionProfileID, searchListID, err
}
