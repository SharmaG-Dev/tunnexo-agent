package agent

import (
	"strings"
	"testing"
)

func TestValidatedPublicURL(t *testing.T) {
	for _, tc := range []struct {
		value any
		token string
		valid bool
	}{
		{"https://demo.tunnexo.live", "", true},
		{"https://demo.tunnexo.live/", "secret", true},
		{nil, "secret", false},
		{"http://demo.tunnexo.live", "secret", false},
		{"https://demo.old.example", "secret", false},
		{"https://demo.tunnexo.live?token=secret", "secret", false},
		{"https://user:secret@demo.tunnexo.live", "secret", false},
		{"https://demo.tunnexo.live/path", "secret", false},
	} {
		got, reason := validatedPublicURL(tc.value, tc.token, DefaultServerURL, "production")
		if (reason == "") != tc.valid {
			t.Fatalf("unexpected validation outcome: %s", reason)
		}
		if !tc.valid && (got != "" || strings.Contains(reason, "secret")) {
			t.Fatal("rejected data exposed")
		}
	}
}

func TestLocalPublicURLs(t *testing.T) {
	for _, tc := range []struct {
		backend, public string
		valid           bool
	}{
		{"http://localhost:3000/tunnel", "http://demo.localhost:3000", true},
		{"http://127.0.0.1:3000/tunnel", "http://localhost:3000/", true},
		{"http://[::1]:3000/tunnel", "http://[::1]:3000", true},
		{"http://localhost:3000/tunnel", "https://demo.tunnexo.live", true},
		{DefaultServerURL, "http://demo.localhost:3000", false},
		{DefaultServerURL, "https://localhost:3000", false},
		{"http://localhost.evil.example/tunnel", "http://localhost:3000", false},
		{"http://localhost:3000/tunnel", "http://evil.example", false},
		{"http://localhost:3000/tunnel", "http://localhost.evil.example", false},
		{"http://localhost:3000/tunnel", "http://localhost:3000?token=secret", false},
		{"http://localhost:3000/tunnel", "http://user:secret@localhost:3000", false},
		{"http://localhost:3000/tunnel", "http://localhost:3000/path", false},
	} {
		got, reason := validatedPublicURL(tc.public, "secret", tc.backend, "development")
		if (reason == "") != tc.valid {
			t.Errorf("backend %s URL %s: %s", tc.backend, tc.public, reason)
		}
		if !tc.valid && (got != "" || strings.Contains(reason, "secret")) {
			t.Fatal("rejected value exposed")
		}
	}
}

func TestProductionRejectsLocalURL(t *testing.T) {
	for _, mode := range []string{"production", "", "invalid"} {
		if _, reason := validatedPublicURL("http://demo.localhost:3000", "", "http://localhost:3000/tunnel", mode); reason == "" {
			t.Fatal("local URL accepted outside development")
		}
	}
}
