package qualys

import "encoding/xml"

// OptionProfiles holds the configurable options that Qualys uses during scanning. In the application, a template OptionProfile should be stored in the
// Qualys API. The application grabs the template OptionProfile, overwrites the search list, and then creates the new option profile to be used in a rescan.
// The search list only contains the QIDs that we wants Qualys to scan for. The template search list and option profile are deleted after the scan completes
type OptionProfiles struct {
	XMLName       xml.Name `xml:"OPTION_PROFILES"`
	Text          string   `xml:",chardata"`
	OptionProfile struct {
		Text      string `xml:",chardata"`
		BasicInfo struct {
			Text              string `xml:",chardata"`
			ID                string `xml:"ID"`
			GroupName         string `xml:"GROUP_NAME"`
			GroupType         string `xml:"GROUP_TYPE"`
			UserID            string `xml:"USER_ID"`
			UnitID            string `xml:"UNIT_ID"`
			SubscriptionID    string `xml:"SUBSCRIPTION_ID"`
			IsDefault         string `xml:"IS_DEFAULT"`
			IsGlobal          string `xml:"IS_GLOBAL"`
			IsOfflineSyncable string `xml:"IS_OFFLINE_SYNCABLE"`
			UpdateDate        string `xml:"UPDATE_DATE"`
		} `xml:"BASIC_INFO"`
		Scan struct {
			Text  string `xml:",chardata"`
			Ports struct {
				Text     string `xml:",chardata"`
				TCPPorts struct {
					Text               string `xml:",chardata"`
					TCPPortsType       string `xml:"TCP_PORTS_TYPE"`
					TCPPortsAdditional struct {
						Text            string `xml:",chardata"`
						HasAdditional   string `xml:"HAS_ADDITIONAL"`
						AdditionalPorts string `xml:"ADDITIONAL_PORTS"`
					} `xml:"TCP_PORTS_ADDITIONAL"`
					ThreeWayHandshake string `xml:"THREE_WAY_HANDSHAKE"`
				} `xml:"TCP_PORTS"`
				UDPPorts struct {
					Text               string `xml:",chardata"`
					UDPPortsType       string `xml:"UDP_PORTS_TYPE"`
					UDPPortsAdditional struct {
						Text            string `xml:",chardata"`
						HasAdditional   string `xml:"HAS_ADDITIONAL"`
						AdditionalPorts string `xml:"ADDITIONAL_PORTS"`
					} `xml:"UDP_PORTS_ADDITIONAL"`
				} `xml:"UDP_PORTS"`
				AuthoritativeOption string `xml:"AUTHORITATIVE_OPTION"`
			} `xml:"PORTS"`
			ScanDeadHosts         string `xml:"SCAN_DEAD_HOSTS"`
			PurgeOldHostOSChanged string `xml:"PURGE_OLD_HOST_OS_CHANGED"`
			Performance           struct {
				Text               string `xml:",chardata"`
				ParallelScaling    string `xml:"PARALLEL_SCALING"`
				OverallPerformance string `xml:"OVERALL_PERFORMANCE"`
				HostsToScan        struct {
					Text              string `xml:",chardata"`
					ExternalScanners  string `xml:"EXTERNAL_SCANNERS"`
					ScannerAppliances string `xml:"SCANNER_APPLIANCES"`
				} `xml:"HOSTS_TO_SCAN"`
				ProcessesToRun struct {
					Text           string `xml:",chardata"`
					TotalProcesses string `xml:"TOTAL_PROCESSES"`
					HTTPProcesses  string `xml:"HTTP_PROCESSES"`
				} `xml:"PROCESSES_TO_RUN"`
				PacketDelay                  string `xml:"PACKET_DELAY"`
				PortScanningAndHostDiscovery string `xml:"PORT_SCANNING_AND_HOST_DISCOVERY"`
			} `xml:"PERFORMANCE"`
			LoadBalancerDetection  string `xml:"LOAD_BALANCER_DETECTION"`
			VulnerabilityDetection struct {
				Text       string `xml:",chardata"`
				CustomList struct {
					Text   string            `xml:",chardata"`
					Custom []SearchListEntry `xml:"CUSTOM"`
				} `xml:"CUSTOM_LIST"`
				DetectionInclude struct {
					Text                string `xml:",chardata"`
					BasicHostInfoChecks string `xml:"BASIC_HOST_INFO_CHECKS"`
					OvalChecks          string `xml:"OVAL_CHECKS"`
				} `xml:"DETECTION_INCLUDE"`
			} `xml:"VULNERABILITY_DETECTION"`
			Authentication    string `xml:"AUTHENTICATION"`
			ADDLCertDetection string `xml:"ADDL_CERT_DETECTION"`
			DissolvableAgent  struct {
				Text                          string `xml:",chardata"`
				DissolvableAgentEnable        string `xml:"DISSOLVABLE_AGENT_ENABLE"`
				WindowsShareEnumerationEnable string `xml:"WINDOWS_SHARE_ENUMERATION_ENABLE"`
			} `xml:"DISSOLVABLE_AGENT"`
			HostAliveTesting string `xml:"HOST_ALIVE_TESTING"`
		} `xml:"SCAN"`
		MAP struct {
			Text                 string `xml:",chardata"`
			BasicInfoGatheringOn string `xml:"BASIC_INFO_GATHERING_ON"`
			TCPPorts             struct {
				Text                 string `xml:",chardata"`
				TCPPortsStandardScan string `xml:"TCP_PORTS_STANDARD_SCAN"`
			} `xml:"TCP_PORTS"`
			UDPPorts struct {
				Text                 string `xml:",chardata"`
				UDPPortsStandardScan string `xml:"UDP_PORTS_STANDARD_SCAN"`
			} `xml:"UDP_PORTS"`
			MapOptions struct {
				Text                 string `xml:",chardata"`
				PerformLiveHostSweep string `xml:"PERFORM_LIVE_HOST_SWEEP"`
				DisableDNSTraffic    string `xml:"DISABLE_DNS_TRAFFIC"`
			} `xml:"MAP_OPTIONS"`
			MapPerformance struct {
				Text               string `xml:",chardata"`
				OverallPerformance string `xml:"OVERALL_PERFORMANCE"`
				MapParallel        struct {
					Text              string `xml:",chardata"`
					ExternalScanners  string `xml:"EXTERNAL_SCANNERS"`
					ScannerAppliances string `xml:"SCANNER_APPLIANCES"`
					NetblockSize      string `xml:"NETBLOCK_SIZE"`
				} `xml:"MAP_PARALLEL"`
				PacketDelay string `xml:"PACKET_DELAY"`
			} `xml:"MAP_PERFORMANCE"`
			MapAuthentication string `xml:"MAP_AUTHENTICATION"`
		} `xml:"MAP"`
		Additional struct {
			Text          string `xml:",chardata"`
			HostDiscovery struct {
				Text      string `xml:",chardata"`
				TCPPPorts struct {
					Text         string `xml:",chardata"`
					StandardScan string `xml:"STANDARD_SCAN"`
				} `xml:"TCP_PORTS"`
				UDPPorts struct {
					Text         string `xml:",chardata"`
					StandardScan string `xml:"STANDARD_SCAN"`
				} `xml:"UDP_PORTS"`
				ICMP string `xml:"ICMP"`
			} `xml:"HOST_DISCOVERY"`
			PacketOptions struct {
				Text                                     string `xml:",chardata"`
				IgnoreFirewallGeneratedTCPRST            string `xml:"IGNORE_FIREWALL_GENERATED_TCP_RST"`
				IgnoreALLTCPRST                          string `xml:"IGNORE_ALL_TCP_RST"`
				IgnoreFirewallGeneratedTCPSYNCACK        string `xml:"IGNORE_FIREWALL_GENERATED_TCP_SYN_ACK"`
				NotSendTCPACKOrSYNACKDuringHostDiscovery string `xml:"NOT_SEND_TCP_ACK_OR_SYN_ACK_DURING_HOST_DISCOVERY"`
			} `xml:"PACKET_OPTIONS"`
		} `xml:"ADDITIONAL"`
	} `xml:"OPTION_PROFILE"`
}

// SearchListEntry is a member of OptionProfiles and must be exported in order to be marshalled
type SearchListEntry struct {
	Text  string `xml:",chardata"`
	ID    string `xml:"ID"`
	Title string `xml:"TITLE"`
}
