package agent

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGuestRegistration(t *testing.T) {
	secret := strings.Repeat("secret", 8)
	for _, tc := range []struct {
		name   string
		status int
		body   string
		want   bool
	}{
		{"success", 201, fmt.Sprintf(`{"token":%q,"expiresAt":%q}`, secret, time.Now().Add(time.Hour).Format(time.RFC3339)), true},
		{"missing backend", 404, secret, false},
		{"rate limit", 429, secret, false},
		{"server error", 500, secret, false},
		{"malformed", 201, secret, false},
		{"expired", 201, fmt.Sprintf(`{"token":%q,"expiresAt":"2020-01-01T00:00:00Z"}`, secret), false},
		{"oversized", 201, strings.Repeat(secret, 400), false},
		{"redirect", 302, secret, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" || r.URL.Path != "/api/v1/agent/guest" || r.Header.Get("Authorization") != "" {
					t.Error("unexpected registration request")
				}
				w.Header().Set("Location", "https://example.invalid/")
				w.WriteHeader(tc.status)
				fmt.Fprint(w, tc.body)
			}))
			defer srv.Close()
			token, err := requestGuestToken(context.Background(), srv.URL+"/tunnel", srv.Client())
			if tc.want {
				if err != nil || token != secret {
					t.Fatal("registration failed")
				}
			} else {
				if err == nil || token != "" {
					t.Fatal("expected failure")
				}
				if strings.Contains(err.Error(), secret) {
					t.Fatal("credential exposed in error")
				}
			}
		})
	}
}

func TestGuestRejectsInsecureOrCredentialURLs(t *testing.T) {
	for _, server := range []string{"http://example.com/tunnel", "https://user:secret@example.com/tunnel", "https://example.com/tunnel?token=secret"} {
		_, err := requestGuestToken(context.Background(), server, &http.Client{})
		if err == nil || strings.Contains(err.Error(), "secret") {
			t.Fatal("unsafe server accepted or exposed")
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := requestGuestToken(ctx, "https://example.invalid/tunnel", &http.Client{})
	if err == nil {
		t.Fatal("cancellation ignored")
	}
}
