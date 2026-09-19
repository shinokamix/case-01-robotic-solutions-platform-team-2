package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/netip"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/httpapi/generated"
)

func NewRouter(handler generated.StrictServerInterface, logger *slog.Logger, webOrigin string, trustedProxyCIDRs ...netip.Prefix) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(15 * time.Second))
	router.Use(sessionCookie)
	router.Use(validateJSONBody)
	router.Use(func(next http.Handler) http.Handler { return cors(webOrigin, next) })
	router.Use(func(next http.Handler) http.Handler { return csrfProtection(webOrigin, next) })

	rateLimiter := NewRateLimiter(10, 15*time.Minute, NewClientIPResolver(trustedProxyCIDRs))
	strict := generated.NewStrictHandlerWithOptions(handler, []generated.StrictMiddlewareFunc{rateLimiter.StrictMiddleware}, generated.StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(response http.ResponseWriter, _ *http.Request, _ error) {
			writeJSONError(response, http.StatusBadRequest, "invalid_request", "Запрос не прошел проверку.")
		},
		ResponseErrorHandlerFunc: func(response http.ResponseWriter, request *http.Request, err error) {
			logger.ErrorContext(request.Context(), "write API response", "error", err)
			writeJSONError(response, http.StatusInternalServerError, "internal_error", "Внутренняя ошибка сервера.")
		},
	})
	return generated.HandlerFromMuxWithBaseURL(strict, router, "/api")
}

func JSON(response http.ResponseWriter, status int, value any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	_ = json.NewEncoder(response).Encode(value)
}
