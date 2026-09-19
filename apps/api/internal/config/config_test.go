package config

import "testing"

func TestLoadDefaultsCookieToSecure(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/robotics")
	t.Setenv("COOKIE_SECURE", "")
	t.Setenv("TRUSTED_PROXY_CIDRS", "")

	config, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if !config.CookieSecure {
		t.Fatal("COOKIE_SECURE must default to true")
	}
}

func TestLoadParsesTrustedProxyCIDRs(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/robotics")
	t.Setenv("COOKIE_SECURE", "false")
	t.Setenv("TRUSTED_PROXY_CIDRS", "10.0.0.0/8, 192.0.2.0/24")

	config, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(config.TrustedProxyCIDRs) != 2 {
		t.Fatalf("trusted proxy CIDRs = %d, want 2", len(config.TrustedProxyCIDRs))
	}
}
