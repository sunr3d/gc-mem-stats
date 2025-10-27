package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sunr3d/gc-mem-stats/gcmemstats"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/metrics", gcmemstats.MetricsHandler())
	gcmemstats.RegisterPprof(mux)

	// oldPercent := gcmemstats.SetGCPercent(150)
	// log.Printf("Настройка частоты срабатывания GC: %d -> 150", oldPercent)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("Старт сервера на :8080")
		log.Println("Метрики в формате Prometheus доступны на: http://localhost:8080/metrics")
		log.Println("Профилирование доступно на: http://localhost:8080/debug/pprof/")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Не удалось поднять сервер: %v", err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	<-done

	log.Println("Остановка сервера...")
	if err := server.Close(); err != nil {
		log.Printf("Не удалось остановить сервер: %v", err)
	}
	log.Println("Сервер остановлен")
}
