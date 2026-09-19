package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/db/migrations"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		logger.Error("DATABASE_URL is required")
		os.Exit(1)
	}
	if err := migrations.Up(context.Background(), databaseURL); err != nil {
		logger.Error("migration failed", "error", err)
		os.Exit(1)
	}
	logger.Info("migrations applied")
}
