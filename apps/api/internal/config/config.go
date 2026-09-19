package config

import (
	"fmt"
	"net/netip"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	APIAddr           string
	DatabaseURL       string
	WebOrigin         string
	CookieSecure      bool
	TrustedProxyCIDRs []netip.Prefix
	ShutdownTimeout   time.Duration
}

func Load() (Config, error) {
	config := Config{
		APIAddr:         valueOrDefault("API_ADDR", ":8080"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		WebOrigin:       valueOrDefault("WEB_ORIGIN", "http://localhost:5173"),
		ShutdownTimeout: 10 * time.Second,
	}
	if config.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	if _, err := url.ParseRequestURI(config.WebOrigin); err != nil {
		return Config{}, fmt.Errorf("WEB_ORIGIN is invalid: %w", err)
	}
	secure, err := strconv.ParseBool(valueOrDefault("COOKIE_SECURE", "true"))
	if err != nil {
		return Config{}, fmt.Errorf("COOKIE_SECURE must be true or false")
	}
	config.CookieSecure = secure
	trustedProxyCIDRs, err := parseCIDRs(os.Getenv("TRUSTED_PROXY_CIDRS"))
	if err != nil {
		return Config{}, err
	}
	config.TrustedProxyCIDRs = trustedProxyCIDRs
	return config, nil
}

func parseCIDRs(value string) ([]netip.Prefix, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	values := strings.Split(value, ",")
	prefixes := make([]netip.Prefix, 0, len(values))
	for _, value := range values {
		prefix, err := netip.ParsePrefix(strings.TrimSpace(value))
		if err != nil {
			return nil, fmt.Errorf("TRUSTED_PROXY_CIDRS contains invalid CIDR %q: %w", value, err)
		}
		prefixes = append(prefixes, prefix.Masked())
	}
	return prefixes, nil
}

func valueOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
