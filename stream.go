package log

import (
	"context"
	"fmt"
	"github.com/nortonlifelock/connection"
	"github.com/nortonlifelock/files"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	directoryPathFileFormat = "2006-01-02"
)

// Stream is a type alias for the log logStream to make it less awkward
type Stream = chan<- Log

type config interface {
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

type logStream struct {
	path            string
	logs            chan Log
	logToFile       bool
	logPreserveFile bool
	fileLogs        chan Log
	logToConsole    bool
	consoleLogs     chan Log
	logToDb         bool
	dbLogs          chan Log
	logToMessageQ   bool
	mqLogs          chan Log
	debug           bool
	snsLogs         chan Log
	logToSNS        bool
	dbconn          connection.DatabaseConnection
	snsClient       *SNSClient
}

// NewLogStream creates and returns a logger
func NewLogStream(ctx context.Context, dbconn connection.DatabaseConnection, logconfig config) (Logger, error) {
	var err error
	// Are there any consequences of creating a buffered channel?
	var lstream = make(chan Log, 100)

	var snsClient *SNSClient
	if len(logconfig.SNSTopicID()) > 0 {
		snsClient, _ = NewSNSClient(ctx, logconfig.SNSTopicID())
	}

	var logger = &logStream{
		path:            logconfig.LogPath(),
		logs:            lstream,
		logToFile:       logconfig.LogFile(),
		logPreserveFile: logconfig.PreserveFileLogs(),
		fileLogs:        make(chan Log),
		logToConsole:    logconfig.LogConsole(),
		consoleLogs:     make(chan Log),
		logToDb:         logconfig.LogDB(),
		dbLogs:          make(chan Log),
		logToMessageQ:   logconfig.LogMQ(),
		mqLogs:          make(chan Log),
		debug:           logconfig.DebugLogs(),
		snsLogs:         make(chan Log),
		logToSNS:        len(logconfig.SNSTopicID()) > 0,
		dbconn:          dbconn,
		snsClient:       snsClient,
	}

	if logger.logToFile {
		go logger.publish(ctx, logger.fileLogs, logger.writeToFile)
	}

	if logger.logToConsole {
		go logger.publish(ctx, logger.consoleLogs, logger.write)
	}

	if logger.logToDb {
		go logger.publish(ctx, logger.dbLogs, logger.writeToDB)
	}

	if logger.logToMessageQ {
		go logger.publish(ctx, logger.mqLogs, logger.writeToMessageQ)
	}

	if logger.logToSNS {
		go logger.publish(ctx, logger.snsLogs, logger.sns)
	}

	go logger.publish(ctx, logger.logs, logger.distribute)

	return logger, err
}

// Send pushes a log onto the log channel
func (stream *logStream) Send(log Log) {
	// Don't want to send as a goroutine, as scheduling could create race conditions in the order logs are pushed
	stream.logs <- log
}

func (stream *logStream) publish(ctx context.Context, inStream <-chan Log, logMethod func(log Log)) {
	defer func() {
		select {
		case <-ctx.Done():
			return
		default:
			handle(stream.logs)
			stream.logs <- Warning("Restarting the Log Publish Stream due to a panic", nil)
			go stream.publish(ctx, inStream, logMethod)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case log := <-inStream:
			// Execute as a go routine so that the sender doesn't have to wait for the
			// log to be written to the source
			logMethod(log)
		}
	}
}

func (stream *logStream) distribute(log Log) {

	// Filter out the debug logs when they're disabled
	if !log.IsDebug || (log.IsDebug && stream.debug) {

		if stream.logToConsole {
			go func() {
				defer handle(stream.logs)
				stream.consoleLogs <- log
			}()
		}

		if stream.logToFile {
			go func() {
				defer handle(stream.logs)
				stream.fileLogs <- log
			}()
		}

		if stream.logToDb {
			go func() {
				defer handle(stream.logs)
				stream.dbLogs <- log
			}()
		}

		if stream.logToMessageQ {
			go func() {
				defer handle(stream.logs)
				stream.mqLogs <- log
			}()
		}
	}
}

func (stream *logStream) sns(log Log) {
	if log.Type == FatalSeverity || log.Type == CritSeverity {
		if stream.snsClient != nil {
			stream.snsClient.PushMessage(log.ToConsoleString())
		}
	}
}

func (stream *logStream) write(log Log) {
	fmt.Println(log.ToConsoleString())
}

func (stream *logStream) writeToFile(log Log) {
	var path = fmt.Sprintf("%s%c%s.log", stream.path, os.PathSeparator, time.Now().Format(directoryPathFileFormat))

	var logContents = fmt.Sprintf("%s\n", log.ToConsoleString())

	if _, err := os.Stat(path); err == nil {

		var logFile *os.File
		logFile, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err == nil {
			defer func() {
				_ = logFile.Close()
			}()
			_, err = logFile.WriteString(logContents)
		} else {
			fmt.Println(fmt.Sprintf("FAILED APPENDING LOG CONTENTS - %s", err.Error()))
		}

	} else { //The file didn't exist yet

		if err = stream.newLog(); err == nil {

			err = files.WriteFile(path, logContents)
			if err == nil {

				if !stream.logPreserveFile {
					// deletes log files older than a day
					_ = stream.deleteLogsOlderThanDay()
				}
			} else {
				fmt.Println("FAILED WRITING INITIAL LOG - ", err.Error())
			}
		} else {
			fmt.Printf("error when creating directory paths for logs | %s\n", err.Error())
		}
	}
}

func (stream *logStream) writeToDB(log Log) {
	var errText string
	if log.Error != nil {
		errText = log.Error.Error()
	}

	var jobID string
	jobID = log.JobID

	_, _, _ = stream.dbconn.CreateLog(log.Type, log.Text, errText, jobID, time.Now())
}

func (stream *logStream) writeToMessageQ(log Log) {

}

func (stream *logStream) newLog() (err error) {

	if len(stream.path) > 0 {
		var dirPath = filepath.Dir(stream.path)

		if _, err = os.Stat(dirPath); err != nil {
			_ = os.MkdirAll(dirPath, 0775)
		}

	} else {
		fmt.Println("WARNING: LOGPATH CONFIG VARIABLE EMPTY - CANNOT LOG")
	}

	return err
}

func (stream *logStream) deleteLogsOlderThanDay() (err error) {
	var dirpath = fmt.Sprintf(
		"%s%s",
		stream.path,
		string(os.PathSeparator),
	)

	var filesInDir []os.FileInfo
	if filesInDir, err = ioutil.ReadDir(stream.path); err == nil {
		for index := range filesInDir {

			file := filesInDir[index]
			var fileTime time.Time
			var filename = file.Name()
			filename = strings.Replace(filename, ".log", "", 1)

			fileTime, err = time.Parse(directoryPathFileFormat, filename)
			if err == nil {
				if time.Since(fileTime) >= time.Hour*24 {
					err = os.Remove(fmt.Sprintf("%s%s", dirpath, file.Name()))
					if err != nil {
						fmt.Println(err.Error())
					}
				}
			}
		}
	}

	return err
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func handle(lstream Stream) {
	if retVal := recover(); retVal != nil {

		var stackWrapper = errors.New(retVal.(string))
		err, ok := errors.Cause(stackWrapper).(stackTracer)

		var newLog = Fatal(retVal.(string), fmt.Errorf("routine panic"))

		if ok {
			stackTrace := err.StackTrace()
			newLog.Stack(fmt.Sprintf("%s\n%+v\n", retVal, stackTrace[3:])) // top two frames
		}

		lstream <- newLog
	}
}
