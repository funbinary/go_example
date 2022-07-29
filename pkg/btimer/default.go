package btimer

import (
	"context"
	"log"
	"time"
)

var defaultTimer = New()

type JobFunc = func(ctx context.Context)

func SetTimeout(ctx context.Context, delay time.Duration, job JobFunc) {
	AddOnce(ctx, delay, job)
}

func SetInterval(ctx context.Context, interval time.Duration, job JobFunc) {
	Add(ctx, interval, job)
}

func Add(ctx context.Context, interval time.Duration, job JobFunc) *Entry {
	return defaultTimer.Add(ctx, interval, job)
}

func AddEntry(ctx context.Context, interval time.Duration, job JobFunc, isSingleton bool, times int32, status int32) *Entry {
	return defaultTimer.AddEntry(ctx, interval, job, isSingleton, times, status)
}

func AddSingleton(ctx context.Context, interval time.Duration, job JobFunc) *Entry {
	return defaultTimer.AddSingleton(ctx, interval, job)
}

func AddOnce(ctx context.Context, interval time.Duration, job JobFunc) *Entry {
	return defaultTimer.AddOnce(ctx, interval, job)
}

func AddTimes(ctx context.Context, interval time.Duration, times int32, job JobFunc) *Entry {
	return defaultTimer.AddTimes(ctx, interval, times, job)
}

func DelayAdd(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc) {
	defaultTimer.DelayAdd(ctx, delay, interval, job)
}

func DelayAddEntry(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc, isSingleton bool, times int32, status int32) {
	defaultTimer.DelayAddEntry(ctx, delay, interval, job, isSingleton, times, status)
}

func DelayAddSingleton(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc) {
	defaultTimer.DelayAddSingleton(ctx, delay, interval, job)
}

func DelayAddOnce(ctx context.Context, delay time.Duration, interval time.Duration, job JobFunc) {
	defaultTimer.DelayAddOnce(ctx, delay, interval, job)
}

func DelayAddTimes(ctx context.Context, delay time.Duration, interval time.Duration, times int32, job JobFunc) {
	defaultTimer.DelayAddTimes(ctx, delay, interval, times, job)
}

func Exit() {
	log.Println(panicExit)
}
