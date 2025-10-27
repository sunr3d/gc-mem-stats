package gcmemstats

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatProm(t *testing.T) {
	stats := collectStats()
	result := formatProm(stats)

	require.NotEmpty(t, result)

	assert.Contains(t, string(result), "gc_alloc_bytes")
	assert.Contains(t, string(result), "gc_total_alloc_bytes")
	assert.Contains(t, string(result), "gc_num_gc")
	assert.Contains(t, string(result), "gc_num_goroutine")

	assert.Contains(t, string(result), "# HELP gc_alloc_bytes")
	assert.Contains(t, string(result), "# TYPE gc_alloc_bytes gauge")

	assert.Contains(t, string(result), "gc_alloc_bytes ")
	assert.Contains(t, string(result), "gc_num_goroutine ")
}

func TestFormatPromStructure(t *testing.T) {
	stats := collectStats()
	result := formatProm(stats)

	lines := strings.Split(string(result), "\n")

	helpCount := 0
	typeCount := 0
	valueCount := 0

	for i, line := range lines {
		if strings.HasPrefix(line, "# HELP") {
			helpCount++
			if i+1 < len(lines) {
				assert.True(t, strings.HasPrefix(lines[i+1], "# TYPE"), "После HELP должна быть TYPE")
			}
		} else if strings.HasPrefix(line, "# TYPE") {
			typeCount++
		} else if !strings.HasPrefix(line, "#") && strings.Contains(line, " ") && line != "" {
			valueCount++
		}
	}

	assert.Equal(t, helpCount, typeCount, "Количество HELP и TYPE должно совпадать")
	assert.Equal(t, helpCount, valueCount, "Количество HELP и VALUES должно совпадать")
}
