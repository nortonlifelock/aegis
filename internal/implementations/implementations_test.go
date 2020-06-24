package implementations

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/aegis/internal/database"
	"github.com/nortonlifelock/aegis/internal/database/dal"
	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
	"regexp"
	"testing"
	"time"
)

const printLog = true

var functionNotImplementedError = fmt.Errorf("function not implemented")

type mockLogger struct {
	errStream chan error
}

func (m *mockLogger) Send(log log.Log) {
	var text string
	if printLog {
		if log.Error == nil {
			text = fmt.Sprintf("%s", log.Text)
		} else {
			text = fmt.Sprintf("%s|%v", log.Text, log.Error)
		}

	}

	if log.Type > 1 { // the log is more severe than a warning
		go func() {
			m.errStream <- fmt.Errorf(text)
		}()
	}
}

func TestScanCloseJob_modifyJiraTicketAccordingToVulnerabilityStatus(t *testing.T) {
	tests := []struct {
		ticket                      domain.Ticket
		scan                        domain.ScanSummary
		deadHostIPToProofMap        map[string]string
		deviceIDToVulnIDToDetection map[string]map[string]domain.Detection
		ipsForCloudDecommissionScan chan string
		funcGetDeviceByAssetOrgID   func(_AssetID string, OrgID string) (domain.Device, error)

		rescanType string

		expectedStatus string
		errExpected    bool
	}{

		// vuln in ticket not in detection slice - ticket should be closed
		{
			&dal.Ticket{
				DeviceIDvar:        "device1",
				IPAddressvar:       addressString("172.0.0.1"),
				VulnerabilityIDvar: "vuln2",
			},
			&dal.ScanSummary{},
			map[string]string{},
			map[string]map[string]domain.Detection{
				"device1": {
					"vuln1;": &mockDetection{},
				},
			},
			make(chan string),
			func(_AssetID string, OrgID string) (device domain.Device, e error) {
				return &mockDevice{
					valTrackingMethod: addressString(IPDevice),
				}, nil
			},
			domain.RescanNormal,
			domain.StatusClosedRemediated,
			false,
		},

		// fixed vuln in ticket is in detection slice - ticket should be closed
		{
			&dal.Ticket{
				DeviceIDvar:        "device1",
				IPAddressvar:       addressString("172.0.0.1"),
				VulnerabilityIDvar: "vuln1",
			},
			&dal.ScanSummary{},
			map[string]string{},
			map[string]map[string]domain.Detection{
				"device1": {
					"vuln1;": &mockDetection{
						valStatus: domain.Fixed,
					},
				},
			},
			make(chan string),
			func(_AssetID string, OrgID string) (device domain.Device, e error) {
				return &mockDevice{
					valTrackingMethod: addressString(IPDevice),
				}, nil
			},
			domain.RescanNormal,
			domain.StatusClosedRemediated,
			false,
		},

		// unfixed vuln in ticket is in detection slice - ticket should be reopened
		{
			&dal.Ticket{
				DeviceIDvar:        "device1",
				IPAddressvar:       addressString("172.0.0.1"),
				VulnerabilityIDvar: "vuln1",
			},
			&dal.ScanSummary{},
			map[string]string{},
			map[string]map[string]domain.Detection{
				"device1": {
					"vuln1;": &mockDetection{
						valStatus: domain.Vulnerable,
					},
				},
			},
			make(chan string),
			func(_AssetID string, OrgID string) (device domain.Device, e error) {
				return &mockDevice{
					valTrackingMethod: addressString(IPDevice),
				}, nil
			},
			domain.RescanNormal,
			domain.StatusReopened,
			false,
		},
	}

	for testIndex, test := range tests {
		scj, errStream := getBaseRescanCloseJob()

		scj.db = &mockDBWrapper{
			FuncGetDeviceByAssetOrgID: test.funcGetDeviceByAssetOrgID,
		}

		var engine = &mockTicketingEngine{
			funcGetStatusMap: func(backendStatus string) (equivalentTicketStatus string) {
				return backendStatus
			},
			funcTransition: func(ticket domain.Ticket, status string, comment string, Assignee string) (err error) {
				if status != test.expectedStatus {
					t.Errorf("[%d] mismatch between transitioned status and expected [%s|%s]", testIndex, status, test.expectedStatus)
				}
				return
			},
		}

		scj.Payload = &ScanClosePayload{}
		scj.Payload.Type = test.rescanType

		scj.modifyJiraTicketAccordingToVulnerabilityStatus(
			engine,
			test.ticket,
			test.scan,
			test.deadHostIPToProofMap,
			test.deviceIDToVulnIDToDetection,
			test.ipsForCloudDecommissionScan,
		)

		errSeen := streamHasErrors(errStream)

		if errSeen != test.errExpected {
			t.Errorf("[%d] mismatch in errSeen/errExpected [%v|%v]", testIndex, errSeen, test.errExpected)
		}
	}
}

func TestTicketingJob_getAssignmentInformation(t *testing.T) {
	tests := []struct {
		tagsForDevice   []domain.Tag
		payload         *vulnerabilityPayload
		assignmentRules []assignmentRule

		successFunction func(testIndex int, assignmentGroup, assignee *string)

		errExpected bool
	}{
		// testing rule matching on a regex
		{
			[]domain.Tag{},
			&vulnerabilityPayload{
				ticket: &dal.Ticket{
					IPAddressvar: addressString(""),
				},
				vuln: &mockVulnerability{
					valName: "a",
				},
			},
			[]assignmentRule{
				{
					AssignmentRules: &dal.AssignmentRules{
						Assigneevar:        addressString("1"),
						AssignmentGroupvar: addressString("2"),
					},
					vulnTitleRegex:        regexp.MustCompile("a"),
					excludeVulnTitleRegex: regexp.MustCompile("b"),
					hostnameRegex:         nil,
					osRegex:               nil,
					categoryRegex:         nil,
					tagKeyRegex:           nil,
					tagKey:                nil,
					ports:                 nil,
					excludePorts:          nil,
				},
			},
			func(testIndex int, assignmentGroup, assignee *string) {
				if assignee != nil && assignmentGroup != nil {
					if *assignee != "1" || *assignmentGroup != "2" {
						t.Errorf("[%d] assignment did not occur properly, got [%s|%s] instead of [1|2]", testIndex, *assignee, *assignmentGroup)
					}
				} else {
					t.Errorf("[%d] either assignee or assignmentGroup not set [%v|%v]", testIndex, assignee, assignmentGroup)
				}
			},
			false,
		},

		// testing rule matching on a regex that doesn't match
		{
			[]domain.Tag{},
			&vulnerabilityPayload{
				ticket: &dal.Ticket{
					IPAddressvar: addressString(""),
				},
				vuln: &mockVulnerability{
					valName: "b",
				},
			},
			[]assignmentRule{
				{
					AssignmentRules: &dal.AssignmentRules{
						Assigneevar:        addressString("1"),
						AssignmentGroupvar: addressString("2"),
					},
					vulnTitleRegex:        regexp.MustCompile("a"),
					excludeVulnTitleRegex: nil,
					hostnameRegex:         nil,
					osRegex:               nil,
					categoryRegex:         nil,
					tagKeyRegex:           nil,
					tagKey:                nil,
					ports:                 nil,
					excludePorts:          nil,
				},
			},
			func(testIndex int, assignmentGroup, assignee *string) {
				if assignee != nil || assignmentGroup != nil {
					t.Errorf("[%d] assignment did not occur properly, got [%v|%v] instead of [nil|nil]", testIndex, assignee, assignmentGroup)
				}
			},
			false,
		},

		// multiple rules that match - should take first
		{
			[]domain.Tag{},
			&vulnerabilityPayload{
				ticket: &dal.Ticket{
					IPAddressvar: addressString(""),
				},
				vuln: &mockVulnerability{
					valName: "ab",
				},
			},
			[]assignmentRule{
				{
					AssignmentRules: &dal.AssignmentRules{
						Assigneevar:        addressString("1"),
						AssignmentGroupvar: addressString("2"),
					},
					vulnTitleRegex:        regexp.MustCompile("a"),
					excludeVulnTitleRegex: nil,
					hostnameRegex:         nil,
					osRegex:               nil,
					categoryRegex:         nil,
					tagKeyRegex:           nil,
					tagKey:                nil,
					ports:                 nil,
					excludePorts:          nil,
				},
				{
					AssignmentRules: &dal.AssignmentRules{
						Assigneevar:        addressString("3"),
						AssignmentGroupvar: addressString("4"),
					},
					vulnTitleRegex:        regexp.MustCompile("b"),
					excludeVulnTitleRegex: nil,
					hostnameRegex:         nil,
					osRegex:               nil,
					categoryRegex:         nil,
					tagKeyRegex:           nil,
					tagKey:                nil,
					ports:                 nil,
					excludePorts:          nil,
				},
			},
			func(testIndex int, assignmentGroup, assignee *string) {
				if assignee != nil && assignmentGroup != nil {
					if *assignee != "1" || *assignmentGroup != "2" {
						t.Errorf("[%d] assignment did not occur properly, got [%s|%s] instead of [1|2]", testIndex, *assignee, *assignmentGroup)
					}
				} else {
					t.Errorf("[%d] either assignee or assignmentGroup not set [%v|%v]", testIndex, assignee, assignmentGroup)
				}
			},
			false,
		},

		// testing exclude regex rule
		{
			[]domain.Tag{},
			&vulnerabilityPayload{
				ticket: &dal.Ticket{
					IPAddressvar: addressString(""),
				},
				vuln: &mockVulnerability{
					valName: "ab",
				},
			},
			[]assignmentRule{
				{
					AssignmentRules: &dal.AssignmentRules{
						Assigneevar:        addressString("1"),
						AssignmentGroupvar: addressString("2"),
					},
					vulnTitleRegex:        regexp.MustCompile("a"),
					excludeVulnTitleRegex: regexp.MustCompile("b"),
					hostnameRegex:         nil,
					osRegex:               nil,
					categoryRegex:         nil,
					tagKeyRegex:           nil,
					tagKey:                nil,
					ports:                 nil,
					excludePorts:          nil,
				},
			},
			func(testIndex int, assignmentGroup, assignee *string) {
				if assignee != nil || assignmentGroup != nil {
					t.Errorf("[%d] assignment did not occur properly, got [%v|%v] instead of [nil|nil]", testIndex, assignee, assignmentGroup)
				}
			},
			false,
		},
	}

	for index, test := range tests {
		tj, errStream := getBaseTicketingJob()

		tj.db = &database.MockSQLDriver{}
		tj.assignmentRules = test.assignmentRules

		tj.getAssignmentInformation(test.tagsForDevice, test.payload)

		errSeen := streamHasErrors(errStream)

		if errSeen != test.errExpected {
			t.Errorf("[%d] mismatch in errSeen/errExpected [%v|%v]", index, errSeen, test.errExpected)
		}

		test.successFunction(index, test.payload.ticket.AssignmentGroup(), test.payload.ticket.AssignedTo())
	}
}

func TestTicketingJob_handleCloudTagMapping(t *testing.T) {
	var ticketHostname = "TICKETHOSTNAME"
	var tagHostname = "TAGHOSTNAME"

	var ticketAG = "TICKETAG"
	var tagAG = "TAGAG"

	tests := []struct {
		ticket domain.Ticket

		funcGetTagsForDevice func(_DeviceID string) ([]domain.Tag, error)
		funcGetTagKeyByID    func(_ID string) (key domain.TagKey, e error)
		tagMaps              []domain.TagMap

		successFunction func(testIndex int, ticket domain.Ticket, tagsForDevice []domain.Tag)
		errExpected     bool
	}{
		// testing Append option
		{
			&dal.Ticket{
				HostNamevar: &ticketHostname,
			},
			func(_DeviceID string) (tags []domain.Tag, e error) {
				return []domain.Tag{
					&dal.Tag{
						DeviceIDvar: "",
						IDvar:       "1",
						TagKeyIDvar: 0,
						Valuevar:    tagHostname,
					},
				}, nil
			},
			func(_ID string) (key domain.TagKey, e error) {
				return &dal.TagKey{
					KeyValuevar: "TEST",
				}, nil
			},
			[]domain.TagMap{
				&dal.TagMap{
					CloudSourceIDvar:     "",
					CloudTagvar:          "TEST",
					IDvar:                "",
					Optionsvar:           Append,
					TicketingSourceIDvar: "",
					TicketingTagvar:      "Hostname",
				},
			},
			func(testIndex int, ticket domain.Ticket, tagsForDevice []domain.Tag) {

				if ticket == nil {
					t.Errorf("[%d] nil ticket returned from method", testIndex)
				} else if ticket.HostName() == nil {
					t.Errorf("[%d] nil hostname found in ticket returned", testIndex)
				} else if *ticket.HostName() != fmt.Sprintf("%s,%s", ticketHostname, tagHostname) {
					t.Errorf("[%d] expected hostname and actual hostname differed [%s|%s]", testIndex, fmt.Sprintf("%s,%s", ticketHostname, tagHostname), *ticket.HostName())
				}
			},
			false,
		},

		// testing Overwrite option
		{
			&dal.Ticket{
				HostNamevar: &ticketHostname,
			},
			func(_DeviceID string) (tags []domain.Tag, e error) {
				return []domain.Tag{&dal.Tag{
					DeviceIDvar: "",
					IDvar:       "",
					TagKeyIDvar: 0,
					Valuevar:    tagHostname,
				}}, nil
			},
			func(_ID string) (key domain.TagKey, e error) {
				return &dal.TagKey{
					KeyValuevar: "TEST",
				}, nil
			},
			[]domain.TagMap{
				&dal.TagMap{
					CloudSourceIDvar:     "",
					CloudTagvar:          "TEST",
					IDvar:                "",
					Optionsvar:           Overwrite,
					TicketingSourceIDvar: "",
					TicketingTagvar:      "Hostname",
				},
			},
			func(testIndex int, ticket domain.Ticket, tagsForDevice []domain.Tag) {

				if ticket == nil {
					t.Errorf("[%d] nil ticket returned from method", testIndex)
				} else if ticket.HostName() == nil {
					t.Errorf("[%d] nil hostname found in ticket returned", testIndex)
				} else if *ticket.HostName() != tagHostname {
					t.Errorf("[%d] expected hostname and actual hostname differed [%s|%s]", testIndex, tagHostname, *ticket.HostName())
				}
			},
			false,
		},

		// testing no changes
		{
			&dal.Ticket{
				HostNamevar: &ticketHostname,
			},
			func(_DeviceID string) (tags []domain.Tag, e error) {
				return []domain.Tag{&dal.Tag{
					DeviceIDvar: "",
					IDvar:       "",
					TagKeyIDvar: 0,
					Valuevar:    tagHostname,
				}}, nil
			},
			func(_ID string) (key domain.TagKey, e error) {
				return &dal.TagKey{
					KeyValuevar: "TEST",
				}, nil
			},
			[]domain.TagMap{
				&dal.TagMap{
					CloudSourceIDvar:     "",
					CloudTagvar:          "TEST",
					IDvar:                "",
					Optionsvar:           Overwrite,
					TicketingSourceIDvar: "",
					TicketingTagvar:      "random value",
				},
			},
			func(testIndex int, ticket domain.Ticket, tagsForDevice []domain.Tag) {

				if ticket == nil {
					t.Errorf("[%d] nil ticket returned from method", testIndex)
				} else if ticket.HostName() == nil {
					t.Errorf("[%d] nil hostname found in ticket returned", testIndex)
				} else if *ticket.HostName() != ticketHostname {
					t.Errorf("[%d] expected hostname and actual hostname differed [%s|%s]", testIndex, ticketHostname, *ticket.HostName())
				}
			},
			false,
		},

		// testing multiple changes
		{
			&dal.Ticket{
				HostNamevar:        &ticketHostname,
				AssignmentGroupvar: &ticketAG,
			},
			func(_DeviceID string) (tags []domain.Tag, e error) {
				return []domain.Tag{
					&dal.Tag{
						DeviceIDvar: "",
						IDvar:       "",
						TagKeyIDvar: 1,
						Valuevar:    tagHostname,
					},
					&dal.Tag{
						DeviceIDvar: "",
						IDvar:       "",
						TagKeyIDvar: 2,
						Valuevar:    tagAG,
					},
				}, nil
			},
			func(_ID string) (key domain.TagKey, e error) {
				vals := map[string]domain.TagKey{
					"1": &dal.TagKey{KeyValuevar: "TEST1"},
					"2": &dal.TagKey{KeyValuevar: "TEST2"},
				}
				return vals[_ID], nil
			},
			[]domain.TagMap{
				&dal.TagMap{
					CloudSourceIDvar:     "",
					CloudTagvar:          "TEST1",
					IDvar:                "",
					Optionsvar:           Overwrite,
					TicketingSourceIDvar: "",
					TicketingTagvar:      "Hostname",
				},
				&dal.TagMap{
					CloudSourceIDvar:     "",
					CloudTagvar:          "TEST2",
					IDvar:                "",
					Optionsvar:           Overwrite,
					TicketingSourceIDvar: "",
					TicketingTagvar:      "AssignmentGroup",
				},
			},
			func(testIndex int, ticket domain.Ticket, tagsForDevice []domain.Tag) {

				if ticket == nil {
					t.Errorf("[%d] nil ticket returned from method", testIndex)
				} else if ticket.HostName() == nil {
					t.Errorf("[%d] nil hostname found in ticket returned", testIndex)
				} else if *ticket.HostName() != tagHostname {
					t.Errorf("[%d] expected hostname and actual hostname differed [%s|%s]", testIndex, tagHostname, *ticket.HostName())
				} else if *ticket.AssignmentGroup() != tagAG {
					t.Errorf("[%d] expected assignment group and actual assignment group differed [%s|%s]", testIndex, tagAG, *ticket.AssignmentGroup())
				}
			},
			false,
		},

		// testing tag return in absence of tag mapping
		{
			&dal.Ticket{
				HostNamevar: &ticketHostname,
			},
			func(_DeviceID string) (tags []domain.Tag, e error) {
				out := make([]domain.Tag, 0)

				for i := 0; i < 10; i++ {
					out = append(out, &dal.Tag{})
				}
				return out, nil
			},
			func(_ID string) (key domain.TagKey, e error) {
				return &dal.TagKey{
					KeyValuevar: "TEST",
				}, nil
			},
			[]domain.TagMap{},
			func(testIndex int, ticket domain.Ticket, tagsForDevice []domain.Tag) {

				if len(tagsForDevice) != 10 {
					t.Errorf("[%d] expected %d tags returned from method but got %d", testIndex, 10, len(tagsForDevice))
				}
			},
			false,
		},
	}

	for index, test := range tests {
		tj, errStream := getBaseTicketingJob()

		tj.tagMaps = test.tagMaps

		tj.db = &database.MockSQLDriver{
			FuncGetTagsForDevice: test.funcGetTagsForDevice,
			FuncGetTagKeyByID:    test.funcGetTagKeyByID,
		}

		ticket, tagsForDevice, err := tj.handleCloudTagMappings(test.ticket, &mockDevice{})

		errSeen := streamHasErrors(errStream)
		errSeen = errSeen || err != nil

		if errSeen != test.errExpected {
			t.Errorf("[%d] mismatch in errSeen/errExpected [%v|%v]", index, errSeen, test.errExpected)
		}

		test.successFunction(index, ticket, tagsForDevice)
	}

}

func TestAssetSyncJob_enterAssetInformationInDB(t *testing.T) {
	tests := []struct {
		asset                          domain.Device
		funcGetDeviceByAssetOrgID      func(_AssetID string, OrgID string) (domain.Device, error)
		funcGetDeviceByInstanceID      func(_InstanceID string, _OrgID string) (devices []domain.Device, e error)
		funcGetDeviceByScannerSourceID func(_IP string, _GroupID string, _OrgID string) (device domain.Device, e error)

		errExpected              bool
		expectCreateMethodCalled bool
		expectUpdateMethodCalled bool
	}{
		{
			&mockDevice{
				valSourceID: addressString("1"),
			},
			func(_AssetID string, OrgID string) (device domain.Device, e error) {
				return
			},
			func(_InstanceID string, _OrgID string) (devices []domain.Device, e error) {
				return
			},
			func(_IP string, _GroupID string, _OrgID string) (device domain.Device, e error) {
				return
			},
			false,
			true, // every db call returns a nil device
			false,
		},

		{
			&mockDevice{
				valSourceID: addressString("1"),
			},
			func(_AssetID string, OrgID string) (device domain.Device, e error) {
				return &mockDevice{
					valSourceID: addressString("2"), // different asset ID
				}, nil
			},
			func(_InstanceID string, _OrgID string) (devices []domain.Device, e error) {
				return
			},
			func(_IP string, _GroupID string, _OrgID string) (device domain.Device, e error) {
				return
			},
			false,
			true, // method call returning a device with a different asset id should still create a device
			false,
		},

		{
			&mockDevice{
				valSourceID: addressString("1"),
				valOS:       "different that what's in the DB so an update occurs",
			},
			func(_AssetID string, OrgID string) (device domain.Device, e error) {
				return &mockDevice{
					valSourceID: addressString("1"), // same asset ID
				}, nil
			},
			func(_InstanceID string, _OrgID string) (devices []domain.Device, e error) {
				return
			},
			func(_IP string, _GroupID string, _OrgID string) (device domain.Device, e error) {
				return
			},
			false,
			false,
			true, // method call returning a device with a same asset id should still update the device when the OS is different
		},

		{
			&mockDevice{
				valSourceID:   addressString("1"),
				valInstanceID: addressString("non-empty value so GetDeviceByInstanceID is called"),
			},
			func(_AssetID string, OrgID string) (device domain.Device, e error) {
				return &mockDevice{
					valSourceID: addressString(""), // device in DB must be empty to cause the update
				}, nil
			},
			func(_InstanceID string, _OrgID string) (devices []domain.Device, e error) {
				return
			},
			func(_IP string, _GroupID string, _OrgID string) (device domain.Device, e error) {
				return
			},
			false,
			false,
			true,
		},
	}

	for index, test := range tests {
		var updateMethodCalled, createMethodCalled bool

		asj, errStream := getBaseAssetSyncJob()

		asj.db = &mockDBWrapper{
			DatabaseConnection: &database.MockSQLDriver{
				FuncCreateDevice: func(_AssetID string, _SourceID string, _Ip string, _Hostname string, inInstanceID string, _MAC string, _GroupID string, _OrgID string, _OS string, _OSTypeID int, inTrackingMethod string) (id int, affectedRows int, err error) {
					createMethodCalled = true
					return
				},
				FuncUpdateAssetIDOsTypeIDOfDevice: func(_ID string, _AssetID string, _ScannerSourceID string, _GroupID string, _OS string, _HostName string, _OsTypeID int, inTrackingMethod string, _OrgID string) (id int, affectedRows int, err error) {
					updateMethodCalled = true
					return
				},
			},
			FuncGetDeviceByAssetOrgID:      test.funcGetDeviceByAssetOrgID,
			FuncGetDeviceByInstanceID:      test.funcGetDeviceByInstanceID,
			FuncGetDeviceByScannerSourceID: test.funcGetDeviceByScannerSourceID,
		}

		err := asj.enterAssetInformationInDB(test.asset, 1, "")
		errSeen := streamHasErrors(errStream)
		errSeen = errSeen || err != nil

		if errSeen != test.errExpected {
			t.Errorf("[%d] mismatch in errSeen/errExpected [%v|%v]", index, errSeen, test.errExpected)
		}

		if createMethodCalled != test.expectCreateMethodCalled {
			t.Errorf("[%d] mismatch in createMethodCalled/expectCreateMethodCalled [%v|%v]", index, createMethodCalled, test.expectCreateMethodCalled)
		}

		if updateMethodCalled != test.expectUpdateMethodCalled {
			t.Errorf("[%d] mismatch in updateMethodCalled/expectUpdateMethodCalled [%v|%v]", index, updateMethodCalled, test.expectUpdateMethodCalled)
		}
	}
}

func TestAssetSyncJob_getDecommIgnoreEntryForAsset(t *testing.T) {
	var ignoreEntryDeleted bool

	tests := []struct {
		funcHasDecommissioned          func(_devID string, _sourceID string, _orgID string) (domain.Ignore, error)
		funcDeleteDecomIgnoreForDevice func(_sourceID string, _devID string, _orgID string) (id int, affectedRows int, err error)
		detections                     func() []domain.Detection

		successFunction            func(testIndex int, ignoreID string)
		errExpected                bool
		ignoreEntryDeletedExpected bool
	}{
		// testing w/ no Ignore
		{
			func(_devID string, _sourceID string, _orgID string) (domain.Ignore, error) {
				return nil, nil
			},
			func(_sourceID string, _devID string, _orgID string) (id int, affectedRows int, err error) {
				t.Errorf("test should not try to delete decom ignore entry")
				return
			},
			func() []domain.Detection {
				return nil
			},
			func(testIndex int, ignoreID string) {
				if len(ignoreID) != 0 {
					t.Errorf("[%d] did not return an Ignore entry but received an IgnoreID [%s]", testIndex, ignoreID)
				}
			},
			false,
			false,
		},

		// testing an Ignore with a DueDate after the found date of all detections - decomm should not be deleted
		{
			func(_devID string, _sourceID string, _orgID string) (domain.Ignore, error) {
				val := time.Now()
				return &dal.Ignore{IDvar: "TEST", DueDatevar: &val}, nil
			},
			func(_sourceID string, _devID string, _orgID string) (id int, affectedRows int, err error) {
				t.Errorf("test should not try to delete decom ignore entry")
				return
			},
			func() []domain.Detection {
				var dets = make([]domain.Detection, 0)

				for i := 0; i < 10; i++ {

					val := time.Now().Add(time.Hour * -10)

					dets = append(dets, &mockDetection{
						valDetected: &val,
					})
				}

				return dets
			},
			func(testIndex int, ignoreID string) {
				if ignoreID != "TEST" {
					t.Errorf("[%d] Expected to find an IgnoreID", testIndex)
				}
			},
			false,
			false,
		},

		// testing an Ignore with a DueDate before the found date of all detections - decomm ignore should be deleted
		{
			func(_devID string, _sourceID string, _orgID string) (domain.Ignore, error) {
				val := time.Now()
				return &dal.Ignore{IDvar: "TEST", DueDatevar: &val}, nil
			},
			func(_sourceID string, _devID string, _orgID string) (id int, affectedRows int, err error) {
				ignoreEntryDeleted = true
				return
			},
			func() []domain.Detection {
				var dets = make([]domain.Detection, 0)

				for i := 0; i < 10; i++ {

					val := time.Now().Add(time.Hour * 10)

					dets = append(dets, &mockDetection{
						valDetected: &val,
					})
				}

				return dets
			},
			func(testIndex int, ignoreID string) {
				if len(ignoreID) != 0 {
					t.Errorf("[%d] did not return an Ignore entry but received an IgnoreID [%s]", testIndex, ignoreID)
				}
			},
			false,
			true,
		},
	}

	for index, test := range tests {
		ignoreEntryDeleted = false
		asj, errStream := getBaseAssetSyncJob()
		asj.db = &database.MockSQLDriver{
			FuncHasDecommissioned:          test.funcHasDecommissioned,
			FuncDeleteDecomIgnoreForDevice: test.funcDeleteDecomIgnoreForDevice,
		}

		ignoreID, err := asj.getDecommIgnoreEntryForAsset("test", "", test.detections())

		errSeen := streamHasErrors(errStream)
		errSeen = errSeen || err != nil

		if errSeen != test.errExpected {
			t.Errorf("[%d] mismatch in errSeen/errExpected [%v|%v]", index, errSeen, test.errExpected)
		}

		test.successFunction(index, ignoreID)

		if test.ignoreEntryDeletedExpected != ignoreEntryDeleted {
			t.Errorf("[%d] there was a disparity between the ignoreEntryDeletedExpected and ignoreEntryDeleted [%v|%v]", index, test.ignoreEntryDeletedExpected, ignoreEntryDeleted)
		}
	}
}

func TestAssetSyncJob_fanInDetections(t *testing.T) {
	tests := []struct {
		dets            func() <-chan domain.Detection
		successFunction func(testIndex int, funcDevIDToDetection map[string][]domain.Detection)
		errExpected     bool
	}{
		{
			func() <-chan domain.Detection {
				out := make(chan domain.Detection)
				go func() {
					defer close(out)

					out <- &mockDetection{
						valDevice: nil, // should cause error
					}
				}()

				return out
			},
			func(testIndex int, funcDevIDToDetection map[string][]domain.Detection) {
				return
			},
			true,
		},

		{
			func() <-chan domain.Detection {
				out := make(chan domain.Detection)
				go func() {
					defer close(out)

					for i := 0; i < 10; i++ {
						out <- &mockDetection{
							valDevice: &mockDevice{valSourceID: addressString("1")},
						}
					}

					for i := 0; i < 100; i++ {
						out <- &mockDetection{
							valDevice: &mockDevice{valSourceID: addressString("2")},
						}
					}

					for i := 0; i < 1000; i++ {
						out <- &mockDetection{
							valDevice: &mockDevice{valSourceID: addressString("3")},
						}
					}
				}()

				return out
			},
			func(testIndex int, funcDevIDToDetection map[string][]domain.Detection) {
				if len(funcDevIDToDetection["1"]) != 10 || len(funcDevIDToDetection["2"]) != 100 || len(funcDevIDToDetection["3"]) != 1000 {
					t.Errorf("[%d] expected 10 detections for device [1], 100 detections for device [2], 1000 detections for device [3], got [%d|%d|%d]",
						testIndex, len(funcDevIDToDetection["1"]), len(funcDevIDToDetection["2"]), len(funcDevIDToDetection["3"]))
				}
			},
			false,
		},

		{
			func() <-chan domain.Detection {
				out := make(chan domain.Detection)
				go func() {
					defer close(out)

					for i := 0; i < 1000; i++ {
						out <- &mockDetection{
							valDevice: &mockDevice{valSourceID: addressString("1")},
							valPort:   i,
						}
					}
				}()

				return out
			},
			func(testIndex int, funcDevIDToDetection map[string][]domain.Detection) {
				for i := 0; i < 1000; i++ {
					var found bool
					for _, detection := range funcDevIDToDetection["1"] {
						if detection.Port() == i {
							found = true
							break
						}
					}

					if !found {
						t.Errorf("[%d] expected to see a detection with the port [%d] but could not find one", testIndex, i)
					}
				}
			},
			false,
		},
	}

	for index, test := range tests {
		asj, errStream := getBaseAssetSyncJob()
		devIDToDetection := asj.fanInDetections(test.dets())
		errSeen := streamHasErrors(errStream)

		if errSeen != test.errExpected {
			t.Errorf("[%d] mismatch in errSeen/errExpected [%v|%v]", index, errSeen, test.errExpected)
		}

		test.successFunction(index, devIDToDetection)
	}
}

type mockScanner struct {
	integrations.Vscanner
	funcDetections func(ctx context.Context, groupsIDs []string) (detections <-chan domain.Detection, err error)
}

func (s *mockScanner) Detections(ctx context.Context, groupsIDs []string) (detections <-chan domain.Detection, err error) {
	if s.funcDetections != nil {
		return s.funcDetections(ctx, groupsIDs)
	} else {
		return nil, functionNotImplementedError
	}
}

func getContext() context.Context {
	return context.Background()
}

type mockDetection struct {
	valID              string
	valVulnerabilityID string
	valStatus          string
	valActiveKernel    *int
	valDetected        *time.Time
	valTimesSeen       int
	valProof           string
	valPort            int
	valProtocol        string
	valIgnoreID        *string
	valLastFound       *time.Time
	valLastUpdated     *time.Time
	valDevice          domain.Device
	valVulnerability   domain.Vulnerability
}

func (d *mockDetection) ID() string {
	return d.valID
}

func (d *mockDetection) VulnerabilityID() string {
	return d.valVulnerabilityID
}

func (d *mockDetection) Status() string {
	return d.valStatus
}

func (d *mockDetection) ActiveKernel() *int {
	return d.valActiveKernel
}

func (d *mockDetection) Detected() (*time.Time, error) {
	return d.valDetected, nil
}

func (d *mockDetection) TimesSeen() int {
	return d.valTimesSeen
}

func (d *mockDetection) Proof() string {
	return d.valProof
}

func (d *mockDetection) Port() int {
	return d.valPort
}

func (d *mockDetection) Protocol() string {
	return d.valProtocol
}

func (d *mockDetection) IgnoreID() (*string, error) {
	return d.valIgnoreID, nil
}

func (d *mockDetection) LastFound() *time.Time {
	return d.valLastFound
}

func (d *mockDetection) LastUpdated() *time.Time {
	return d.valLastUpdated
}

func (d *mockDetection) Device() (domain.Device, error) {
	return d.valDevice, nil
}

func (d *mockDetection) Vulnerability() (domain.Vulnerability, error) {
	return d.valVulnerability, nil
}

type mockVulnerability struct {
	valID                   string
	valSourceID             string
	valName                 string
	valDescription          string
	valCVSS2                float32
	valCVSS3                *float32
	valUpdated              time.Time
	valSoftware             string
	valPatchable            *string
	valDetectionInformation string
	valThreat               *string
	valCategory             *string
}

func (m *mockVulnerability) ID() string {
	return m.valID
}

func (m *mockVulnerability) SourceID() string {
	return m.valSourceID
}

func (m *mockVulnerability) Name() string {
	return m.valName
}

func (m *mockVulnerability) Description() string {
	return m.valDescription
}

func (m *mockVulnerability) CVSS2() float32 {
	return m.valCVSS2
}

func (m *mockVulnerability) CVSS3() *float32 {
	return m.valCVSS3
}

func (m *mockVulnerability) Updated() time.Time {
	return m.valUpdated
}

func (m *mockVulnerability) Solutions(ctx context.Context) (<-chan domain.Solution, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockVulnerability) References(ctx context.Context) (<-chan domain.VulnerabilityReference, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockVulnerability) Software() string {
	return m.valSoftware
}

func (m *mockVulnerability) Patchable() *string {
	return m.valPatchable
}

func (m *mockVulnerability) DetectionInformation() string {
	return m.valDetectionInformation
}

func (m *mockVulnerability) Threat() *string {
	return m.valThreat
}

func (m *mockVulnerability) Category() (param *string) {
	return m.valCategory
}

type mockDevice struct {
	valID             string
	valSourceID       *string
	valOS             string
	valHostName       string
	valMAC            string
	valIP             string
	valRegion         *string
	valInstanceID     *string
	valGroupID        *string
	valTrackingMethod *string
}

func (m *mockDevice) ID() string {
	return m.valID
}
func (m *mockDevice) SourceID() *string {
	return m.valSourceID
}
func (m *mockDevice) OS() string {
	return m.valOS
}
func (m *mockDevice) HostName() string {
	return m.valHostName
}
func (m *mockDevice) MAC() string {
	return m.valMAC
}
func (m *mockDevice) IP() string {
	return m.valIP
}
func (m *mockDevice) Vulnerabilities(ctx context.Context) (param <-chan domain.Detection, err error) {
	return nil, fmt.Errorf("not implemented")
}
func (m *mockDevice) Region() *string {
	return m.valRegion
}
func (m *mockDevice) InstanceID() *string {
	return m.valInstanceID
}
func (m *mockDevice) GroupID() *string {
	return m.valGroupID
}
func (m *mockDevice) TrackingMethod() *string {
	return m.valTrackingMethod
}

func addressString(in string) *string {
	return &in
}

type mockDBWrapper struct {
	domain.DatabaseConnection
	FuncGetDeviceByAssetOrgID         func(_AssetID string, OrgID string) (domain.Device, error)
	FuncGetDeviceByIP                 func(_IP string, _OrgID string) (domain.Device, error)
	FuncGetDeviceByCloudSourceIDAndIP func(_IP string, _CloudSourceID string, _OrgID string) ([]domain.Device, error)
	FuncGetDeviceByScannerSourceID    func(_IP string, _GroupID string, _OrgID string) (domain.Device, error)
	FuncGetDeviceByInstanceID         func(_InstanceID string, _OrgID string) ([]domain.Device, error)
	FuncGetDevicesBySourceID          func(_SourceID string, _OrgID string) ([]domain.Device, error)
	FuncGetDevicesByCloudSourceID     func(_CloudSourceID string, _OrgID string) ([]domain.Device, error)
	FuncGetDetection                  func(_SourceDeviceID string, _VulnerabilityID string, _Port int, _Protocol string) (domain.Detection, error)
	FuncGetDetectionBySourceVulnID    func(_SourceDeviceID string, _SourceVulnerabilityID string, _Port int, _Protocol string) (domain.Detection, error)
	FuncGetDetectionsForDevice        func(_DeviceID string) ([]domain.Detection, error)
	FuncGetDetectionsAfter            func(after time.Time, orgID string) (detections []domain.Detection, err error)
	FuncGetDetectionForGroupAfter     func(_After time.Time, _OrgID string, inGroupID string, ticketInactiveKernels bool) ([]domain.Detection, error)
	FuncGetVulnReferences             func(vulnInfoID string, sourceID string) (references []domain.VulnerabilityReference, err error)
	FuncGetVulnRef                    func(vulnInfoID string, sourceID string, reference string) (existing domain.VulnerabilityReference, err error)
	FuncGetVulnBySourceVulnID         func(_SourceVulnID string) (vulnerability domain.Vulnerability, err error)
}

func (m *mockDBWrapper) GetDeviceByAssetOrgID(_AssetID string, OrgID string) (domain.Device, error) {
	if m.FuncGetDeviceByAssetOrgID != nil {
		return m.FuncGetDeviceByAssetOrgID(_AssetID, OrgID)
	} else {
		panic("method not implemented")
	}
}
func (m *mockDBWrapper) GetDeviceByIP(_IP string, _OrgID string) (domain.Device, error) {
	if m.FuncGetDeviceByIP != nil {
		return m.FuncGetDeviceByIP(_IP, _OrgID)
	} else {
		panic("method not implemented")
	}
}
func (m *mockDBWrapper) GetDeviceByCloudSourceIDAndIP(_IP string, _CloudSourceID string, _OrgID string) ([]domain.Device, error) {
	if m.FuncGetDeviceByCloudSourceIDAndIP != nil {
		return m.FuncGetDeviceByCloudSourceIDAndIP(_IP, _CloudSourceID, _OrgID)
	} else {
		panic("method not implemented")
	}
}
func (m *mockDBWrapper) GetDeviceByScannerSourceID(_IP string, _GroupID string, _OrgID string) (domain.Device, error) {
	if m.FuncGetDeviceByScannerSourceID != nil {
		return m.FuncGetDeviceByScannerSourceID(_IP, _GroupID, _OrgID)
	} else {
		panic("method not implemented")
	}
}
func (m *mockDBWrapper) GetDeviceByInstanceID(_InstanceID string, _OrgID string) ([]domain.Device, error) {
	if m.FuncGetDeviceByInstanceID != nil {
		return m.FuncGetDeviceByInstanceID(_InstanceID, _OrgID)
	} else {
		panic("method not implemented")
	}
}
func (m *mockDBWrapper) GetDevicesBySourceID(_SourceID string, _OrgID string) ([]domain.Device, error) {
	if m.FuncGetDevicesBySourceID != nil {
		return m.FuncGetDevicesBySourceID(_SourceID, _OrgID)
	} else {
		panic("method not implemented")
	}
}
func (m *mockDBWrapper) GetDevicesByCloudSourceID(_CloudSourceID string, _OrgID string) ([]domain.Device, error) {
	if m.FuncGetDevicesByCloudSourceID != nil {
		return m.FuncGetDevicesByCloudSourceID(_CloudSourceID, _OrgID)
	} else {
		panic("method not implemented")
	}
}
func (m *mockDBWrapper) GetDetection(_SourceDeviceID string, _VulnerabilityID string, _Port int, _Protocol string) (domain.Detection, error) {
	if m.FuncGetDetection != nil {
		return m.FuncGetDetection(_SourceDeviceID, _VulnerabilityID, _Port, _Protocol)
	} else {
		panic("method not implemented")
	}
}
func (m *mockDBWrapper) GetDetectionBySourceVulnID(_SourceDeviceID string, _SourceVulnerabilityID string, _Port int, _Protocol string) (domain.Detection, error) {
	if m.FuncGetDetectionBySourceVulnID != nil {
		return m.FuncGetDetectionBySourceVulnID(_SourceDeviceID, _SourceVulnerabilityID, _Port, _Protocol)
	} else {
		panic("method not implemented")
	}
}
func (m *mockDBWrapper) GetDetectionsForDevice(_DeviceID string) ([]domain.Detection, error) {
	if m.FuncGetDetectionsForDevice != nil {
		return m.FuncGetDetectionsForDevice(_DeviceID)
	} else {
		panic("method not implemented")
	}
}
func (m *mockDBWrapper) GetDetectionsAfter(after time.Time, orgID string) (detections []domain.Detection, err error) {
	if m.FuncGetDetectionsAfter != nil {
		return m.FuncGetDetectionsAfter(after, orgID)
	} else {
		panic("method not implemented")
	}
}
func (m *mockDBWrapper) GetDetectionForGroupAfter(_After time.Time, _OrgID string, inGroupID string, ticketInactiveKernels bool) ([]domain.Detection, error) {
	if m.FuncGetDetectionForGroupAfter != nil {
		return m.FuncGetDetectionForGroupAfter(_After, _OrgID, inGroupID, ticketInactiveKernels)
	} else {
		panic("method not implemented")
	}
}
func (m *mockDBWrapper) GetVulnReferences(vulnInfoID string, sourceID string) (references []domain.VulnerabilityReference, err error) {
	if m.FuncGetVulnReferences != nil {
		return m.FuncGetVulnReferences(vulnInfoID, sourceID)
	} else {
		panic("method not implemented")
	}
}
func (m *mockDBWrapper) GetVulnRef(vulnInfoID string, sourceID string, reference string) (existing domain.VulnerabilityReference, err error) {
	if m.FuncGetVulnRef != nil {
		return m.FuncGetVulnRef(vulnInfoID, sourceID, reference)
	} else {
		panic("method not implemented")
	}
}
func (m *mockDBWrapper) GetVulnBySourceVulnID(_SourceVulnID string) (vulnerability domain.Vulnerability, err error) {
	if m.FuncGetVulnBySourceVulnID != nil {
		return m.FuncGetVulnBySourceVulnID(_SourceVulnID)
	} else {
		panic("method not implemented")
	}
}

func getBaseRescanCloseJob() (*ScanCloseJob, chan error) {
	var errStream = make(chan error)
	sc := &ScanCloseJob{}
	sc.ctx = getContext()
	sc.lstream = &mockLogger{
		errStream,
	}
	sc.insource = &dal.SourceConfig{}
	sc.outsource = &dal.SourceConfig{}
	sc.config = &dal.JobConfig{}
	return sc, errStream
}

func getBaseTicketingJob() (*TicketingJob, chan error) {
	var errStream = make(chan error)
	tj := &TicketingJob{}
	tj.ctx = getContext()
	tj.lstream = &mockLogger{
		errStream,
	}
	tj.insource = &dal.SourceConfig{}
	tj.outsource = &dal.SourceConfig{}
	tj.config = &dal.JobConfig{}
	return tj, errStream
}

func getBaseAssetSyncJob() (*AssetSyncJob, chan error) {
	var errStream = make(chan error)
	asj := &AssetSyncJob{}
	asj.ctx = getContext()
	asj.lstream = &mockLogger{
		errStream,
	}
	asj.insources = &dal.SourceConfig{}
	asj.outsource = &dal.SourceConfig{}
	asj.config = &dal.JobConfig{}
	return asj, errStream
}

func streamHasErrors(errStream chan error) bool {
	time.Sleep(time.Millisecond * 10) // give a moment for asynchronous logger to push logs onto stream
	var errSeen bool
	select {
	case _, ok := <-errStream:
		if ok {
			errSeen = true
		}
	default:
	}

	return errSeen
}

type mockTicketingEngine struct {
	funcCreateTicket               func(ticket domain.Ticket) (sourceID int, sourceKey string, err error)
	funcUpdateTicket               func(ticket domain.Ticket, comment string) (sourceID int, sourceKey string, err error)
	funcTransition                 func(ticket domain.Ticket, status string, comment string, Assignee string) (err error)
	funcGetTicket                  func(sourceKey string) (ticket domain.Ticket, err error)
	funcGetTicketsByClosedStatus   func(orgCode string, methodOfDiscovery string, startDate time.Time) (tix <-chan domain.Ticket)
	funcGetTicketsUpdatedSince     func(since time.Time, orgCode string, methodOfDiscovery string) <-chan domain.Ticket
	funcGetTicketsForRescan        func(cerfs []domain.CERF, methodOfDiscovery string, orgCode string, algorithm string) (issues <-chan domain.Ticket, err error)
	funcGetTicketsByDeviceIDVulnID func(methodOfDiscovery string, orgCode string, deviceID string, vulnID string, statuses map[string]bool, port int, protocol string) (issues <-chan domain.Ticket, err error)
	funcGetCERFExpirationUpdates   func(startDate time.Time) (cerfs map[string]time.Time, err error)
	funcGetOpenTicketsByGroupID    func(methodOfDiscovery string, orgCode string, groupID string) (tickets <-chan domain.Ticket, err error)
	funcGetRelatedTicketsForRescan func(tickets []domain.Ticket, groupID string, methodOfDiscovery string, orgCode string, rescanType string) (issues <-chan domain.Ticket, err error)
	funcAssignmentGroupExists      func(groupName string) (exists bool, err error)
	funcGetStatusMap               func(backendStatus string) (equivalentTicketStatus string)
}

func (m *mockTicketingEngine) CreateTicket(ticket domain.Ticket) (sourceID int, sourceKey string, err error) {
	if m.funcCreateTicket != nil {
		return m.funcCreateTicket(ticket)
	} else {
		panic("method not implemented")
	}
}
func (m *mockTicketingEngine) UpdateTicket(ticket domain.Ticket, comment string) (sourceID int, sourceKey string, err error) {
	if m.funcUpdateTicket != nil {
		return m.funcUpdateTicket(ticket, comment)
	} else {
		panic("method not implemented")
	}
}
func (m *mockTicketingEngine) Transition(ticket domain.Ticket, status string, comment string, Assignee string) (err error) {
	if m.funcTransition != nil {
		return m.funcTransition(ticket, status, comment, Assignee)
	} else {
		panic("method not implemented")
	}
}
func (m *mockTicketingEngine) GetTicket(sourceKey string) (ticket domain.Ticket, err error) {
	if m.funcGetTicket != nil {
		return m.funcGetTicket(sourceKey)
	} else {
		panic("method not implemented")
	}
}
func (m *mockTicketingEngine) GetTicketsByClosedStatus(orgCode string, methodOfDiscovery string, startDate time.Time) (tix <-chan domain.Ticket) {
	if m.funcGetTicketsByClosedStatus != nil {
		return m.funcGetTicketsByClosedStatus(orgCode, methodOfDiscovery, startDate)
	} else {
		panic("method not implemented")
	}
}
func (m *mockTicketingEngine) GetTicketsUpdatedSince(since time.Time, orgCode string, methodOfDiscovery string) <-chan domain.Ticket {
	if m.funcGetTicketsUpdatedSince != nil {
		return m.funcGetTicketsUpdatedSince(since, orgCode, methodOfDiscovery)
	} else {
		panic("method not implemented")
	}
}
func (m *mockTicketingEngine) GetTicketsForRescan(cerfs []domain.CERF, methodOfDiscovery string, orgCode string, algorithm string) (issues <-chan domain.Ticket, err error) {
	if m.funcGetTicketsForRescan != nil {
		return m.funcGetTicketsForRescan(cerfs, methodOfDiscovery, orgCode, algorithm)
	} else {
		panic("method not implemented")
	}
}
func (m *mockTicketingEngine) GetTicketsByDeviceIDVulnID(methodOfDiscovery string, orgCode string, deviceID string, vulnID string, statuses map[string]bool, port int, protocol string) (issues <-chan domain.Ticket, err error) {
	if m.funcGetTicketsByDeviceIDVulnID != nil {
		return m.funcGetTicketsByDeviceIDVulnID(methodOfDiscovery, orgCode, deviceID, vulnID, statuses, port, protocol)
	} else {
		panic("method not implemented")
	}
}
func (m *mockTicketingEngine) GetCERFExpirationUpdates(startDate time.Time) (cerfs map[string]time.Time, err error) {
	if m.funcGetCERFExpirationUpdates != nil {
		return m.funcGetCERFExpirationUpdates(startDate)
	} else {
		panic("method not implemented")
	}
}
func (m *mockTicketingEngine) GetOpenTicketsByGroupID(methodOfDiscovery string, orgCode string, groupID string) (tickets <-chan domain.Ticket, err error) {
	if m.funcGetOpenTicketsByGroupID != nil {
		return m.funcGetOpenTicketsByGroupID(methodOfDiscovery, orgCode, groupID)
	} else {
		panic("method not implemented")
	}
}
func (m *mockTicketingEngine) GetRelatedTicketsForRescan(tickets []domain.Ticket, groupID string, methodOfDiscovery string, orgCode string, rescanType string) (issues <-chan domain.Ticket, err error) {
	if m.funcGetRelatedTicketsForRescan != nil {
		return m.funcGetRelatedTicketsForRescan(tickets, groupID, methodOfDiscovery, orgCode, rescanType)
	} else {
		panic("method not implemented")
	}
}
func (m *mockTicketingEngine) AssignmentGroupExists(groupName string) (exists bool, err error) {
	if m.funcAssignmentGroupExists != nil {
		return m.funcAssignmentGroupExists(groupName)
	} else {
		panic("method not implemented")
	}
}
func (m *mockTicketingEngine) GetStatusMap(backendStatus string) (equivalentTicketStatus string) {
	if m.funcGetStatusMap != nil {
		return m.funcGetStatusMap(backendStatus)
	} else {
		panic("method not implemented")
	}
}
