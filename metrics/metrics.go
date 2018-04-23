package metrics

import (
	"expvar"
	"time"
)

var (
	// Uptime stores the timestamp of the SqsProcessor boot
	Uptime = expvar.NewInt("uptime")
)

// Init starts metrics collection
func Init() {
	Uptime.Set(time.Now().Unix())
}