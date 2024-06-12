package app

import (
	"app/internal/config"
	"app/internal/delivery/http/handlers/text/getText"
	"app/internal/delivery/http/handlers/text/saveText"
	"app/internal/lib/schedule"
	"app/internal/services/minio/client"
	"app/internal/storage/sqlite"
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Run() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log = log.With("Env", cfg.Env)

	log.Info("initializing logger")
	if envLocal == cfg.Env {
		log.Debug("debug logging enabled")
	}

	// storage DB
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		fmt.Printf("error initializing storage %s", err)
	}
	// S3Minio
	minio, err := client.New(log, cfg)
	if err != nil {
		fmt.Printf("error initializing minio client %s", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{alias}", getText.ReadTextHandler(log, storage, minio))
	r.Post("/", saveText.SaveTextHandler(log, storage, minio))

	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      r,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Run schedule
	schedule.StartSchedule(log, storage, minio, time.Minute)

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Error("graceful shutdown timed out.. forcing exit.", slog.String("error", err.Error()))
			}
		}()

		// Trigger graceful shutdown
		err = server.Shutdown(shutdownCtx)
		if err != nil {
			log.Error("error shutting down server gracefully", slog.String("error", err.Error()))
		}
		serverStopCtx()
	}()

	// Run the server
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("error starting server", slog.String("error", err.Error()))
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
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
