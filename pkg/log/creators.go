package log

import (
	"fmt"
)

// Info creates a log for non-serious information
func Info(text string) (log Log) {

	log = newLog()
	log.Type = InfoSeverity
	log.Text = text

	return log
}

// Infof creates a log for non-serious information allowing a formatting string
// to be passed for the log rather than a straight string parameter
func Infof(format string, a ...interface{}) (log Log) {
	return Info(fmt.Sprintf(format, a...))
}

// Debug creates a log for helpful information that might not be needed
func Debug(text string) (log Log) {

	log = newLog()
	log.Type = InfoSeverity
	log.Text = text
	log.IsDebug = true

	return log
}

// Debugf creates a log for helpful information that might not be needed allowing a formatting string
// to be passed for the log rather than a straight string parameter
func Debugf(format string, a ...interface{}) (log Log) {
	return Debug(fmt.Sprintf(format, a...))
}

// Warning creates a log for concerning information that does not stop the flow of the application but is concerning
func Warning(text string, err error) (log Log) {

	log = newLog()
	log.Type = WarnSeverity
	log.Text = text
	log.Error = err

	return log
}

// Warningf creates a log for concerning information that does not stop the flow of the application but is concerning
// allowing a formatting string to be passed for the log rather than a straight string parameter
func Warningf(err error, format string, a ...interface{}) (log Log) {
	return Warning(fmt.Sprintf(format, a...), err)
}

// Error creates a log for events that prevent normal flow of the application
func Error(text string, err error) (log Log) {

	log = newLog()
	log.Type = ErrorSeverity
	log.Text = text
	log.Error = err

	return log
}

// Errorf creates a log for events that prevent normal flow of the application
// allowing a formatting string to be passed for the log rather than a straight string parameter
func Errorf(err error, format string, a ...interface{}) (log Log) {
	return Error(fmt.Sprintf(format, a...), err)
}

// Critical creates a log for series issues that occur in the application
func Critical(text string, err error) (log Log) {

	log = newLog()
	log.Type = CritSeverity
	log.Text = text
	log.Error = err

	return log
}

// Criticalf creates a log for series issues that occur in the application allowing a formatting string
// to be passed for the log rather than a straight string parameter
func Criticalf(err error, format string, a ...interface{}) (log Log) {
	return Critical(fmt.Sprintf(format, a...), err)
}

// Fatal creates a log for serious events that occur during execution
func Fatal(text string, err error) (log Log) {

	log = newLog()
	log.Type = FatalSeverity
	log.Text = text
	log.Error = err

	return log
}

// Fatalf creates a log for serious events that occur during execution allowing a formatting string
// to be passed for the log rather than a straight string parameter
//noinspection GoUnusedExportedFunction
func Fatalf(err error, format string, a ...interface{}) (log Log) {
	return Fatal(fmt.Sprintf(format, a...), err)
}
