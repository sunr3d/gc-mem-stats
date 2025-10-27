package gcmemstats

import (
	"runtime"
)

type memStats struct {
	stats        runtime.MemStats
	numGoroutine int
}

func collectStats() *memStats {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	return &memStats{
		stats:        ms,
		numGoroutine: runtime.NumGoroutine(),
	}
}
