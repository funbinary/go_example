package btimer

import "time"

const (
	StatusReady   = 0
	StatusRunning = 1
	StatusStopped = 2
	StatusClosed  = -1
)

var defaultInterval = 100 * time.Millisecond
