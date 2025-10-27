# gc-mem-stats

Go-библиотека для мониторинга памяти и сборщика мусора. Предоставляет HTTP endpoint с метриками в формате Prometheus и поддержкой профилирования pprof.

## Функциональность

- ✅ **Метрики памяти** - аллокации, heap, stack, системная память
- ✅ **Метрики GC** - количество циклов, время пауз, последний GC
- ✅ **Метрики горутин** - текущее количество активных горутин
- ✅ **Prometheus формат** - стандартный text exposition format
- ✅ **pprof профилирование** - все стандартные endpoints Go
- ✅ **Настройка GC** - изменение агрессивности сборщика мусора
- ✅ **Пример сервера** - готовый к запуску HTTP сервер

## Особенности

- **Простая библиотека** для интеграции в существующие проекты
- **HTTP handler** для метрик в формате Prometheus  
- **pprof endpoints** для профилирования производительности
- **Настройка GC** через `debug.SetGCPercent`
- **Обновление метрик по запросу** - данные собираются при каждом HTTP запросе

## Установка и запуск

### 1. Быстрый старт

```bash
# Установка библиотеки
go get github.com/sunr3d/gc-mem-stats/gcmemstats

# Запуск примера сервера
cd example
go run main.go

# Или из корня проекта
go run example/main.go
```

### 2. Интеграция в проект

```go
package main

import (
    "net/http"

    "github.com/sunr3d/gc-mem-stats/gcmemstats" // импорт пакета
)

func main() {
    mux := http.NewServeMux()

    // Метрики Prometheus
    mux.Handle("/metrics", gcmemstats.MetricsHandler())

    // pprof профилирование
    gcmemstats.RegisterPprof(mux)

    http.ListenAndServe(":8080", mux)
}
```

## API

### Метрики Prometheus

```bash
GET /metrics

# Ответ в формате Prometheus:
# HELP gc_alloc_bytes Объём памяти, выделенной в данный момент (используемой прямо сейчас)
# TYPE gc_alloc_bytes gauge
gc_alloc_bytes 2147768

# HELP gc_num_gc Количество завершённых циклов сборки мусора (GC)
# TYPE gc_num_gc counter
gc_num_gc 7

# HELP gc_num_goroutine Текущее количество горутин (активные + ожидающие)
# TYPE gc_num_goroutine gauge
gc_num_goroutine 5
```

### pprof профилирование

```bash
# Список профилей
GET /debug/pprof/

# Heap profile
GET /debug/pprof/heap

# CPU profile (30 секунд)
GET /debug/pprof/profile?seconds=30

# Goroutines dump
GET /debug/pprof/goroutine?debug=1

# Memory allocations
GET /debug/pprof/allocs
```

## Примеры использования

### Базовые команды

```bash
# Получение метрик
curl http://localhost:8080/metrics

# Создание нагрузки для тестирования GC
for i in {1..1000}; do
  curl http://localhost:8080/metrics > /dev/null
done

# Heap профиль через pprof
go tool pprof http://localhost:8080/debug/pprof/heap

# CPU профиль (30 секунд)
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30

# Dump горутин
curl http://localhost:8080/debug/pprof/goroutine?debug=2
```

### Настройка GC

```go
package main

import (
    "log"

    "github.com/sunr3d/gc-mem-stats/gcmemstats"
)

func main() {
    // Уменьшить потребление памяти (GC чаще)
    oldPercent := gcmemstats.SetGCPercent(50)
    log.Printf("GC percent: %d -> 50", oldPercent)

    // Уменьшить нагрузку на CPU (GC реже)
    gcmemstats.SetGCPercent(200)

    // Отключить GC (не рекомендуется)
    gcmemstats.SetGCPercent(-1)
}
```

## Список метрик

| Метрика                  | Тип     | Описание                                         |
| ------------------------ | ------- | ------------------------------------------------ |
| `gc_alloc_bytes`         | gauge   | Объём памяти, выделенной в данный момент         |
| `gc_total_alloc_bytes`   | counter | Общий объём памяти, выделенный с момента запуска |
| `gc_sys_bytes`           | gauge   | Объём памяти, полученной от ОС                   |
| `gc_heap_alloc_bytes`    | gauge   | Память кучи, выделенная и используемая           |
| `gc_heap_sys_bytes`      | gauge   | Память кучи, полученная от ОС                    |
| `gc_heap_idle_bytes`     | gauge   | Неиспользуемая память кучи                       |
| `gc_heap_inuse_bytes`    | gauge   | Память кучи, занятая приложением                 |
| `gc_heap_released_bytes` | gauge   | Память из кучи, возвращённая ОС                  |
| `gc_heap_objects`        | gauge   | Количество объектов в куче                       |
| `gc_stack_inuse_bytes`   | gauge   | Память стеков, используемых горутинами           |
| `gc_stack_sys_bytes`     | gauge   | Память стеков, выделенная от ОС                  |
| `gc_mallocs_total`       | counter | Количество вызовов malloc                        |
| `gc_frees_total`         | counter | Количество освобождений памяти                   |
| `gc_num_gc`              | counter | Количество завершённых циклов GC                 |
| `gc_pause_total_ns`      | counter | Общее время пауз GC в наносекундах               |
| `gc_last_gc_timestamp`   | gauge   | Время последнего GC (timestamp)                  |
| `gc_num_goroutine`       | gauge   | Текущее количество горутин                       |

## Структура проекта

```
├── gcmemstats/              # Основная библиотека
│   ├── collector.go         # Сбор метрик из runtime.ReadMemStats
│   ├── prom_formatter.go    # Форматирование в Prometheus
│   ├── handlers.go          # HTTP handlers и pprof (API)
├── example/                 # Пример использования
```

### Команды разработки

```bash
# Запуск тестов
go test -v ./...

# Запуск примера
go run example/main.go

# Нагрузочное тестирование
./example/smoke_test.sh

# Проверка метрик
curl http://localhost:8080/metrics | head -20
```

### Тестирование

```bash
# Все тесты
go test -v ./gcmemstats/...

# Тесты с race detector
go test -race -v ./gcmemstats/...

# Нагрузочное тестирование (требуется запущенный сервер)
cd example && ./smoke_test.sh
```
