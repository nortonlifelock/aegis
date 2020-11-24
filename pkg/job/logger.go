package job

import (
	"fmt"
	"github.com/nortonlifelock/aegis/pkg/log"
)

type logger struct {
	jobID   int
	jobName string
	logs    log.Logger
}

// Send updates the log that is passed in with job information added to the beginning for job
// specific logging and then drops it on the system logger
func (l logger) Send(log log.Log) {

	log.Text = fmt.Sprintf("[%s: %v] %s", l.jobName, l.jobID, log.Text)

	l.logs.Send(log)
}
