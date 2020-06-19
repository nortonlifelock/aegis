package implementations

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/aegis/internal/database"
	"github.com/nortonlifelock/aegis/internal/database/dal"
	"github.com/nortonlifelock/aegis/internal/integrations"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
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

func TestTicketingJob_handleCloudTagMapping(t *testing.T) {
	tests := []struct {
		ticket domain.Ticket
		device domain.Device

		funcGetTagsForDevice func(_DeviceID string) ([]domain.Tag, error)
		funcGetTagKeyByID    func(_ID string) (key domain.TagKey, e error)
		tagMaps              []domain.TagMap

		successFunction func(testIndex int, ticket domain.Ticket, tagsForDevice []domain.Tag)
		errExpected     bool
	}{
		// testing Append option
		{
			&dal.Ticket{},
			&mockDevice{},
			func(_DeviceID string) (tags []domain.Tag, e error) {
				return []domain.Tag{&dal.Tag{
					DeviceIDvar: "",
					IDvar:       "",
					TagKeyIDvar: 0,
					Valuevar:    "",
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
					Optionsvar:           Append,
					TicketingSourceIDvar: "",
					TicketingTagvar:      "",
				},
			},
			func(testIndex int, ticket domain.Ticket, tagsForDevice []domain.Tag) {

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

		ticket, tagsForDevice, err := tj.handleCloudTagMappings(test.ticket, test.device)

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
