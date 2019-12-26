package job

import (
	"context"
	"github.com/pkg/errors"
	"github.com/nortonlifelock/log"
)

// worker represents the worker that executes the job
type worker struct {
	pool    chan chan *jobWrapper
	queue   chan *jobWrapper
	lstream log.Logger
	ctx     context.Context
}

func newWorker(ctx context.Context, lstream log.Logger, pool chan chan *jobWrapper) (w worker, err error) {

	// TODO: The worker should have it's own context initialized
	w = worker{
		pool:    pool,
		queue:   make(chan *jobWrapper),
		ctx:     ctx,
		lstream: lstream,
	}

	return w, err
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (worker *worker) start() {
	// todo need to restart in the event of a failure

	go func() {
		// TODO: Need to add a handler here
		for {
			// TODO: look for a better way to do this
			// register the current worker into the worker queue.
			worker.pool <- worker.queue

			select {
			case <-worker.ctx.Done():
				// we have received a signal to stop
				return

			case wrapper := <-worker.queue:
				worker.process(wrapper)
			}
		}
	}()
}
func (worker *worker) process(wrapper *jobWrapper) {
	defer func() {
		if recoverMessage := recover(); recoverMessage != nil {
			var err = errors.Errorf("Job Runner: Job [%s:%v] panicked during execution with the error - [%s]", wrapper.name, wrapper.id, recoverMessage)
			worker.lstream.Send(log.Fatal(err.Error(), err))

			// TODO: Drop the job back on the channel to be processed, or update to error
		}
	}()

	worker.lstream.Send(log.Infof("Job Runner: Pulled job [%s:%v] from job channel - Beginning processing...", wrapper.name, wrapper.id))
	if err := wrapper.Execute(); err == nil {
		worker.lstream.Send(log.Infof("Job Runner: Job [%s:%v] successfully completed", wrapper.name, wrapper.id))
	} else {
		worker.lstream.Send(log.Errorf(err, "Job Runner: An error occurred while processing job [%s:%v]", wrapper.name, wrapper.id))
	}
}
