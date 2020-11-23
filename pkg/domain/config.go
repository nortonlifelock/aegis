package domain

// Config holds all the configurations from within app.json
// to avoid giving the entire application access to the entire configuration, we use interfaces that implement a subset of the
// methods within the Config interface
type Config interface {
	DBConfig
	EncryptionConfig
	LogConfig
	WebAppConfig
}

// DBConfig defines an interface which contains methods for building connection strings
type DBConfig interface {
	DBPath() string
	DBPort() string
	DBUsername() string
	DBPassword() string
	DBSchema() string
}

// EncryptionConfig defines an interface which returns methods for encryption / decryption
type EncryptionConfig interface {
	EncryptionKey() string
	EncryptionType() string
	KMSRegion() string
	KMSProfile() string
}

// LogConfig defines an interface which returns the methods for logging configurations for the application
type LogConfig interface {
	LogPath() string
	LogFile() bool
	LogConsole() bool
	LogDB() bool
	LogMQ() bool
	DebugLogs() bool
	SNSTopicID() string
	SNSRegion() string
	PreserveFileLogs() bool
}

// WebAppConfig defines the interface for hosting the api and UI of the web application
type WebAppConfig interface {
	APIPort() int
	WebSocketProtocol() string
	TransportProtocol() string
	UILocation() string
	AegisPath() string
}
