package agent

import (
	"net"
	"net/url"
	"strings"
)

// Return only fixed diagnostics; never echo a rejected server value.
func validatedPublicURL(value any, token string, serverURL string, environment string) (string, string) {
	raw, ok := value.(string)
	if !ok || strings.TrimSpace(raw) == "" {
		return "", "server acknowledgement must include a nonempty string field named url"
	}
	if token != "" && strings.Contains(raw, token) {
		return "", "URL contains authentication data"
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", "server returned malformed URL syntax"
	}
	server, serverErr := url.Parse(serverURL)
	localBackend := serverErr == nil && (server.Scheme == "http" || server.Scheme == "https") && server.User == nil && isLocalHost(server.Hostname())
	localURL := environment == "development" && localBackend && isLocalHost(u.Hostname())
	if localURL {
		if u.Scheme != "http" && u.Scheme != "https" {
			return "", "local development URL must use HTTP or HTTPS"
		}
	} else {
		if u.Scheme != "https" {
			return "", "server must return an HTTPS URL; HTTP localhost URLs require ENV_ENVIRONMENT=development and a local backend"
		}
		if !strings.HasSuffix(strings.ToLower(u.Hostname()), ".tunnexo.live") {
			return "", "hostname must be a subdomain of tunnexo.live; localhost URLs require ENV_ENVIRONMENT=development and a local backend"
		}
	}
	if u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" {
		return "", "URL must not contain credentials, query parameters or a fragment"
	}
	if u.Path != "" && u.Path != "/" {
		return "", "server must return the tunnel root URL without an application path"
	}
	return u.String(), ""
}

// No DNS lookup: only explicit localhost names and loopback IPs are local.
func isLocalHost(host string) bool {
	host = strings.ToLower(host)
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
