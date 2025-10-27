package gcmemstats

import (
	"bytes"
	"fmt"
)

// formatProm - собирает статистику в формате Prometheus
func formatProm(sts *memStats) []byte {
	var buf bytes.Buffer

	metric := func(name, typ, help string, value any) {
		buf.WriteString(fmt.Sprintf("# HELP %s %s\n", name, help))
		buf.WriteString(fmt.Sprintf("# TYPE %s %s\n", name, typ))
		buf.WriteString(fmt.Sprintf("%s %v\n\n", name, value))
	}

	// Аллокации памяти общее
	metric("gc_alloc_bytes", "gauge", "Объём памяти, выделенной в данный момент (используемой прямо сейчас)", sts.stats.Alloc)
	metric("gc_total_alloc_bytes", "counter", "Общий объём памяти, выделенный с момента запуска программы", sts.stats.TotalAlloc)
	metric("gc_sys_bytes", "gauge", "Объём памяти, полученной от ОС (включая кучу, стеки и прочее)", sts.stats.Sys)

	// Статистика Heap
	metric("gc_heap_alloc_bytes", "gauge", "Память кучи, выделенная и всё ещё используемая", sts.stats.HeapAlloc)
	metric("gc_heap_sys_bytes", "gauge", "Память кучи, полученная от ОС", sts.stats.HeapSys)
	metric("gc_heap_idle_bytes", "gauge", "Неиспользуемая память кучи, которую можно вернуть ОС", sts.stats.HeapIdle)
	metric("gc_heap_inuse_bytes", "gauge", "Память кучи, реально занятая приложением", sts.stats.HeapInuse)
	metric("gc_heap_released_bytes", "gauge", "Память из кучи, возвращённая ОС", sts.stats.HeapReleased)
	metric("gc_heap_objects", "gauge", "Количество объектов в куче (mallocs − frees)", sts.stats.HeapObjects)

	// Статистика Stack
	metric("gc_stack_inuse_bytes", "gauge", "Память стеков, используемых активными горутинами", sts.stats.StackInuse)
	metric("gc_stack_sys_bytes", "gauge", "Память стеков, выделенные от ОС (включая неиспользуемую)", sts.stats.StackSys)

	// Malloc/Free
	metric("gc_mallocs_total", "counter", "Количество вызовов malloc с начала работы программы", sts.stats.Mallocs)
	metric("gc_frees_total", "counter", "Количество освобождений памяти (free) с начала работы программы", sts.stats.Frees)

	// GC
	metric("gc_num_gc", "counter", "Количество завершённых циклов сборки мусора (GC)", sts.stats.NumGC)
	metric("gc_pause_total_ns", "counter", "Общее время (в наносекундах), проведённое в паузах GC (STW)", sts.stats.PauseTotalNs)

	// Время вызова последнего цикла GC
	if sts.stats.LastGC > 0 {
		metric("gc_last_gc_timestamp", "gauge", "Время вызова последнего цикла GC (timestamp)", sts.stats.LastGC)
	}

	// Горутины
	metric("gc_num_goroutine", "gauge", "Текущее количество горутин (активные + ожидающие)", sts.numGoroutine)

	return buf.Bytes()
}
