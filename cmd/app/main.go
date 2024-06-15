package main

import (
	"UrlShortener/configs"
	"UrlShortener/internal/app"
	"UrlShortener/pkg/logger"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Init configuration
	cfg := configs.MustLoad()

	fmt.Println(cfg)

	// Init logger
	log := logger.SetupLogger(cfg.Log.Level)

	application := app.New(log, cfg)

	// Run servers
	go func() {
		application.HTTPServer.MustRun()
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-stop:
	}

	log.Info("starting graceful shutdown")

	if err := application.HTTPServer.Shutdown(); err != nil {
		log.Error("HTTPServer.Shutdown error", logger.Err(err))
	}

	err := application.DB.Close()

	if err != nil {
		log.Error("redis.Close error", logger.Err(err))
	}

	log.Info("gracefully stopped")
}
