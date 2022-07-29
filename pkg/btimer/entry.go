package btimer

import (
	"context"
	"log"

	"github.com/funbinary/go_example/pkg/errors"
	"go.uber.org/atomic"
)

type Entry struct {
	job         JobFunc         // The job function.
	ctx         context.Context // The context for the job, for READ ONLY.
	timer       *Timer          // Belonged timer.
	ticks       int64           // The job runs every tick.
	times       *atomic.Int32   // Limit running times.
	status      *atomic.Int32   // Job status.
	isSingleton *atomic.Bool    // Singleton mode.
	nextTicks   *atomic.Int64   // Next run ticks of the job.
	infinite    *atomic.Bool    // No times limit.
}

func (entry *Entry) Run() {
	if !entry.infinite.Load() {
		leftRunningTimes := entry.times.Add(-1)
		// It checks its running times exceeding.
		if leftRunningTimes < 0 {
			entry.SetStatus(StatusClosed)
			return
		}
	}
	go func() {
		defer func() {
			if exception := recover(); exception != nil {
				if exception != panicExit {
					if v, ok := exception.(error); ok {
						log.Panicln(v)
					} else {
						panic(errors.Errorf(`exception recovered: %+v`, exception))
					}
				} else {
					entry.Close()
					return
				}
			}
			if entry.Status() == StatusRunning {
				entry.SetStatus(StatusReady)
			}
		}()
		entry.job(entry.ctx)
	}()
}

func (entry *Entry) Start() {
	entry.status.Store(StatusReady)
}

// Stop stops the job.
func (entry *Entry) Stop() {
	entry.status.Store(StatusStopped)
}

// Close closes the job, and then it will be removed from the timer.
func (entry *Entry) Close() {
	entry.status.Store(StatusClosed)
}

func (entry *Entry) Reset() {
	entry.nextTicks.Store(entry.timer.ticks.Load() + entry.ticks)
}

func (entry *Entry) Status() int32 {
	return entry.status.Load()
}

func (entry *Entry) SetStatus(status int32) int32 {
	entry.status.Store(status)
	return entry.status.Load()
}

func (entry *Entry) IsSingleton() bool {
	return entry.isSingleton.Load()
}

func (entry *Entry) Job() JobFunc {
	return entry.job
}

// Ctx returns the initialized context of this job.
func (entry *Entry) Ctx() context.Context {
	return entry.ctx
}

// SetTimes sets the limit running times for the job.
func (entry *Entry) SetTimes(times int32) {
	entry.times.Store(times)
	entry.infinite.Store(false)
}

func (entry *Entry) doCheckAndRunByTicks(currentTimerTicks int64) {
	// Ticks check.
	if currentTimerTicks < entry.nextTicks.Load() {
		return
	}
	entry.nextTicks.Store(currentTimerTicks + entry.ticks)
	// Perform job checking.
	switch entry.status.Load() {
	case StatusRunning:
		if entry.IsSingleton() {
			return
		}
	case StatusReady:
		if !entry.status.CAS(StatusReady, StatusRunning) {
			return
		}
	case StatusStopped:
		return
	case StatusClosed:
		return
	}
	// Perform job running.
	entry.Run()
}
