package httpapi

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"testing"
	"time"

	"github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/auth"
	"github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/httpapi/generated"
)

type fakeAuthService struct {
	registered auth.Authenticated
	login      auth.Authenticated
	current    auth.User
	loginErr   error
}

func (service *fakeAuthService) Register(_ context.Context, input auth.RegisterInput) (auth.Authenticated, error) {
	if input.Email == "taken@example.com" {
		return auth.Authenticated{}, auth.ErrEmailAlreadyExists
	}
	return service.registered, nil
}

func (service *fakeAuthService) Login(context.Context, auth.LoginInput) (auth.Authenticated, error) {
	return service.login, service.loginErr
}

func (service *fakeAuthService) CurrentUser(_ context.Context, token string) (auth.User, error) {
	if token != "session-token" {
		return auth.User{}, auth.ErrUnauthenticated
	}
	return service.current, nil
}

func TestRegisterSetsSecureOpaqueCookie(t *testing.T) {
	user := auth.User{ID: "01K5", Email: "user@example.com", FirstName: "Иван", LastName: "Петров", Role: auth.RoleUser}
	service := &fakeAuthService{registered: auth.Authenticated{User: user, SessionToken: "session-token", ExpiresAt: time.Now().Add(30 * 24 * time.Hour)}}
	router := NewRouter(NewAuthHandler(service, testLogger(), true), testLogger(), "http://localhost:5173")

	request := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{"email":"user@example.com","firstName":"Иван","lastName":"Петров","password":"very-secure-password"}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Origin", "http://localhost:5173")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	cookie := response.Result().Cookies()[0]
	if cookie.Name != "session_id" || cookie.Value != "session-token" || !cookie.HttpOnly || !cookie.Secure || cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("unexpected cookie: %+v", cookie)
	}
	var got generated.User
	if err := json.NewDecoder(response.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Id != user.ID || got.Role == nil || *got.Role != generated.UserRoleUser {
		t.Fatalf("unexpected user response: %+v", got)
	}
}

func TestRegisterRejectsRoleField(t *testing.T) {
	service := &fakeAuthService{}
	router := NewRouter(NewAuthHandler(service, testLogger(), false), testLogger(), "http://localhost:5173")
	request := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{"email":"user@example.com","firstName":"Иван","lastName":"Петров","password":"very-secure-password","role":"admin"}`))
	request.Header.Set("Origin", "http://localhost:5173")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestUnsafeRequestRequiresAllowedOrigin(t *testing.T) {
	router := NewRouter(NewAuthHandler(&fakeAuthService{}, testLogger(), false), testLogger(), "http://localhost:5173")
	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"email":"user@example.com","password":"very-secure-password"}`))
	request.Header.Set("Origin", "https://attacker.example")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d", response.Code)
	}
}

func TestLoginUsesGenericCredentialError(t *testing.T) {
	service := &fakeAuthService{loginErr: auth.ErrInvalidCredentials}
	router := NewRouter(NewAuthHandler(service, testLogger(), false), testLogger(), "http://localhost:5173")
	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"email":"unknown@example.com","password":"very-secure-password"}`))
	request.Header.Set("Origin", "http://localhost:5173")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized || !strings.Contains(response.Body.String(), "invalid_credentials") {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestCurrentUserReadsSessionCookie(t *testing.T) {
	user := auth.User{ID: "01K5", Email: "user@example.com", FirstName: "Иван", LastName: "Петров", Role: auth.RoleAdmin}
	router := NewRouter(NewAuthHandler(&fakeAuthService{current: user}, testLogger(), false), testLogger(), "http://localhost:5173")
	request := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	request.AddCookie(&http.Cookie{Name: "session_id", Value: "session-token"})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"role":"admin"`) {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestCurrentUserWithoutSessionIsUnauthorized(t *testing.T) {
	router := NewRouter(NewAuthHandler(&fakeAuthService{}, testLogger(), false), testLogger(), "http://localhost:5173")
	request := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized || !strings.Contains(response.Body.String(), "unauthenticated") {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(2, time.Minute, NewClientIPResolver(nil))
	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	request.RemoteAddr = "192.0.2.10:1234"
	value := generated.LoginRequestObject{Body: &generated.LoginRequest{Email: "user@example.com"}}
	handler := limiter.StrictMiddleware(
		func(_ context.Context, response http.ResponseWriter, _ *http.Request, _ any) (any, error) {
			response.WriteHeader(http.StatusNoContent)
			return nil, nil
		},
		"Login",
	)
	for index, expected := range []int{http.StatusNoContent, http.StatusNoContent, http.StatusTooManyRequests} {
		response := httptest.NewRecorder()
		_, _ = handler(request.Context(), response, request, value)
		if response.Code != expected {
			t.Fatalf("request %d status = %d, want %d", index+1, response.Code, expected)
		}
	}
}

func TestClientIPResolverTrustsOnlyConfiguredProxies(t *testing.T) {
	resolver := NewClientIPResolver([]netip.Prefix{netip.MustParsePrefix("10.0.0.0/8")})

	trusted := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	trusted.RemoteAddr = "10.0.0.2:8080"
	trusted.Header.Set("X-Forwarded-For", "198.51.100.10")
	if got := resolver.Resolve(trusted); got != "198.51.100.10" {
		t.Fatalf("trusted proxy IP = %s", got)
	}

	direct := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	direct.RemoteAddr = "192.0.2.10:8080"
	direct.Header.Set("X-Forwarded-For", "198.51.100.10")
	if got := resolver.Resolve(direct); got != "192.0.2.10" {
		t.Fatalf("direct client IP = %s", got)
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

var _ AuthService = (*fakeAuthService)(nil)
