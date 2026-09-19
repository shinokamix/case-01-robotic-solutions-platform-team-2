package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/auth"
	authpostgres "github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/auth/postgres"
	"github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/config"
	"github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/httpapi"
	appPostgres "github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/postgres"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		logger.Error("API stopped", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	rootContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := appPostgres.Open(rootContext, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	authService, err := auth.NewService(authpostgres.New(pool), time.Now)
	if err != nil {
		return err
	}
	handler := httpapi.NewAuthHandler(authService, logger, cfg.CookieSecure)
	server := &http.Server{
		Addr:              cfg.APIAddr,
		Handler:           httpapi.NewRouter(handler, logger, cfg.WebOrigin, cfg.TrustedProxyCIDRs...),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("API listening", "address", cfg.APIAddr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case <-rootContext.Done():
		shutdownContext, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownContext); err != nil {
			return err
		}
		return nil
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
