package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/netip"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/auth"
	"github.com/shinokamix/case-01-robotic-solutions-platform-team-2/apps/api/internal/httpapi/generated"
)

type sessionTokenKey struct{}

func sessionCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie("session_id")
		if err == nil {
			request = request.WithContext(contextWithSessionToken(request, cookie.Value))
		}
		next.ServeHTTP(response, request)
	})
}

func contextWithSessionToken(request *http.Request, token string) context.Context {
	return context.WithValue(request.Context(), sessionTokenKey{}, token)
}

func validateJSONBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost || request.URL.Path != "/api/auth/register" && request.URL.Path != "/api/auth/login" {
			next.ServeHTTP(response, request)
			return
		}

		body, err := io.ReadAll(io.LimitReader(request.Body, 1<<20+1))
		if err != nil || len(body) > 1<<20 {
			writeJSONError(response, http.StatusBadRequest, "invalid_request", "Запрос не прошел проверку.")
			return
		}
		var target any
		if request.URL.Path == "/api/auth/register" {
			target = &generated.RegisterRequest{}
		} else {
			target = &generated.LoginRequest{}
		}
		decoder := json.NewDecoder(bytes.NewReader(body))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(target); err != nil {
			writeJSONError(response, http.StatusBadRequest, "invalid_request", "Запрос не прошел проверку.")
			return
		}
		if err := decoder.Decode(&struct{}{}); err != io.EOF {
			writeJSONError(response, http.StatusBadRequest, "invalid_request", "Запрос не прошел проверку.")
			return
		}
		request.Body = io.NopCloser(bytes.NewReader(body))
		next.ServeHTTP(response, request)
	})
}

func csrfProtection(allowedOrigin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost || request.Method == http.MethodPut || request.Method == http.MethodPatch || request.Method == http.MethodDelete {
			if request.Header.Get("Origin") != allowedOrigin {
				writeJSONError(response, http.StatusForbidden, "csrf_failed", "Источник запроса не разрешен.")
				return
			}
		}
		next.ServeHTTP(response, request)
	})
}

func cors(allowedOrigin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Origin") == allowedOrigin {
			response.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			response.Header().Set("Access-Control-Allow-Credentials", "true")
			response.Header().Set("Vary", "Origin")
		}
		if request.Method == http.MethodOptions {
			response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			response.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(response, request)
	})
}

type ClientIPResolver struct {
	trustedProxyCIDRs []netip.Prefix
}

func NewClientIPResolver(trustedProxyCIDRs []netip.Prefix) ClientIPResolver {
	return ClientIPResolver{trustedProxyCIDRs: trustedProxyCIDRs}
}

func (resolver ClientIPResolver) Resolve(request *http.Request) string {
	remoteIP, err := parseIP(request.RemoteAddr)
	if err != nil || !resolver.trusted(remoteIP) {
		return remoteIP.String()
	}

	forwarded := strings.Split(request.Header.Get("X-Forwarded-For"), ",")
	if len(forwarded) == 1 && strings.TrimSpace(forwarded[0]) == "" {
		return remoteIP.String()
	}
	addresses := make([]netip.Addr, 0, len(forwarded))
	for _, value := range forwarded {
		address, err := netip.ParseAddr(strings.TrimSpace(value))
		if err != nil {
			return remoteIP.String()
		}
		addresses = append(addresses, address.Unmap())
	}
	for index := len(addresses) - 1; index >= 0; index-- {
		if !resolver.trusted(addresses[index]) {
			return addresses[index].String()
		}
	}
	return addresses[0].String()
}

func (resolver ClientIPResolver) trusted(address netip.Addr) bool {
	for _, prefix := range resolver.trustedProxyCIDRs {
		if prefix.Contains(address) {
			return true
		}
	}
	return false
}

func parseIP(remoteAddress string) (netip.Addr, error) {
	if addressPort, err := netip.ParseAddrPort(remoteAddress); err == nil {
		return addressPort.Addr().Unmap(), nil
	}
	address, err := netip.ParseAddr(remoteAddress)
	return address.Unmap(), err
}

type rateBucket struct {
	started time.Time
	count   int
}

type RateLimiter struct {
	mu          sync.Mutex
	buckets     map[string]rateBucket
	now         func() time.Time
	limit       int
	window      time.Duration
	clientIP    ClientIPResolver
	lastCleanup time.Time
}

func NewRateLimiter(limit int, window time.Duration, clientIP ClientIPResolver) *RateLimiter {
	return &RateLimiter{
		buckets:  make(map[string]rateBucket),
		now:      time.Now,
		limit:    limit,
		window:   window,
		clientIP: clientIP,
	}
}

func (limiter *RateLimiter) StrictMiddleware(next generated.StrictHandlerFunc, operationID string) generated.StrictHandlerFunc {
	return func(ctx context.Context, response http.ResponseWriter, request *http.Request, value any) (any, error) {
		key, limited := limiter.key(operationID, request, value)
		if !limited {
			return next(ctx, response, request, value)
		}
		allowed, retryAfter := limiter.allow(key)
		if !allowed {
			seconds := max(1, int(retryAfter.Round(time.Second)/time.Second))
			response.Header().Set("Retry-After", strconv.Itoa(seconds))
			writeJSONError(response, http.StatusTooManyRequests, "rate_limited", "Слишком много запросов. Повторите попытку позже.")
			return nil, nil
		}
		return next(ctx, response, request, value)
	}
}

func (limiter *RateLimiter) key(operationID string, request *http.Request, value any) (string, bool) {
	clientIP := limiter.clientIP.Resolve(request)
	switch operationID {
	case "Register":
		return operationID + ":" + clientIP, true
	case "Login":
		login, ok := value.(generated.LoginRequestObject)
		if !ok || login.Body == nil {
			return operationID + ":" + clientIP, true
		}
		email, err := auth.CanonicalEmail(string(login.Body.Email))
		if err != nil {
			email = strings.ToLower(strings.TrimSpace(string(login.Body.Email)))
		}
		return operationID + ":" + clientIP + ":" + email, true
	default:
		return "", false
	}
}

func (limiter *RateLimiter) allow(key string) (bool, time.Duration) {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	now := limiter.now()
	if limiter.lastCleanup.IsZero() || now.Sub(limiter.lastCleanup) >= limiter.window {
		for key, bucket := range limiter.buckets {
			if now.Sub(bucket.started) >= limiter.window {
				delete(limiter.buckets, key)
			}
		}
		limiter.lastCleanup = now
	}

	bucket := limiter.buckets[key]
	if bucket.started.IsZero() || now.Sub(bucket.started) >= limiter.window {
		limiter.buckets[key] = rateBucket{started: now, count: 1}
		return true, 0
	}
	if bucket.count >= limiter.limit {
		return false, limiter.window - now.Sub(bucket.started)
	}
	bucket.count++
	limiter.buckets[key] = bucket
	return true, 0
}

func writeJSONError(response http.ResponseWriter, status int, code, message string) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	_ = json.NewEncoder(response).Encode(generated.ErrorResponse{Code: code, Message: message})
}
