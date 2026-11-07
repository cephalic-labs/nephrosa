package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cephalic-labs/nephrosa/backend/internal/config"
	"github.com/cephalic-labs/nephrosa/backend/internal/database"
	"github.com/cephalic-labs/nephrosa/backend/internal/httpapi"
	"github.com/cephalic-labs/nephrosa/backend/internal/migrate"
	"github.com/cephalic-labs/nephrosa/backend/internal/sarvam"
	"github.com/cephalic-labs/nephrosa/backend/internal/seed"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		slog.Error("renal relay command failed", "error", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	command := "serve"
	if len(args) > 0 {
		command = args[0]
	}

	switch command {
	case "healthcheck":
		return healthcheck()
	case "serve", "migrate", "seed", "reset-demo":
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		pool, err := database.Open(ctx, cfg.Database, command != "serve")
		if err != nil {
			return err
		}
		defer pool.Close()

		switch command {
		case "migrate":
			return migrate.Apply(ctx, pool)
		case "reset-demo":
			if cfg.AppEnv != "development" || cfg.Database.Name != "renalrelay" || len(args) != 2 || args[1] != "--confirm-fictional-reset" {
				return errors.New("reset requires development, database renalrelay, and --confirm-fictional-reset")
			}
			return seed.Reset(ctx, pool)
		case "seed":
			if cfg.AppEnv != "development" {
				return errors.New("fictional demo data can only be seeded when APP_ENV=development")
			}
			return seed.Apply(ctx, pool)
		default:
			return serve(cfg, pool)
		}
	default:
		return fmt.Errorf("unknown command %q; use serve, migrate, seed, reset-demo, or healthcheck", command)
	}
}

func serve(cfg config.Config, pool *database.Pool) error {
	voiceCfg, err := sarvam.Load()
	if err != nil {
		return err
	}
	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           httpapi.New(pool, sarvam.New(voiceCfg)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serveErr := make(chan error, 1)
	go func() {
		slog.Info("Renal Relay API listening", "address", cfg.HTTPAddr, "environment", cfg.AppEnv)
		serveErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serveErr:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shut down HTTP server: %w", err)
		}
		return nil
	}
}

func healthcheck() error {
	client := &http.Client{Timeout: 2 * time.Second}
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://127.0.0.1:8080/health/live", nil)
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("health endpoint returned %s", response.Status)
	}
	return nil
}
