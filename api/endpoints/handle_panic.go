package endpoints

import (
	"fmt"
	"github.com/nortonlifelock/log"
	"github.com/pkg/errors"
	"net/http"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// HandleRoutinePanic is used in a defer statements for goroutines to prevent a single panic from bringing down the entire application
func handleRoutinePanic(trans *transaction, w http.ResponseWriter, endpointName string) {

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

		trans.err = fmt.Errorf("%s - %s", panicErr.Error(), retValString)
		trans.wrapper.addError(trans.err, backendError)

		respondToUserWithStatusCode(trans.user, newResponse(trans.obj, trans.totalRecords), w, trans.wrapper, endpointName, http.StatusInternalServerError)
	}
}
