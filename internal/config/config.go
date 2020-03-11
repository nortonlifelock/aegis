package config

// AppConfig is a struct that implements the Config interface. It is used to parse the json contents of the app config
// when it is loaded by another package, it should be loaded into an interface such that only the necessary fields can
// be accessed
type AppConfig struct {
	DatabasePath     string `json:"db_path"`
	DatabasePort     string `json:"db_port"`
	DatabaseUser     string `json:"db_username"`
	DatabasePassword string `json:"db_password"`
	DatabaseSchema   string `json:"db_schema"`

	EKey       string `json:"key_id"`
	EType      string `json:"encryption_type"`
	TopicKey   string `json:"sns_id"`
	RegionKMS  string `json:"kms_region"`
	RegionSNS  string `json:"sns_region"`
	ProfileKMS string `json:"kms_profile"`

	// Logging
	LogFilePath   string `json:"logpath"`
	LogToFile     bool   `json:"log_to_file"`
	LogToConsole  bool   `json:"log_to_console"`
	LogToDb       bool   `json:"log_to_db"`
	LogToMessageQ bool   `json:"log_to_messageQ"`
	LogDontDelete bool   `json:"log_no_delete"`

	APIServicePort int  `json:"api_port,omitempty"`
	Debug          bool `json:"debug,omitempty"`

	SocketProtocol string `json:"websocket_protocol,omitempty"`
	Protocol       string `json:"transport_protocol,omitempty"`
	UI             string `json:"ui_location,omitempty"`

	PathToAegis string `json:"path_to_aegis"`
}

// Validate implements the IValidator interface and defines whether the app config is considered valid or not
// it contains the bare essentials - the information required to connect to the database
func (cfg AppConfig) Validate() (valid bool) {

	if len(cfg.DatabasePath) > 0 {
		if len(cfg.DatabasePort) > 0 {
			if len(cfg.DatabaseUser) > 0 {
				if len(cfg.EKey) > 0 {
					if len(cfg.PathToAegis) > 0 {
						valid = true
					}
				}
			}
		}
	}

	return valid
}

// AegisPath returns the path up to an including the Aegis directory. This is used over os.WorkingDir() as it is dependent on the location in which
// the process is started
func (cfg AppConfig) AegisPath() string {
	return cfg.PathToAegis
}

// DBPath returns the IP/URL of the database
func (cfg AppConfig) DBPath() string {
	return cfg.DatabasePath
}

// DBPort returns the port on which the database is connected to
func (cfg AppConfig) DBPort() string {
	return cfg.DatabasePort
}

// DBUsername returns the username to connect to the database
func (cfg AppConfig) DBUsername() string {
	return cfg.DatabaseUser
}

// DBPassword returns the unencrypted password to connect to the db. It should be encrypted in the app config file,
// but decrypted before loaded into the struct
func (cfg AppConfig) DBPassword() string {
	return cfg.DatabasePassword
}

// DBSchema returns the schema of the database used by the project. It tracks schema changes, and holds the data and stored
// procedures utilized by the project
func (cfg AppConfig) DBSchema() string {
	return cfg.DatabaseSchema
}

// EncryptionKey returns the AWS key used for KMS encryption/decryption
func (cfg AppConfig) EncryptionKey() string {
	return cfg.EKey
}

// LogPath returns the path to which the logs are stored
func (cfg AppConfig) LogPath() string {
	return cfg.LogFilePath
}

// LogFile returns a boolean which controls whether or not logs are stored in files
func (cfg AppConfig) LogFile() bool {
	return cfg.LogToFile
}

// LogConsole returns a boolean which controls whether or not logs are printed to the console
func (cfg AppConfig) LogConsole() bool {
	return cfg.LogToConsole
}

// LogDB returns a boolean which controls whether or not logs are stored in the database
func (cfg AppConfig) LogDB() bool {
	return cfg.LogToDb
}

// LogMQ returns a boolean which controls whether are not logs are pushed onto a message queue
func (cfg AppConfig) LogMQ() bool {
	return cfg.LogToMessageQ
}

// APIPort returns the port on which the API listens
func (cfg AppConfig) APIPort() int {
	return cfg.APIServicePort
}

// DebugLogs returns a boolean which controls whether debug logs are processed or not
func (cfg AppConfig) DebugLogs() bool {
	return cfg.Debug
}

// WebSocketProtocol returns which websocket protocol is used (ws/wss)
func (cfg AppConfig) WebSocketProtocol() string {
	return cfg.SocketProtocol
}

// TransportProtocol returns which transport protocol is used (http/https)
func (cfg AppConfig) TransportProtocol() string {
	return cfg.Protocol
}

// UILocation defines the UI URL so the API can confirm the origin of the request
func (cfg AppConfig) UILocation() string {
	return cfg.UI
}

// SNSTopicID returns the topic ID for an SNS
func (cfg AppConfig) SNSTopicID() string {
	return cfg.TopicKey
}

func (cfg AppConfig) EncryptionType() string {
	return cfg.EType
}

func (cfg AppConfig) KMSRegion() string {
	return cfg.RegionKMS
}

func (cfg AppConfig) KMSProfile() string {
	return cfg.ProfileKMS
}

func (cfg AppConfig) SNSRegion() string {
	return cfg.RegionSNS
}

func (cfg AppConfig) PreserveFileLogs() bool {
	return cfg.LogDontDelete
}
