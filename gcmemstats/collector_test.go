package gcmemstats

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollect(t *testing.T) {
	sts := collectStats()

	require.NotNil(t, sts)

	assert.Greater(t, sts.stats.Alloc, uint64(0), "Alloc должен быть больше 0")
	assert.Greater(t, sts.stats.TotalAlloc, uint64(0), "TotalAlloc должен быть больше 0")
	assert.Greater(t, sts.stats.Sys, uint64(0), "Sys должен быть больше 0")

	assert.GreaterOrEqual(t, sts.stats.HeapAlloc, uint64(0), "HeapAlloc должен быть >= 0")
	assert.Greater(t, sts.stats.HeapSys, uint64(0), "HeapSys должен быть больше 0")

	assert.Greater(t, sts.numGoroutine, 0, "NumGoroutine должен быть больше 0")

	assert.GreaterOrEqual(t, sts.stats.NumGC, uint32(0), "NumGC должен быть >= 0")
	assert.GreaterOrEqual(t, sts.stats.PauseTotalNs, uint64(0), "PauseTotalNs должен быть >= 0")
}

func TestCollectConsistency(t *testing.T) {
	sts1 := collectStats()
	sts2 := collectStats()

	assert.GreaterOrEqual(t, sts2.stats.TotalAlloc, sts1.stats.TotalAlloc, "TotalAlloc должен только расти")

	assert.GreaterOrEqual(t, sts2.stats.Mallocs, sts1.stats.Mallocs, "Mallocs должен только расти")
}
