package log

import (
	"fmt"
	"time"
)

const (
	// InfoSeverity is the severity for non-serious information
	InfoSeverity = iota

	// WarnSeverity is for concerning information that does not stop the flow of the application but is concerning
	WarnSeverity

	// ErrorSeverity is for events that prevent normal flow of the application
	ErrorSeverity

	// CritSeverity is for series issues that occur in the application
	CritSeverity

	// FatalSeverity is for serious events that occur during execution
	FatalSeverity

	// InfoDetail is used in logs to discern that it is informational
	InfoDetail = "INFO"

	// WarnDetail is used in logs to discern that it is a warning
	WarnDetail = "WARN"

	// ErrorDetail is used in logs to discern that it is an error
	ErrorDetail = "ERROR"

	// CriticalDetail is used in logs to discern that it is a critical error
	CriticalDetail = "CRITICAL"

	// FatalDetail is used in logs to discern that it is a fatal error
	FatalDetail = "FATAL"

	// InfoColor defines the color for info logs
	InfoColor = "\033[0;32m"

	// WarnColor defines the color for warning logs
	WarnColor = "\033[1;33m"

	// ErrorColor defines the color for error logs
	ErrorColor = "\033[0;31m"

	// CriticalColor defines the color for critical logs
	CriticalColor = "\033[1;31m"

	// FatalColor defines the color for fatal logs
	FatalColor = "\033[0;35m"

	logTimestampFormat = "2006-01-02 15:04:05"
)

// Logger is the interface that defines the required method for processing a log
type Logger interface {
	Send(log Log)
}

// Log is a struct that holds all the information required to log an event
type Log struct {
	Type          int       `json:"type"`
	Text          string    `json:"text"`
	TimeStamp     time.Time `json:"timestamp"`
	Error         error     `json:"error"`
	StackTrace    string    `json:"stackTrace"`
	SpecialAppend string    `json:"special"`

	JobID        string `json:"jobId"`
	Job          string `json:"job"`
	OrgCode      string `json:"org_code"`
	JobInSource  string `json:"inSource"`
	JobOutSource string `json:"outSource"`

	IsDebug bool `json:"isDebug"`
	Slack   bool `json:"slack"`
	SNS     bool `json:"send_sns"`
}

func newLog() Log {
	return Log{
		TimeStamp: time.Now().UTC(),
	}
}

// Special sets the SpecialAppend field which is used when extra information needs to be included in the log
func (log Log) Special(special string) Log {
	if len(special) > 0 {
		log.SpecialAppend = fmt.Sprintf("%s", special)
	}

	return log
}

// Stack sets the stacktrace for the log
func (log Log) Stack(stacktrace string) Log {
	log.StackTrace = stacktrace

	return log
}

// Debug turns on the debug flag in the log and returns the log so that the methods may be chained
func (log Log) Debug() Log {
	log.IsDebug = true

	return log
}

func (log Log) SendSNS() Log {
	log.SNS = true
	return log
}

// ToConsoleString returns the string formatted for the console
func (log Log) ToConsoleString() (logString string) {

	logString = fmt.Sprintf("%s | [%s%s\033[0m] - %s", log.TimeStamp.Format(logTimestampFormat), log.TypeToColorCode(log.Type), log.TypeToString(log.Type), log.Text)

	if len(log.JobID) > 0 { //Is a Job thread, include job information
		logString = fmt.Sprintf("%s | [%s%s\033[0m] [%s:%s:%s] - %s", log.TimeStamp.Format(logTimestampFormat), log.TypeToColorCode(log.Type), log.TypeToString(log.Type), log.Job, log.OrgCode, log.JobID, log.Text)
	}

	if log.Error != nil {
		if len(log.Error.Error()) > 0 {
			logString = fmt.Sprintf("%s - [%s]", logString, log.Error.Error())
		}
	}

	return logString
}

// ToString returns a string-formatted version of the log
func (log Log) ToString() (logString string) {

	logString = fmt.Sprintf("%s | %s - %s", log.TimeStamp.Format(logTimestampFormat), log.TypeToString(log.Type), log.Text)

	//logContents = fmt.Sprintf("SEVERITY: %s - %s", logType, logMessage)
	if len(log.JobID) > 0 { //Is a Job thread, include job information
		logString = fmt.Sprintf("%s | %s %s [%s:%s] - %s", log.TimeStamp.Format(logTimestampFormat), log.TypeToString(log.Type), log.Job, log.OrgCode, log.JobID, log.Text)
	}

	return logString
}

// TypeToString returns a readable version of the log type
func (log Log) TypeToString(Type int) (TypeString string) {

	switch Type {
	case InfoSeverity:
		TypeString = InfoDetail
		break
	case WarnSeverity:
		TypeString = WarnDetail
		break
	case ErrorSeverity:
		TypeString = ErrorDetail
		break
	case CritSeverity:
		TypeString = CriticalDetail
		break
	case FatalSeverity:
		TypeString = FatalDetail
		break
	}

	return TypeString
}

// TypeToColorCode returns the characters that change the terminal output color according to the log type
func (log Log) TypeToColorCode(Type int) (color string) {

	switch Type {
	case InfoSeverity:
		color = InfoColor
		break
	case WarnSeverity:
		color = WarnColor
		break
	case ErrorSeverity:
		color = ErrorColor
		break
	case CritSeverity:
		color = CriticalColor
		break
	case FatalSeverity:
		color = FatalColor
		break
	}

	return color
}
