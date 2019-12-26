package job

import (
	"context"
	"github.com/pkg/errors"
	"github.com/nortonlifelock/log"
	"sync"
)

// Dispatcher is the struct that holds the job queue and helps kick off jobs
type Dispatcher interface {
	Queue(wrapper *jobWrapper) (err error)
	Run() (err error)
	Cancel(ID string)
}

// NewDispatcher returns a dispatcher holds queued jobs and runs them when possible
func NewDispatcher(ctx context.Context, lstream log.Logger, maxWorkers int) (spatcher Dispatcher, err error) {

	if maxWorkers > 0 {

		if lstream != nil {

			// Initialize background context in case of nil context
			if ctx == nil {
				ctx = context.Background()
			}

			// Create the instance of a dispatcher
			spatcher = &dispatcher{
				maxworkers:        maxWorkers,
				workerpool:        make(chan chan *jobWrapper),
				queued:            make(map[string]*jobWrapper),
				queue:             make(chan *jobWrapper),
				ctx:               ctx,
				lstream:           lstream,
				contextMu:         sync.RWMutex{},
				contextRegistries: make(map[string]ctxWrapper),
			}
		} else {
			err = errors.New("log stream was passed nil to the dispatcher")
		}
	} else {
		err = errors.New("max workers must be greater than 0 for dispatcher to initialize")
	}

	return
}

type dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	workerpool chan chan *jobWrapper
	maxworkers int
	queue      chan *jobWrapper
	queued     map[string]*jobWrapper
	ctx        context.Context
	lstream    log.Logger

	contextMu         sync.RWMutex
	contextRegistries map[string]ctxWrapper
	cancelled         sync.Map
}

// Queue adds new jobs to the queue for processing
func (disp *dispatcher) Queue(wrapper *jobWrapper) (err error) {

	// Only queue jobs that exist and aren't already queued
	if disp.queued[wrapper.id] == nil {

		// Create the context for the job so that it can be cancelled
		var cancel context.CancelFunc
		wrapper.ctx, cancel = context.WithCancel(disp.ctx)

		// Add this job as queued so that it doesn't get re-queued
		disp.queued[wrapper.id] = wrapper

		disp.registerContext(
			wrapper.id,
			ctxWrapper{
				wrapper.ctx,
				cancel,
			})

		disp.queue <- wrapper
	} else {
		err = errors.Errorf("Dispatcher: Job [%v] has already been queued", wrapper.id)
	}

	return err
}

// Cancel cancelled the context of a job based on the ID passed
func (disp *dispatcher) Cancel(ID string) {

	// Build a cache of jobs we've already cancelled so we don't attempt to cancel more than once
	if _, loaded := disp.cancelled.LoadOrStore(ID, true); !loaded {

		if disp.contextRegistries[ID].cancel != nil {

			// Pull registered context
			disp.contextMu.RLock()
			disp.lstream.Send(log.Infof("Dispatcher: Cancelling Job with Id [%v]", ID))
			disp.contextRegistries[ID].cancel()
			disp.contextMu.RUnlock()
		}
	}
}

// Run starts all the workers that are responsible for starting jobs
func (disp *dispatcher) Run() (err error) {

	// TODO: should this be moved into the creation of the dispatcher?
	for i := 0; i < disp.maxworkers; i++ {
		var w worker

		if w, err = newWorker(disp.ctx, disp.lstream, disp.workerpool); err == nil {
			w.start()
		} else {
			break
		}
	}

	if err == nil {
		go disp.dispatch()
	}

	return err
}

// Register the context of the job so that it cancelled after having been started
func (disp *dispatcher) registerContext(id string, context ctxWrapper) {
	disp.contextMu.Lock()
	defer disp.contextMu.Unlock()
	if disp.contextRegistries != nil {
		if oldCtx, dup := disp.contextRegistries[id]; dup {
			oldCtx.cancel()
		}

		disp.contextRegistries[id] = context
	} else {
		panic("Job Runner: registerContext context map is nil")
	}
}

func (disp *dispatcher) dispatch() {
	// todo: add handler here and restart the loop

	for {
		select {
		case <-disp.ctx.Done():
			return
		case wrapper := <-disp.queue:

			// a job request has been received
			go func(ctx context.Context, wrapper *jobWrapper) {
				// TODO: add handler here
				select {
				case <-ctx.Done():
					return
				default:
					// try to obtain a worker job channel that is available.
					// disp will block until a worker is idle
					wrapperQueue := <-disp.workerpool

					// dispatch the job to the worker job channel
					wrapperQueue <- wrapper
				}
			}(disp.ctx, wrapper)
		}
	}
}
