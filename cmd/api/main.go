package main

import (
	"context"
	"fmt"
	"golang-rest-api/config"
	"golang-rest-api/pkg/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	config.LoadEnvConfig()
	log.InitLogger(log.LoggerMetaData{
		LogLevel:   "",
		Service:    "golang-rest-api",
		AppVersion: "v0.0.0",
	})

	r := chi.NewRouter()

	httpServer := http.Server{
		Addr:              ":" + config.Get().AppPort,
		Handler:           r,
		ReadTimeout:       5 * time.Minute,
		ReadHeaderTimeout: 5 * time.Minute,
		WriteTimeout:      5 * time.Minute,
	}

	// Start HTTP server
	go func() {
		log.Info(context.Background(), "Start Http Server")

		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(context.Background(), fmt.Sprint("Error HTTP server ListenAndServe: ", err))
		}
	}()

	// Capture SIGINT and SIGTERM signals for graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Wait until we receive a shutdown signal
	<-signals

	// Create a context with a 10-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt a graceful shutdown
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error(ctx, "Error server forced to shutdown: ", err)
	}

	log.Info(context.Background(), "Http server exiting gracefully")
}
