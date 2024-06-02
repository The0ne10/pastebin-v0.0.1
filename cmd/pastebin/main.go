package main

import (
	"fmt"
	"log/slog"
	"os"
	"pastebin-v0.0.1/internal/config"
	"pastebin-v0.0.1/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: init cfg
	cfg := config.MustLoad()

	// TODO: init log/slog
	log := setupLogger(cfg.Env)
	log = log.With("Env", cfg.Env)

	log.Info("initializing logger")
	if envLocal == cfg.Env {
		log.Debug("debug logging enabled")
	}

	// TODO: init storage
	_, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		fmt.Printf("error initializing storage %s", err)
	}

	// TODO: init routes (chi)

	// TODO: run application server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
