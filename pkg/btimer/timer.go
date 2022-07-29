package btimer

import (
	"context"
	"go.uber.org/atomic"
	"sync"
	"time"
)

const panicExit = "exit"

type Timer struct {
	sync.RWMutex
	queue   *priorityQueue
	status  *atomic.Int32
	ticks   *atomic.Int64
	options Option
}

func New(opt ...Option) *Timer {
	t := &Timer{
		queue:  newPriorityQueue(),
		status: atomic.NewInt32(StatusRunning),
		ticks:  atomic.NewInt64(0),
	}
	if len(opt) > 0 {
		t.options = opt[0]
	} else {
		t.options = DefaultOptions()
	}
	go t.loop()
	return t
}

type Option struct {
	Interval time.Duration
}

func DefaultOptions() Option {
	return Option{
		Interval: defaultInterval,
	}
}

func (t *Timer) Start() {
	t.status.Store(StatusRunning)
}

// Stop stops the timer.
func (t *Timer) Stop() {
	t.status.Store(StatusStopped)
}

// Close closes the timer.
func (t *Timer) Close() {
	t.status.Store(StatusClosed)
}

func (t *Timer) Add(ctx context.Context, interval time.Duration, job JobFunc) *Entry {
	return t.createEntry(createEntryInput{
		Ctx:         ctx,
		Interval:    interval,
		Job:         job,
		IsSingleton: false,
		Times:       -1,
		Status:      StatusReady,
	})
}

func (t *Timer) AddEntry(ctx context.Context, interval time.Duration, job JobFunc, isSingleton bool, times int32, status int32) *Entry {
	return t.createEntry(createEntryInput{
		Ctx:         ctx,
		Interval:    interval,
		Job:         job,
		IsSingleton: isSingleton,
		Times:       times,
		Status:      status,
	})
}

func (t *Timer) AddSingleton(ctx context.Context, interval time.Duration, job JobFunc) *Entry {
	return t.createEntry(createEntryInput{
		Ctx:         ctx,
		Interval:    interval,
		Job:         job,
		IsSingleton: true,
		Times:       -1,
		Status:      StatusReady,
	})
}

func (t *Timer) AddOnce(ctx context.Context, interval time.Duration, job JobFunc) *Entry {
	return t.createEntry(createEntryInput{
		Ctx:         ctx,
		Interval:    interval,
		Job:         job,
		IsSingleton: true,
		Times:       1,
		Status:      StatusReady,
	})
}

func (t *Timer) AddTimes(ctx context.Context, interval time.Duration, times int32, job JobFunc) *Entry {
	return t.createEntry(createEntryInput{
		Ctx:         ctx,
		Interval:    interval,
		Job:         job,
		IsSingleton: true,
		Times:       times,
		Status:      StatusReady,
	})
}

func (t *Timer) DelayAdd(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(ctx, delay, func(ctx context.Context) {
		t.Add(ctx, interval, job)
	})
}

func (t *Timer) DelayAddEntry(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc, isSingleton bool, times int32, status int32) {
	t.AddOnce(ctx, delay, func(ctx context.Context) {
		t.AddEntry(ctx, interval, job, isSingleton, times, status)
	})
}

func (t *Timer) DelayAddSingleton(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(ctx, delay, func(ctx context.Context) {
		t.AddSingleton(ctx, interval, job)
	})
}

func (t *Timer) DelayAddOnce(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc) {
	t.AddOnce(ctx, delay, func(ctx context.Context) {
		t.AddOnce(ctx, interval, job)
	})
}

func (t *Timer) DelayAddTimes(ctx context.Context, delay time.Duration, interval time.Duration, times int32, job JobFunc) {
	t.AddOnce(ctx, delay, func(ctx context.Context) {
		t.AddTimes(ctx, interval, times, job)
	})
}

func (self *Timer) loop() {
	go func() {
		var (
			currentTimerTicks   int64
			timerIntervalTicker = time.NewTicker(self.options.Interval)
		)
		defer timerIntervalTicker.Stop()
		for {
			select {
			case <-timerIntervalTicker.C:
				// 校验定时器状态
				switch self.status.Load() {
				case StatusRunning:
					// Timer proceeding.
					if currentTimerTicks = self.ticks.Add(1); currentTimerTicks >= self.queue.NextPriority() {
						self.run(currentTimerTicks)
					}
				case StatusStopped:

				case StatusClosed:
					// 定时器退出
					return
				}
			}
		}
	}()
}

func (self *Timer) run(currentTimerTicks int64) {
	var (
		value interface{}
	)
	for {
		value = self.queue.Pop()
		if value == nil {
			break
		}
		entry := value.(*Entry)
		// It checks if it meets the ticks' requirement.
		if jobNextTicks := entry.nextTicks.Load(); currentTimerTicks < jobNextTicks {
			// It pushes the job back if current ticks does not meet its running ticks requirement.
			self.queue.Push(entry, entry.nextTicks.Load())
			break
		}
		// It checks the job running requirements and then does asynchronous running.
		entry.doCheckAndRunByTicks(currentTimerTicks)
		// Status check: push back or ignore it.
		if entry.Status() != StatusClosed {
			// It pushes the job back to queue for next running.
			self.queue.Push(entry, entry.nextTicks.Load())
		}
	}
}

type createEntryInput struct {
	Ctx         context.Context
	Interval    time.Duration
	Job         JobFunc
	IsSingleton bool
	Times       int32
	Status      int32
}

// createEntry creates and adds a timing job to the timer.
func (self *Timer) createEntry(in createEntryInput) *Entry {
	var (
		infinite = false
	)
	if in.Times <= 0 {
		infinite = true
	}
	var (
		intervalTicksOfJob = int64(in.Interval / self.options.Interval)
	)
	if intervalTicksOfJob == 0 {
		// If the given interval is lesser than the one of the wheel,
		// then sets it to one tick, which means it will be run in one interval.
		intervalTicksOfJob = 1
	}
	var (
		nextTicks = self.ticks.Load() + intervalTicksOfJob
		entry     = &Entry{
			job:         in.Job,
			ctx:         in.Ctx,
			timer:       self,
			ticks:       intervalTicksOfJob,
			times:       atomic.NewInt32(in.Times),
			status:      atomic.NewInt32(in.Status),
			isSingleton: atomic.NewBool(in.IsSingleton),
			nextTicks:   atomic.NewInt64(nextTicks),
			infinite:    atomic.NewBool(infinite),
		}
	)
	self.queue.Push(entry, nextTicks)
	return entry
}
