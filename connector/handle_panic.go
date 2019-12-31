package connector

import (
	"fmt"
	"github.com/nortonlifelock/log"
	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// HandleRoutinePanic is used in a defer statements for goroutines to prevent a single panic from bringing down the entire application
func handleRoutinePanic(lstream log.Logger) {

	if retVal := recover(); retVal != nil {

		var retValString string
		var ok bool
		retValString, ok = retVal.(string)
		if !ok {
			retValString = fmt.Sprintf("%v", retVal)
		}

		var stackWrapper = errors.New(retValString)
		err, ok := errors.Cause(stackWrapper).(stackTracer)

		panicErr := fmt.Errorf("unhandled panic")
		var newLog = log.Fatal(retValString, panicErr)

		if ok {
			stackTrace := err.StackTrace()

			panicErr = fmt.Errorf(fmt.Sprintf("%s\n%+v\n", retVal, stackTrace[3:]))
			newLog = log.Fatal(retValString, panicErr)

			newLog.Stack(panicErr.Error()) // top two frames
		}

		lstream.Send(newLog)
	}
}
