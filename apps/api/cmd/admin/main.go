package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"syscall"
	"time"

	"golang.org/x/term"

	"github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/auth"
	authpostgres "github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/auth/postgres"
	"github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/config"
	appPostgres "github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/postgres"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if len(os.Args) != 2 || os.Args[1] != "create" {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/admin create")
		os.Exit(2)
	}
	if err := createAdmin(logger); err != nil {
		logger.Error("create administrator", "error", err)
		os.Exit(1)
	}
}

func createAdmin(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	input, err := adminInput()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	pool, err := appPostgres.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	service, err := auth.NewService(authpostgres.New(pool), time.Now)
	if err != nil {
		return err
	}
	user, err := service.CreateAdmin(ctx, input)
	if err != nil {
		return err
	}
	logger.Info("administrator created", "id", user.ID, "email", user.Email)
	return nil
}

func adminInput() (auth.AdminInput, error) {
	reader := bufio.NewReader(os.Stdin)
	email, err := readValue(reader, "ADMIN_EMAIL", "Email: ", "")
	if err != nil {
		return auth.AdminInput{}, err
	}
	firstName, err := readValue(reader, "ADMIN_FIRST_NAME", "Имя: ", "Администратор")
	if err != nil {
		return auth.AdminInput{}, err
	}
	lastName, err := readValue(reader, "ADMIN_LAST_NAME", "Фамилия: ", "Системы")
	if err != nil {
		return auth.AdminInput{}, err
	}
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		fmt.Fprint(os.Stderr, "Пароль: ")
		bytes, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return auth.AdminInput{}, err
		}
		password = string(bytes)
	}
	return auth.AdminInput{Email: email, FirstName: firstName, LastName: lastName, Password: password}, nil
}

func readValue(reader *bufio.Reader, environmentName, prompt, fallback string) (string, error) {
	if value := os.Getenv(environmentName); value != "" {
		return value, nil
	}
	fmt.Fprint(os.Stderr, prompt)
	value, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback, nil
	}
	return value, nil
}
