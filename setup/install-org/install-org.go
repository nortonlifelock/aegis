package install_org

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/nortonlifelock/config"
	"github.com/nortonlifelock/crypto"
	"github.com/nortonlifelock/database"
	"github.com/nortonlifelock/domain"
	"github.com/nortonlifelock/implementations"
	"github.com/nortonlifelock/integrations"
	"github.com/nortonlifelock/jira"
	nexpose "github.com/nortonlifelock/nexpose/connector"
	qualys "github.com/nortonlifelock/qualys/connector"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	rescanJob      = "RescanJob"
	rescanQueueJob = "RescanQueueJob"
	ticketingJob   = "TicketingJob"
	exceptionJob   = "ExceptionJob"
	scanSyncJob    = "ScanSyncJob"
	scanCloseJob   = "ScanCloseJob"
	bulkUpdateJob  = "BulkUpdateJob"
	vulnSyncJob    = "VulnSyncJob"
	assetSyncJob   = "AssetSyncJob"
	ticketSyncJob  = "TicketSyncJob"
)

func InstallOrg(path string) {
	if len(path) == 0 {
		panic("must provide a -path flag leading to the application config (app.json)")
	}

	db, appConfig := loadDatabaseConnectionAppConfig(path)

	reader := bufio.NewReader(os.Stdin)

	organization := createOrganization(reader, appConfig, db)

	aesClient, err := crypto.NewEncryptionClient(crypto.AES256, db, appConfig.EncryptionKey(), organization.ID(), appConfig.KMSProfile(), appConfig.KMSRegion())
	check(err)

	fmt.Println("Now let's get your sources setup! We'll start with the essentials")
	scannerSC, ticketSC, assetGroups := createSources(reader, aesClient, db, organization.ID())

	fmt.Println("Creating job configs...")
	createJobConfigs(db, organization.ID(), scannerSC.ID(), ticketSC.ID(), assetGroups)
}

func createJobConfigs(db domain.DatabaseConnection, orgID, scannerSCID, ticketSCID string, assetGroups []int) {
	jobRegistrations, err := db.GetJobs()
	check(err)

	jobNameToJobID := make(map[string]int)
	for _, registration := range jobRegistrations {
		jobNameToJobID[registration.GoStruct()] = registration.ID()
	}

	normalRsqPayload := implementations.RescanQueuePayload{
		Type: domain.RescanNormal,
	}
	normalRsqBody, err := json.Marshal(&normalRsqPayload)
	check(err)
	decomRsqPayload := implementations.RescanQueuePayload{
		Type: domain.RescanDecommission,
	}
	decomRsqBody, err := json.Marshal(&decomRsqPayload)
	check(err)
	passiveRsqPayload := implementations.RescanQueuePayload{
		Type: domain.RescanPassive,
	}
	passiveRsqBody, err := json.Marshal(&passiveRsqPayload)
	check(err)
	exceptionRsqPayload := implementations.RescanQueuePayload{
		Type: domain.RescanExceptions,
	}
	exceptionRsqBody, err := json.Marshal(&exceptionRsqPayload)
	check(err)
	now := time.Now()
	ticketingPayload := implementations.TicketingPayload{
		MinDate: &now,
	}
	ticketingBody, err := json.Marshal(&ticketingPayload)
	check(err)
	assetSyncPayload := implementations.AssetSyncPayload{
		GroupIDs: assetGroups,
	}
	assetSyncBody, err := json.Marshal(&assetSyncPayload)
	check(err)

	_, _, err = db.CreateJobConfigWPayload(jobNameToJobID[rescanQueueJob], orgID, 0, true, 60, 1, true, "INSTALLER", ticketSCID, scannerSCID, string(normalRsqBody))
	check(err)
	_, _, err = db.CreateJobConfigWPayload(jobNameToJobID[rescanQueueJob], orgID, 0, true, 60, 1, true, "INSTALLER", ticketSCID, scannerSCID, string(decomRsqBody))
	check(err)
	_, _, err = db.CreateJobConfigWPayload(jobNameToJobID[rescanQueueJob], orgID, 0, false, 0, 1, false, "INSTALLER", ticketSCID, scannerSCID, string(passiveRsqBody))
	check(err)
	_, _, err = db.CreateJobConfigWPayload(jobNameToJobID[rescanQueueJob], orgID, 0, false, 0, 1, false, "INSTALLER", ticketSCID, scannerSCID, string(exceptionRsqBody))
	check(err)

	fmt.Println("Finished job config creation")

	_, _, err = db.CreateJobConfigWPayload(jobNameToJobID[rescanJob], orgID, 0, false, 0, 0, false, "INSTALLER", scannerSCID, ticketSCID, "{}")
	check(err)
	_, _, err = db.CreateJobConfigWPayload(jobNameToJobID[ticketingJob], orgID, 0, false, 0, 0, false, "INSTALLER", scannerSCID, ticketSCID, string(ticketingBody))
	check(err)
	_, _, err = db.CreateJobConfigWPayload(jobNameToJobID[scanSyncJob], orgID, 0, true, 120, 1, true, "INSTALLER", scannerSCID, scannerSCID, "{}")
	check(err)
	_, _, err = db.CreateJobConfigWPayload(jobNameToJobID[scanCloseJob], orgID, 0, false, 0, 0, false, "INSTALLER", scannerSCID, ticketSCID, "{}")
	check(err)
	_, _, err = db.CreateJobConfigWPayload(jobNameToJobID[bulkUpdateJob], orgID, 0, false, 0, 0, false, "INSTALLER", ticketSCID, ticketSCID, "{}")
	check(err)
	_, _, err = db.CreateJobConfigWPayload(jobNameToJobID[vulnSyncJob], orgID, 0, false, 0, 0, false, "INSTALLER", scannerSCID, scannerSCID, "{}")
	check(err)
	_, _, err = db.CreateJobConfigWPayload(jobNameToJobID[assetSyncJob], orgID, 0, false, 0, 0, false, "INSTALLER", scannerSCID, scannerSCID, string(assetSyncBody))
	check(err)
	_, _, err = db.CreateJobConfigWPayload(jobNameToJobID[ticketSyncJob], orgID, 0, false, 0, 1, false, "INSTALLER", ticketSCID, ticketSCID, "{}")
	check(err)

	// TODO set AutoStart to true once CERFs are handled in the DB (in-case they don't have a CERF project setup)
	_, _, err = db.CreateJobConfigWPayload(jobNameToJobID[exceptionJob], orgID, 0, true, 60, 1, false, "INSTALLER", ticketSCID, scannerSCID, "{}")
	check(err)
}

func createSources(reader *bufio.Reader, aesClient crypto.Client, db domain.DatabaseConnection, orgID string) (scannerSC, ticketSC domain.SourceConfig, assetGroups []int) {
	sources, err := db.GetSources()
	check(err)

	var sourceNameToID = make(map[string]string)
	for _, source := range sources {
		sourceNameToID[source.Source()] = source.ID()
	}

	scannerSC, assetGroups = createScannerSourceConfig(reader, aesClient, db, orgID, sourceNameToID)
	ticketSC = createTicketSourceConfig(reader, aesClient, db, orgID, sourceNameToID)

	_, _, err = db.UpdateSourceConfigConcurrencyByID(scannerSC.ID(), 10, 0, 10)
	check(err)

	_, _, err = db.UpdateSourceConfigConcurrencyByID(ticketSC.ID(), 10, 0, 10)
	check(err)

	return scannerSC, ticketSC, assetGroups
}

func createTicketSourceConfig(reader *bufio.Reader, aesClient crypto.Client, db domain.DatabaseConnection, orgID string, sourceNameToID map[string]string) (ticketSourceConfig domain.SourceConfig) {
	for {
		fmt.Println("Setting up your JIRA source configurations...")

		fmt.Println("What's the address for the API of your ticketing instance? (do not include port)")
		address := getInput(reader)

		if strings.Index(address, "http://") < 0 && strings.Index(address, "https://") < 0 {
			address = fmt.Sprintf("https://%s", address)
		}

		fmt.Println("What port does the API listen on? (leave empty if NA)")
		port := getInput(reader)

		fmt.Println("What's the username of the service account that will be authenticating for the API?")
		user := getInput(reader)

		fmt.Println("What's the password of the service account that will be authenticating for the API?")
		passByte, err := terminal.ReadPassword(0)
		check(err)

		pass, err := aesClient.Encrypt(string(passByte))
		check(err)

		payload := jira.PayloadJira{
			Project:        "",
			MappableFields: []string{},
			StatusMap:      nil,
			FieldMap:       nil,
		}

		fmt.Println("What JIRA project are you using?")
		payload.Project = getInput(reader)

		// If your ticket status do not align with the statuses that Aegis uses for tracking, you can map your status to our status in the payload of the JIRA source config in the database
		// the Aegis recognized statuses are the keys, and your custom status is the value
		payload.StatusMap = map[string]string{
			"scan-error":              "Scan-Error",
			"open":                    "Open",
			"reopened":                "Reopened",
			"in-progress":             "In-Progress",
			"resolved-remediated":     "Resolved-Remediated",
			"resolved-falsepositive":  "Resolved-FalsePositive",
			"resolved-decommissioned": "Resolved-Decommissioned",
			"resolved-exception":      "Resolved-Exception",
			"closed-remediated":       "Closed-Remediated",
			"closed-false-positive":   "Closed-False-Positive",
			"closed-decommission":     "Closed-Decommission",
			"closed-exception":        "Closed-Exception",
			"closed-cerf":             "Closed-CERF",
			"closed-error":            "Closed-Error",
		}

		// If your ticket fields do not align with the fields that Aegis uses for tracking, you can map your status to our status in the payload of the JIRA source config in the database
		// the Aegis recognized fields are the keys, and your custom field is the value
		payload.FieldMap = map[string]string{
			"Method of Discovery":    "Method of Discovery",
			"Summary":                "Summary",
			"Hostname":               "Hostname",
			"IP Address":             "IP Address",
			"MAC Address":            "MAC Address",
			"Service Port":           "Service Port",
			"Description":            "Description",
			"Solution":               "Solution",
			"VRR Priority":           "VRR Priority",
			"Scan/Alert Date":        "Scan/Alert Date",
			"Assignment Group":       "Assignment Group",
			"Resolution Date":        "Resolution Date",
			"Operating System":       "Operating System",
			"Vulnerablility":         "Vulnerablility",
			"cvss_score":             "cvss_score",
			"VulnerabilityID":        "VulnerabilityID",
			"GroupID":                "GroupID",
			"DeviceID":               "DeviceID",
			"ScanID":                 "ScanID",
			"Org":                    "Org",
			"CERF":                   "CERF",
			"Actual Expiration Date": "Actual Expiration Date",
			"CVE References":         "CVE References",
			"VendorRef":              "VendorRef",
			"OS_Detailed":            "OS_Detailed",
			"Config":                 "Config",
			"LastChecked":            "LastChecked",
			"CloudID":                "CloudID",
		}

		body, err := json.Marshal(&payload)

		source := integrations.JIRA
		fmt.Println("Finished! Creating source config")
		_, _, err = db.CreateSourceConfig(source, sourceNameToID[source], orgID, address, port, user, pass, "", "", "", string(body))
		check(err)

		sourceConfigs, err := db.GetSourceConfigBySourceID(orgID, sourceNameToID[source])
		check(err)

		if len(sourceConfigs) > 0 {
			ticketSourceConfig = sourceConfigs[0]
			break
		} else {
			fmt.Println("Source config could not be found in database")
		}
	}

	return ticketSourceConfig
}

func createScannerSourceConfig(reader *bufio.Reader, aesClient crypto.Client, db domain.DatabaseConnection, orgID string, sourceNameToID map[string]string) (scannerSourceConfig domain.SourceConfig, assetGroups []int) {
	for {
		fmt.Println("What vulnerability scanner do you plan on using, Nexpose or Qualys")

		vulnScanner := getInput(reader)
		if len(vulnScanner) > 0 {
			vulnScanner = strings.ToLower(vulnScanner)

			var nex = vulnScanner[0] == 'n'
			var qual = vulnScanner[0] == 'q'

			if nex || qual {
				var source string
				if nex {
					source = "Nexpose"
				} else {
					source = "Qualys"
				}

				fmt.Println("What's the address for the API of your vulnerability scanner instance? (do not include port)")
				address := getInput(reader)

				if strings.Index(address, "http://") < 0 && strings.Index(address, "https://") < 0 {
					address = fmt.Sprintf("https://%s", address)
				}

				fmt.Println("What port does the API listen on? (leave empty if NA)")
				port := getInput(reader)

				fmt.Println("What's the username of the service account that will be authenticating for the API?")
				user := getInput(reader)

				fmt.Println("What's the password of the service account that will be authenticating for the API?")
				passByte, err := terminal.ReadPassword(0)
				check(err)

				pass, err := aesClient.Encrypt(string(passByte))
				check(err)

				var body []byte
				if nex {
					payload := nexpose.Payload{
						ScanTemplate:        "",
						ScanNameFormat:      "",
						RescanSite:          0,
						DiscoveryTemplate:   "",
						DiscoveryNameFormat: "",
						EngineCacheTTL:      nil,
					}

					fmt.Println("Enter the ID of the scan template that will be used for vulnerability scans")
					payload.ScanTemplate = getInput(reader)

					fmt.Println("Enter the ID of the scan template that will be used for discovery scans")
					payload.DiscoveryTemplate = getInput(reader)

					fmt.Println("Enter the name of the scan format you would like your vulnerability scans to have (will have the creation date appended by automation)")
					payload.ScanNameFormat = getInput(reader)

					fmt.Println("Enter the name of the scan format you would like your discovery scans to have (will have the creation date appended by automation)")
					payload.DiscoveryNameFormat = getInput(reader)

					payload.RescanSite = getInputInt(reader, "Enter the site ID you will be using for rescans")

					assetGroups = getInputIntArray(reader, "Please enter all integer IDs for the asset groups you plan on doing vulnerability scans on (comma separated)")
					body, err = json.Marshal(&payload)
				} else {
					payload := qualys.QSPayload{
						SearchListID:              0,
						OptionProfileID:           0,
						DiscoveryOptionProfileID:  0,
						AssetGroups:               nil,
						KernelFilter:              0,
						ScanNameFormatString:      "",
						OptionProfileFormatString: "",
						SearchListFormatString:    "",
						ExternalGroups:            nil,
					}

					fmt.Println("Enter the name of the scan format you would like your vulnerability scans to have (must contain a single '%s' to be replaced with the scan creation date)")
					payload.ScanNameFormatString = getInput(reader)

					fmt.Println("Enter the name format you would like your Option Profiles to have (must contain a single '%s' to be replaced with the scan creation date)")
					payload.OptionProfileFormatString = getInput(reader)

					fmt.Println("Enter the name format you would like your Search Lists to have (must contain a single '%s' to be replaced with the scan creation date)")
					payload.SearchListFormatString = getInput(reader)

					payload.SearchListID = getInputInt(reader, "Enter the integer search list ID vulnerability scans")
					payload.OptionProfileID = getInputInt(reader, "Enter the integer option profile ID you will be using for vulnerability scans")
					payload.DiscoveryOptionProfileID = getInputInt(reader, "Enter the integer option profile ID you will be using for discovery scans")

					payload.AssetGroups = getInputIntArray(reader, "Please enter all integer IDs for the asset groups you plan on doing vulnerability scans on (comma separated)")
					payload.ExternalGroups = getInputIntArray(reader, "Please enter all integer IDs for the asset groups that will require external scanners (comma separated)")
					assetGroups = payload.AssetGroups
					body, err = json.Marshal(&payload)
				}
				check(err)

				fmt.Println("Finished! Creating source config")
				_, _, err = db.CreateSourceConfig(source, sourceNameToID[source], orgID, address, port, user, pass, "", "", "", string(body))
				check(err)

				sourceConfigs, err := db.GetSourceConfigBySourceID(orgID, sourceNameToID[source])
				check(err)

				if len(sourceConfigs) > 0 {
					scannerSourceConfig = sourceConfigs[0]
					break
				} else {
					fmt.Println("Source config could not be found in database")
				}
			} else {
				fmt.Println("Must choose either Nexpose or Qualys")
			}
		} else {
			fmt.Println("Empty input")
		}
	}

	return scannerSourceConfig, assetGroups
}

func createOrganization(reader *bufio.Reader, appConfig config.AppConfig, db domain.DatabaseConnection) (organization domain.Organization) {
	fmt.Print("Enter the name of your organization: ")
	name := getInput(reader)

	fmt.Print("Enter a shorthand code for the name of your organization: ")
	code := getInput(reader)
	code = strings.ToUpper(code)

	fmt.Println("Encrypting Organization encryption key using KMS...")
	kmsClient, err := crypto.CreateKMSClientWithProfile(appConfig.EncryptionKey(), appConfig.KMSProfile(), appConfig.KMSRegion())
	check(err)
	fmt.Println("Done")

	encryptedOrganizationKey, err := kmsClient.Encrypt(generateSecureEncryptionKey())
	check(err)

	const defaultPayload = `{"lowest_ticketed_cvss":4,"cvss_version":2,"severities":[{"name":"Medium","duration":90,"cvss_min":4},{"name":"High","duration":60,"cvss_min":7},{"name":"Critical","duration":30,"cvss_min":9}],"ad_servers":[""],"ad_ldap_tls_port":636,"ad_base_dn":"","ad_skip_verify":false,"ad_member_of_attribute":"memberOf","ad_search_string":"(accountName=%s)"}`

	_, _, err = db.CreateOrganizationWithPayloadEkey(code, name, 0, defaultPayload, encryptedOrganizationKey, "initializer")
	check(err)

	fmt.Printf("Organization [%s] created!\nYou can manage how the organization handles CVSS severity during ticketing, or the Active Directory configurations in the Payload column of the Organization table!\n", code)
	organization, err = db.GetOrganizationByCode(code)
	check(err)

	return organization
}

func generateSecureEncryptionKey() (ekey string) {
	fmt.Println("Generating secure Organization encryption key using crypto/rand")

	b := make([]byte, 32)
	_, err := rand.Read(b)
	check(err)
	ekey = base64.StdEncoding.EncodeToString(b)
	if len(ekey) > 32 {
		ekey = ekey[:32]
	}

	return ekey
}

func loadDatabaseConnectionAppConfig(path string) (db domain.DatabaseConnection, appConfig config.AppConfig) {
	var err error
	appConfig, err = config.LoadConfigByPath(path)
	check(err)
	db, err = database.NewConnection(appConfig)
	check(err)

	return db, appConfig
}

func getInput(reader *bufio.Reader) (userInput string) {
	var err error
	userInput, err = reader.ReadString('\n')
	check(err)

	userInput = strings.TrimSuffix(userInput, "\n")
	return userInput
}

func getInputInt(reader *bufio.Reader, message string) (userInput int) {
	for {
		fmt.Println(message)
		idString := getInput(reader)
		if idInt, err := strconv.Atoi(idString); err == nil {
			userInput = idInt
			break
		} else {
			fmt.Println("Please enter an integer")
		}
	}

	return userInput
}

func getInputIntArray(reader *bufio.Reader, message string) (userInput []int) {
	for {
		fmt.Println(message)
		userList := getInput(reader)
		userList = strings.Replace(userList, " ", "", -1)

		var noElementsThatWerentInts = true
		idList := strings.Split(userList, ",")
		userInput = make([]int, 0)
		for _, id := range idList {
			if idInt, err := strconv.Atoi(id); err == nil {
				userInput = append(userInput, idInt)
			} else {
				fmt.Printf("List included [%s] which is not an integer", id)
				noElementsThatWerentInts = false
				break
			}
		}

		if noElementsThatWerentInts {
			break
		}
	}

	return userInput
}

func check(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
