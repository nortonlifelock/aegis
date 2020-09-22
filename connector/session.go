package connector

import (
	"context"
	"encoding/json"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/log"
	"github.com/nortonlifelock/qualys"
	"sync"
)

type logger interface {
	Send(log log.Log)
}

// QSPayload is the Qualys Source Config Payload that is stored in the database for Qualys configurations
type QSPayload struct {
	// The default search list to be used for rescans
	SearchListID int `json:"searchListId"`

	// The default option profile to be used for rescans
	OptionProfileID int `json:"optionProfileId"`

	// The default option profile to be used for rescans
	DiscoveryOptionProfileID int `json:"discoveryOptionProfileId"`

	// The known asset groups for the specific implementation of scanning appliances for use in rescans
	AssetGroups []int `json:"groups"`

	// KernelFilter sets the arf_kernel_filter flag in the host detection API calls. Can hold values [0,4]
	// 0 vulnerabilities are not filtered based on kernel activity

	// TODO move this field to the organization level?
	// 1 exclude kernel related vulnerabilities that are not exploitable (found on non-running kernels)
	// 2 only include kernel related vulnerabilities that are not exploitable (found on non-running kernels)
	// 3 only include kernel related vulnerabilities that are exploitable (found on running kernels)
	// 4 only include kernel related vulnerabilities
	KernelFilter int `json:"kernel_filter"`

	// RescanNameTemplate holds the format string (which should contain a single %s) to name the rescans. The %s is replaced with a timestamp by the application
	ScanNameFormatString string `json:"rescan_name_format"`

	// OptionProfileFormatString holds the format string (which should contain a single %s) to name the option profiles. The %s is replaced with a timestamp
	// by the application. The option profile is deleted after the scan is completed
	OptionProfileFormatString string `json:"option_profile_format"`

	// SearchListFormatString holds the format string (which should contain a single %s) to name the search lists. The %s is replaced with a timestamp
	// by the application. The search list is deleted after the scan is completed
	SearchListFormatString string `json:"search_list_format"`

	// ExternalGroups holds a list of asset groups that must be scanned with the external scanner
	ExternalGroups []int `json:"external_groups"`

	// WebAppOptionProfile holds the ID of the option profile that you'd like to use for web application scans (WAS - optional)
	WebAppOptionProfile string `json:"web_app_option_profile"`

	// EC2ScanSettings controls the parameters used to create the ec2 scans
	EC2ScanSettings map[string]*struct {
		ConnectorName string `json:"connector_name"`
		ScannerName   string `json:"scanner_name"`
	} `json:"ec2_scan_settings"`
}

// QsSession is the struct that is responsible for making Qualys API calls
type QsSession struct {
	apiSession *qualys.Session

	// The map of the vulnerability knowledge base which will eventually be replaced by a database
	// for handling the vulnerabilities from Qualys
	vulnerabilities   map[int]*qualys.QVulnerability
	vulnerabilityLock *sync.Mutex

	lstream logger

	// The Qualys payload which came in from the Source config
	payload *QSPayload

	// Cache of the appliances loaded for the system
	appliances map[int][]int

	// Cache of asset groups (corresponding to the asset group slice in the QSPayload)
	assetGroupCache []*qualys.QSAssetGroup
}

// Connect returns a QsSession, which is used to process information returned from the Qualys API
func Connect(ctx context.Context, lstream logger, sourceConfig domain.SourceConfig) (session *QsSession, err error) {
	session = &QsSession{
		lstream:           lstream,
		vulnerabilities:   make(map[int]*qualys.QVulnerability),
		vulnerabilityLock: &sync.Mutex{},
		appliances:        make(map[int][]int),
	}

	var payload = &QSPayload{}
	if err = json.Unmarshal([]byte(sord(sourceConfig.Payload())), payload); err == nil {
		session.payload = payload
		session.apiSession, err = qualys.NewQualysAPISession(ctx, lstream, sourceConfig)
	}

	return session, err
}
