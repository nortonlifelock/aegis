package nexpose

// ScanTemplates is the json response for the list of scan templates in nexpose
type ScanTemplates struct {
	// Hypermedia links to corresponding or related resources.
	Links []Link `json:"links,omitempty"`

	// The list of can templates
	Templates []*ScanTemplate `json:"resources,omitempty"`
}

// ScanTemplate is the singular json response for scan templates from nexpose
type ScanTemplate struct {

	// Settings for which vulnerability checks to run during a scan.
	// The rules for inclusion of checks is as follows:
	// - Enabled checks by category and by check type are included
	// - Disabled checks in by category and by check type are removed
	// - Enabled checks in by individual check are added (even if they are disabled in by category or check type)
	// - Disabled checks in by individual check are removed
	// - If unsafe is disabled, unsafe checks are removed
	// - If potential is disabled, potential checks are removed
	Checks *struct {

		// The vulnerability check categories enabled or disabled during a scan.
		Categories *struct {

			// The categories of vulnerability checks to disable during a scan.
			Disabled []string `json:"disabled,omitempty"`

			// The categories of vulnerability checks to enable during a scan.
			Enabled []string `json:"enabled,omitempty"`

			// Hypermedia links to corresponding or related resources.
			Links []Link `json:"links,omitempty"`
		} `json:"categories,omitempty"`

		// Whether an extra step is performed at the end of the scan where more trust is put in OS patch checks to
		// attempt to override the results of other checks which could be less reliable.
		Correlate bool `json:"correlate,omitempty"`

		// The individual vulnerability checks enabled or disabled during a scan.
		Individual *struct {

			// The individual vulnerability checks to disable during a scan.
			Disabled []string `json:"disabled,omitempty"`

			// The individual vulnerability checks to enable during a scan.
			Enabled []string `json:"enabled,omitempty"`

			// Hypermedia links to corresponding or related resources.
			Links []Link `json:"links,omitempty"`
		} `json:"individual,omitempty"`

		// Hypermedia links to corresponding or related resources.
		Links []Link `json:"links,omitempty"`

		// Whether checks that result in potential vulnerabilities are assessed during a scan.
		Potential bool `json:"potential,omitempty"`

		// The vulnerability check types enabled or disabled during a scan.
		Types *struct {

			// The types of vulnerability checks to disable during a scan.
			Disabled []string `json:"disabled,omitempty"`

			// The types of vulnerability checks to enable during a scan.
			Enabled []string `json:"enabled,omitempty"`

			// Hypermedia links to corresponding or related resources.
			Links []Link `json:"links,omitempty"`
		} `json:"types,omitempty"`

		// Whether checks considered \"unsafe\" are assessed during a scan.
		Unsafe bool `json:"unsafe,omitempty"`
	} `json:"checks,omitempty"`

	// Settings for discovery databases.
	Database *struct {

		// Database name for DB2 database instance.
		Db2 string `json:"db2,omitempty"`

		// Hypermedia links to corresponding or related resources.
		Links []Link `json:"links,omitempty"`

		// Database name (SID) for an Oracle database instance.
		Oracle []string `json:"oracle,omitempty"`

		// Database name for PostgesSQL database instance.
		Postgres string `json:"postgres,omitempty"`
	} `json:"database,omitempty"`

	// A verbose description of the scan template..
	Description string `json:"description,omitempty"`

	// Discovery settings used during a scan.
	Discovery *struct {

		// Asset discovery settings used during a scan.
		Asset *struct {

			// Whether to query Whois during discovery.
			// Defaults to `false`.
			CollectWhoisInformation bool `json:"collectWhoisInformation,omitempty"`

			// The minimum certainty required for a fingerprint to be considered valid during a scan.
			// Defaults to `0.16`.
			FingerprintMinimumCertainty float32 `json:"fingerprintMinimumCertainty,omitempty"`

			// The number of fingerprinting attempts made to determine the operating system fingerprint.
			// Defaults to `4`.
			FingerprintRetries int32 `json:"fingerprintRetries,omitempty"`

			// Whether to fingerprint TCP/IP stacks for hardware, operating system and software information.
			IPFingerprintingEnabled bool `json:"ipFingerprintingEnabled,omitempty"`

			// Whether ARP pings are sent during asset discovery.
			// Defaults to `true`.
			SendArpPings bool `json:"sendArpPings,omitempty"`

			// Whether ICMP pings are sent during asset discovery.
			// Defaults to `false`.
			SendIcmpPings bool `json:"sendIcmpPings,omitempty"`

			// TCP ports to send packets and perform discovery.
			// Defaults to no ports.
			TCPPorts []int32 `json:"tcpPorts,omitempty"`

			// Whether TCP reset responses are treated as live assets.
			// Defaults to `true`.
			TreatTCPResetAsAsset bool `json:"treatTcpResetAsAsset,omitempty"`

			// UDP ports to send packets and perform discovery.
			// Defaults to no ports.
			UDPPorts []int32 `json:"udpPorts,omitempty"`
		} `json:"asset,omitempty"`

		// Discovery performance settings used during a scan.
		Performance *struct {

			// The number of packets to send per second during scanning.
			PacketRate *struct {

				// Whether defeat rate limit (defeat-rst-ratelimit) is enforced on the minimum packet setting, which
				// can improve scan speed. If it is disabled, the minimum packet rate setting may be ignored when a
				// target limits its rate of RST (reset) responses to a port scan. This can increase scan accuracy by
				// preventing the scan from missing ports.
				// Defaults to `true`.
				DefeatRateLimit bool `json:"defeatRateLimit,omitempty"`

				// The minimum number of packets to send each second during discovery attempts.
				// Defaults to `0`.
				Maximum int32 `json:"maximum,omitempty"`

				// The minimum number of packets to send each second during discovery attempts.
				// Defaults to `0`.
				Minimum int32 `json:"minimum,omitempty"`
			} `json:"packetRate,omitempty"`

			// The number of discovery connection requests to be sent to target host simultaneously. These settings
			// has no effect if values have been set for `scanDelay`.
			Parallelism *struct {

				// The maximum number of discovery connection requests send in parallel.
				// Defaults to `0`.
				Maximum int32 `json:"maximum,omitempty"`

				// The minimum number of discovery connection requests send in parallel.
				// Defaults to `0`.
				Minimum int32 `json:"minimum,omitempty"`
			} `json:"parallelism,omitempty"`

			// The maximum number of attempts to contact target assets. If the limit is exceeded with no response,
			// the given asset is not scanned.
			// Defaults to `3`.
			RetryLimit int32 `json:"retryLimit,omitempty"`

			// The duration to wait between sending packets to each target host during a scan.
			ScanDelay *struct {

				// The minimum duration to wait between sending packets to each target host. The value is specified
				// as a ISO8601 duration and can range from `PT0S` (0ms) to `P30S` (30s).
				// Defaults to `PT0S`.
				Maximum string `json:"maximum,omitempty"`

				// The maximum duration to wait between sending packets to each target host. The value is specified
				// as a ISO8601 duration and can range from `PT0S` (0ms) to `P30S` (30s).
				// Defaults to `PT0S`.
				Minimum string `json:"minimum,omitempty"`
			} `json:"scanDelay,omitempty"`

			// The duration to wait between retry attempts.
			Timeout *struct {

				// The initial timeout to wait between retry attempts. The value is specified as a ISO8601 duration
				// and can range from `PT0.5S` (500ms) to `P30S` (30s).
				// Defaults to `PT0.5S`.
				Initial string `json:"initial,omitempty"`

				// The maximum time to wait between retries. The value is specified as a ISO8601 duration and can
				// range from `PT0.5S` (500ms) to `P30S` (30s).
				// Defaults to `PT3S`.
				Maximum string `json:"maximum,omitempty"`

				// The minimum time to wait between retries. The value is specified as a ISO8601 duration and can
				// range from `PT0.5S` (500ms) to `P30S` (30s).
				// Defaults to `PT0.5S`.
				Minimum string `json:"minimum,omitempty"`
			} `json:"timeout,omitempty"`
		} `json:"performance,omitempty"`

		// Service discovery settings used during a scan.
		Service *struct {

			// An optional file that lists each port and the service that commonly resides on it. If scans cannot
			// identify actual services on ports, service names will be derived from this file in scan results.
			// Defaults to empty.
			ServiceNameFile string `json:"serviceNameFile,omitempty"`

			// TCP service discovery settings.
			TCP *struct {

				// Additional TCP ports to scan. Individual ports can be specified as numbers or a string, but port
				// ranges must be strings (e.g. 7892-7898).
				// Defaults to empty.
				AdditionalPorts []interface{} `json:"additionalPorts,omitempty"`

				// TCP ports to exclude from scanning. Individual ports can be specified as numbers or a string,
				// but port ranges must be strings (e.g. 7892-7898).
				// Defaults to empty.
				ExcludedPorts []interface{} `json:"excludedPorts,omitempty"`

				// Hypermedia links to corresponding or related resources.
				Links []Link `json:"links,omitempty"`

				// The method of TCP discovery.
				// Defaults to `SYN`.
				Method string `json:"method,omitempty"`

				// The TCP ports to scan.
				// Defaults to `well-known`.
				Ports string `json:"ports,omitempty"`
			} `json:"tcp,omitempty"`

			// UDP service discovery settings.
			UDP *struct {

				// Additional UDP ports to scan. Individual ports can be specified as numbers or a string, but port
				// ranges must be strings (e.g. 7892-7898).
				// Defaults to empty.
				AdditionalPorts []interface{} `json:"additionalPorts,omitempty"`

				// UDP ports to exclude from scanning. Individual ports can be specified as numbers or a string, but
				// port ranges must be strings (e.g. 7892-7898).
				// Defaults to empty.
				ExcludedPorts []interface{} `json:"excludedPorts,omitempty"`

				// Hypermedia links to corresponding or related resources.
				Links []Link `json:"links,omitempty"`

				// The UDP ports to scan.
				// Defaults to `well-known`.
				Ports string `json:"ports,omitempty"`
			} `json:"udp,omitempty"`
		} `json:"service,omitempty"`
	} `json:"discovery,omitempty"`

	// Whether only discovery is performed during a scan.
	DiscoveryOnly bool `json:"discoveryOnly,omitempty"`

	// Whether Windows services are enabled during a scan. Windows services will be temporarily reconfigured when
	// this option is selected. Original settings will be restored after the scan completes, unless it is interrupted.
	EnableWindowsServices bool `json:"enableWindowsServices,omitempty"`

	// Whether enhanced logging is gathered during scanning. Collection of enhanced logs may greatly increase the
	// disk space used by a scan.
	EnhancedLogging bool `json:"enhancedLogging,omitempty"`

	// The identifier of the scan template
	ID string `json:"id,omitempty"`

	// Hypermedia links to corresponding or related resources.
	Links []Link `json:"links,omitempty"`

	// The maximum number of assets scanned simultaneously per scan engine during a scan.
	MaxParallelAssets int32 `json:"maxParallelAssets,omitempty"`

	// The maximum number of scan processes simultaneously allowed against each asset during a scan.
	MaxScanProcesses int32 `json:"maxScanProcesses,omitempty"`

	// A concise name for the scan template.
	Name string `json:"name,omitempty"`

	// Policy configuration settings used during a scan.
	Policy *struct {

		// The identifiers of the policies enabled to be checked during a scan. No policies are enabled by default.
		Enabled []int64 `json:"enabled,omitempty"`

		// Hypermedia links to corresponding or related resources.
		Links []Link `json:"links,omitempty"`

		// Whether recursive windows file searches are enabled, if your internal security practices require this
		// capability. Recursive file searches can increase scan times by several hours, depending on the number
		// of files and other factors, so this setting is disabled for Windows systems by default.
		// Defaults to `false`.
		RecursiveWindowsFSSearch bool `json:"recursiveWindowsFSSearch,omitempty"`

		// Whether Asset Reporting Format (ARF) results are stored. If you are required to submit reports of your
		// policy scan results to the U.S. government in ARF for SCAP certification, you will need to store SCAP
		// data so that it can be exported in this format. Note that stored SCAP data can accumulate rapidly,
		// which can have a significant impact on file storage.
		// Defaults to `false`.
		StoreSCAP bool `json:"storeSCAP,omitempty"`
	} `json:"policy,omitempty"`

	// Whether policy assessment is performed during a scan.
	PolicyEnabled bool `json:"policyEnabled,omitempty"`

	// Settings for interacting with the Telnet protocol.
	Telnet *struct {

		// The character set to use.
		CharacterSet string `json:"characterSet,omitempty"`

		// Regular expression to match a failed login response.
		FailedLoginRegex string `json:"failedLoginRegex,omitempty"`

		// Hypermedia links to corresponding or related resources.
		Links []Link `json:"links,omitempty"`

		// Regular expression to match a login response.
		LoginRegex string `json:"loginRegex,omitempty"`

		// Regular expression to match a password prompt.
		PasswordPromptRegex string `json:"passwordPromptRegex,omitempty"`

		// Regular expression to match a potential false negative login response.
		QuestionableLoginRegex string `json:"questionableLoginRegex,omitempty"`
	} `json:"telnet,omitempty"`

	// Whether vulnerability assessment is performed during a scan.
	VulnerabilityEnabled bool `json:"vulnerabilityEnabled,omitempty"`

	// Web spider settings used during a scan.
	Web *struct {

		// Whether scanning of multi-use devices, such as printers or print servers should be avoided.
		DontScanMultiUseDevices bool `json:"dontScanMultiUseDevices,omitempty"`

		// Whether query strings are using in URLs when web spidering. This causes the web spider to make many more
		// requests to the Web server. This will increase overall scan time and possibly affect the Web server's
		// performance for legitimate users.
		IncludeQueryStrings bool `json:"includeQueryStrings,omitempty"`

		// Paths to use when web spidering.
		Paths *struct {

			// Paths to bootstrap spidering with.
			Boostrap string `json:"boostrap,omitempty"`

			// Paths excluded from spidering.
			Excluded string `json:"excluded,omitempty"`

			// ${scan.template.web.spider.paths.robot.directives.description}
			HonorRobotDirectives bool `json:"honorRobotDirectives,omitempty"`
		} `json:"paths,omitempty"`

		// Patterns to match responses during web spidering.
		Patterns *struct {

			// A regular expression that is used to find sensitive content on a page.
			SensitiveContent string `json:"sensitiveContent,omitempty"`

			// A regular expression that is used to find fields that may contain sensitive input.
			// Defaults to (p|pass)(word|phrase|wd|code).
			SensitiveField string `json:"sensitiveField,omitempty"`
		} `json:"patterns,omitempty"`

		// Performance settings used during web spidering.
		Performance *struct {

			// The names of HTTP Daemons (HTTPd) to skip when spidering. For example, CUPS.
			HTTPDaemonsToSkip []string `json:"httpDaemonsToSkip,omitempty"`

			// The directory depth limit for web spidering. Limiting directory depth can save significant time,
			// especially with large sites. A value of `0` signifies unlimited directory traversal.
			// Defaults to `6`.
			MaximumDirectoryLevels int32 `json:"maximumDirectoryLevels,omitempty"`

			// The maximum number of unique host names that the spider may resolve. This function adds substantial
			// time to the spidering process, especially with large Web sites, because of frequent cross-link
			// checking involved. Defaults to `100`.
			MaximumForeignHosts int32 `json:"maximumForeignHosts,omitempty"`

			// The maximum depth of links to traverse when spidering. Defaults to `6`.
			MaximumLinkDepth int32 `json:"maximumLinkDepth,omitempty"`

			// The maximum the number of pages that are spidered. This is a time-saving measure for large sites.
			// Defaults to `3000`.
			MaximumPages int32 `json:"maximumPages,omitempty"`

			// The maximum the number of times to retry a request after a failure. A value of `0` means no retry
			// attempts are made.
			// Defaults to `2`.
			MaximumRetries int32 `json:"maximumRetries,omitempty"`

			// The maximum length of time to web spider. This limit prevents scans from taking longer than the
			// allotted scan schedule. A value of `PT0S` means no limit is applied.
			// The acceptable range is `PT1M` to `PT16666.6667H`.
			MaximumTime string `json:"maximumTime,omitempty"`

			// The duration to wait for a response from a target web server. The value is specified as a ISO8601
			// duration and can range from `PT0S` (0ms) to `P1H` (1 hour).
			// Defaults to `PT2M`.
			ResponseTimeout string `json:"responseTimeout,omitempty"`

			// The number of threads to use per web server being spidered.
			// Defaults to `3`.
			ThreadsPerServer int32 `json:"threadsPerServer,omitempty"`
		} `json:"performance,omitempty"`

		// Whether to determine if discovered logon forms accept commonly used user names or passwords. The process
		// may cause authentication services with certain security policies to lock out accounts with these credentials.
		TestCommonUsernamesAndPasswords bool `json:"testCommonUsernamesAndPasswords,omitempty"`

		// Whether to test for persistent cross-site scripting during a single scan. This test helps to reduce the
		// risk of dangerous attacks via malicious code stored on Web servers. Enabling it may increase
		// Web spider scan times.
		TestXSSInSingleScan bool `json:"testXssInSingleScan,omitempty"`

		// The `User-Agent` to use when web spidering.
		// Defaults to Mozilla/5.0 (compatible; MSIE 7.0; Windows NT 6.0; .NET CLR 1.1.4322; .NET CLR 2.0.50727).
		UserAgent string `json:"userAgent,omitempty"`
	} `json:"web,omitempty"`

	// Whether web spidering and assessment are performed during a scan.
	WebEnabled bool `json:"webEnabled,omitempty"`
}
