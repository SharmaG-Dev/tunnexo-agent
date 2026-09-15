package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"
)

// Guest credentials are session-only and never written to disk or logged.
func requestGuestToken(ctx context.Context, server string, client *http.Client) (string, error) {
	u, err := url.Parse(server)
	if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return "", fmt.Errorf("guest registration requires an HTTPS server URL without credentials, query or fragment")
	}
	u.Path, u.RawPath = "/api/v1/agent/guest", ""
	body, _ := json.Marshal(map[string]string{"platform": runtime.GOOS, "architecture": runtime.GOARCH})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("could not prepare guest registration")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	// Never follow redirects carrying registration requests to another endpoint.
	bounded := *client
	bounded.Timeout = 15 * time.Second
	bounded.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	resp, err := bounded.Do(req)
	if err != nil {
		return "", fmt.Errorf("guest registration connection failed; check network and server availability")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		switch resp.StatusCode {
		case http.StatusNotFound, http.StatusMethodNotAllowed:
			return "", fmt.Errorf("guest registration is unavailable; the server must deploy POST /api/v1/agent/guest")
		case http.StatusTooManyRequests:
			return "", fmt.Errorf("guest registration limit reached; try again later")
		case http.StatusUnauthorized, http.StatusForbidden:
			return "", fmt.Errorf("guest registration denied by server")
		default:
			return "", fmt.Errorf("guest registration failed (HTTP %d)", resp.StatusCode)
		}
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, 16385))
	if err != nil || len(data) > 16384 {
		return "", fmt.Errorf("invalid guest registration response")
	}
	var result struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expiresAt"`
	}
	if json.Unmarshal(data, &result) != nil || len(result.Token) < 32 || len(result.Token) > 4096 || strings.ContainsAny(result.Token, " \t\r\n") || !result.ExpiresAt.After(time.Now().Add(time.Minute)) {
		return "", fmt.Errorf("invalid guest registration response")
	}
	return result.Token, nil
}
