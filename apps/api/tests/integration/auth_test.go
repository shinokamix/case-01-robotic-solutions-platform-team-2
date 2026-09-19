package integration

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/db/migrations"
	"github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/auth"
	authpostgres "github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/auth/postgres"
	"github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/httpapi"
)

const testOrigin = "http://localhost:5173"

func TestAuthenticationHTTP(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	if err := migrations.Up(ctx, databaseURL); err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	if _, err := pool.Exec(ctx, "TRUNCATE users CASCADE"); err != nil {
		t.Fatal(err)
	}

	now := time.Date(2026, 9, 19, 12, 0, 0, 0, time.UTC)
	service, err := auth.NewService(authpostgres.New(pool), func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server := httptest.NewServer(httpapi.NewRouter(httpapi.NewAuthHandler(service, logger, false), logger, testOrigin))
	defer server.Close()

	response := requestJSON(t, server.Client(), http.MethodPost, server.URL+"/api/auth/register", `{
		"email":" USER@EXAMPLE.COM ","firstName":" Иван ","lastName":" Петров ","password":"very-secure-password"
	}`, "")
	if response.StatusCode != http.StatusCreated {
		t.Fatalf("register status = %d, body = %s", response.StatusCode, readBody(t, response))
	}
	registeredCookie := response.Cookies()[0]
	user := decodeBody(t, response)
	assertUser(t, user, "user@example.com", "Иван", "Петров", "user")
	if _, exists := user["password"]; exists {
		t.Fatal("response contains password")
	}

	invalidRegistrations := []struct {
		name string
		body string
		code string
	}{
		{name: "empty first name", body: `{"email":"one@example.com","firstName":" ","lastName":"Петров","password":"very-secure-password"}`, code: "invalid_first_name"},
		{name: "long first name", body: `{"email":"two@example.com","firstName":"` + strings.Repeat("я", 101) + `","lastName":"Петров","password":"very-secure-password"}`, code: "invalid_first_name"},
		{name: "empty last name", body: `{"email":"three@example.com","firstName":"Иван","lastName":" ","password":"very-secure-password"}`, code: "invalid_last_name"},
		{name: "short password", body: `{"email":"four@example.com","firstName":"Иван","lastName":"Петров","password":"short"}`, code: "invalid_password"},
	}
	for _, test := range invalidRegistrations {
		t.Run(test.name, func(t *testing.T) {
			response := requestJSON(t, server.Client(), http.MethodPost, server.URL+"/api/auth/register", test.body, "")
			if response.StatusCode != http.StatusBadRequest || decodeBody(t, response)["code"] != test.code {
				t.Fatalf("unexpected validation response")
			}
		})
	}

	response = requestJSON(t, server.Client(), http.MethodPost, server.URL+"/api/auth/register", `{
		"email":"user@example.com","firstName":"Другой","lastName":"Пользователь","password":"very-secure-password"
	}`, "")
	if response.StatusCode != http.StatusConflict {
		t.Fatalf("duplicate status = %d, body = %s", response.StatusCode, readBody(t, response))
	}
	_ = response.Body.Close()

	wrongPassword := requestJSON(t, server.Client(), http.MethodPost, server.URL+"/api/auth/login", `{
		"email":"user@example.com","password":"wrong-password"
	}`, "")
	unknownEmail := requestJSON(t, server.Client(), http.MethodPost, server.URL+"/api/auth/login", `{
		"email":"missing@example.com","password":"wrong-password"
	}`, "")
	wrongBody := decodeBody(t, wrongPassword)
	unknownBody := decodeBody(t, unknownEmail)
	if wrongPassword.StatusCode != http.StatusUnauthorized || unknownEmail.StatusCode != http.StatusUnauthorized || wrongBody["code"] != unknownBody["code"] {
		t.Fatalf("credential errors differ: %v, %v", wrongBody, unknownBody)
	}

	response = requestJSON(t, server.Client(), http.MethodPost, server.URL+"/api/auth/login", `{
		"email":"user@example.com","password":"very-secure-password"
	}`, "")
	if response.StatusCode != http.StatusOK {
		t.Fatalf("login status = %d, body = %s", response.StatusCode, readBody(t, response))
	}
	loginCookie := response.Cookies()[0]
	assertUser(t, decodeBody(t, response), "user@example.com", "Иван", "Петров", "user")

	response = requestJSON(t, server.Client(), http.MethodGet, server.URL+"/api/auth/me", "", loginCookie.String())
	if response.StatusCode != http.StatusOK {
		t.Fatalf("current user status = %d, body = %s", response.StatusCode, readBody(t, response))
	}
	assertUser(t, decodeBody(t, response), "user@example.com", "Иван", "Петров", "user")

	now = now.Add(14 * 24 * time.Hour)
	response = requestJSON(t, server.Client(), http.MethodGet, server.URL+"/api/auth/me", "", loginCookie.String())
	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expired idle session status = %d", response.StatusCode)
	}
	_ = response.Body.Close()

	now = time.Date(2026, 9, 19, 12, 0, 0, 0, time.UTC)
	admin, err := service.CreateAdmin(ctx, auth.AdminInput{
		Email: "admin@example.com", FirstName: "Системный", LastName: "Администратор", Password: "very-secure-admin-password",
	})
	if err != nil || admin.Role != auth.RoleAdmin {
		t.Fatalf("create admin = %+v, %v", admin, err)
	}
	secondAdmin, err := service.CreateAdmin(ctx, auth.AdminInput{
		Email: "another-admin@example.com", FirstName: "Другой", LastName: "Администратор", Password: "another-secure-password",
	})
	if err != nil || secondAdmin.Role != auth.RoleAdmin {
		t.Fatalf("create second admin = %+v, %v", secondAdmin, err)
	}
	response = requestJSON(t, server.Client(), http.MethodPost, server.URL+"/api/auth/login", `{
		"email":"admin@example.com","password":"very-secure-admin-password"
	}`, "")
	if response.StatusCode != http.StatusOK {
		t.Fatalf("admin login status = %d, body = %s", response.StatusCode, readBody(t, response))
	}
	assertUser(t, decodeBody(t, response), "admin@example.com", "Системный", "Администратор", "admin")

	response = requestJSON(t, server.Client(), http.MethodGet, server.URL+"/api/auth/me", "", registeredCookie.String())
	_ = response.Body.Close()

	response = requestJSON(t, server.Client(), http.MethodPost, server.URL+"/auth/register", `{}`, "")
	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("legacy auth route status = %d", response.StatusCode)
	}
	_ = response.Body.Close()
}

func requestJSON(t *testing.T, client *http.Client, method, url, body, cookie string) *http.Response {
	t.Helper()
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	if method != http.MethodGet {
		request.Header.Set("Origin", testOrigin)
	}
	if cookie != "" {
		request.Header.Set("Cookie", cookie)
	}
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	return response
}

func decodeBody(t *testing.T, response *http.Response) map[string]any {
	t.Helper()
	defer response.Body.Close()
	var value map[string]any
	if err := json.NewDecoder(response.Body).Decode(&value); err != nil {
		t.Fatal(err)
	}
	return value
}

func readBody(t *testing.T, response *http.Response) string {
	t.Helper()
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(body)
}

func assertUser(t *testing.T, user map[string]any, email, firstName, lastName, role string) {
	t.Helper()
	if user["email"] != email || user["firstName"] != firstName || user["lastName"] != lastName || user["role"] != role {
		t.Fatalf("unexpected user: %#v", user)
	}
}
