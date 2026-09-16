package agent

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

const DefaultServerURL = "https://tunnexo.live/tunnel"
const DefaultAgentNamePrefix = "tunnexo"

func LoadConfig(path, target string) (Config, error) {
	values, err := readEnv(path)
	if err != nil {
		return Config{}, err
	}
	get := func(key string) string {
		return configValue(values, key)
	}
	environment := get("ENV_ENVIRONMENT")
	if environment == "" {
		environment = "production"
	}
	if environment != "production" && environment != "development" {
		return Config{}, fmt.Errorf("ENV_ENVIRONMENT must be production or development")
	}
	defaultServer := DefaultServerURL
	if environment == "development" {
		defaultServer = "http://localhost:3000/tunnel"
	}
	config := Config{
		Environment: environment,
		ServerURL:   configValueOrDefault(values, "TUNNEXO_SERVER_URL", defaultServer),
		Token:       get("TUNNEXO_AGENT_TOKEN"),
		LocalTarget: strings.TrimSpace(target),
	}
	prefix := configValueOrDefault(values, "TUNNEXO_AGENT_NAME_PREFIX", DefaultAgentNamePrefix)

	if _, err := validateURL(config.ServerURL); err != nil {
		return Config{}, fmt.Errorf("TUNNEXO_SERVER_URL must be an absolute HTTP(S) URL without credentials or fragment")
	}
	u, err := validateURL(config.LocalTarget)
	if err != nil {
		return Config{}, fmt.Errorf("--target must be an absolute HTTP(S) URL without credentials or fragment")
	}
	if u.RawQuery != "" || u.ForceQuery {
		return Config{}, fmt.Errorf("--target must not contain a query string")
	}
	for _, c := range prefix {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-') {
			return Config{}, fmt.Errorf("TUNNEXO_AGENT_NAME_PREFIX must contain only letters, digits or hyphens")
		}
	}
	config.AgentIP, err = localIP()
	if err != nil {
		return Config{}, err
	}
	config.AgentName = config.AgentIP + "-" + prefix + "-agent"
	config.TunnelName = u.Hostname()
	return config, nil
}

// New names take precedence over legacy names; environment overrides files
// for each name. Explicit empty values are preserved for validation.
func configValue(values map[string]string, key string) string {
	legacy := strings.Replace(key, "TUNNEXO_", "PORTUNE_", 1)
	for _, name := range []string{key, legacy} {
		if value, ok := os.LookupEnv(name); ok {
			return strings.TrimSpace(value)
		}
		if value, ok := values[name]; ok {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func configValueOrDefault(values map[string]string, key, fallback string) string {
	if value := configValue(values, key); value != "" {
		return value
	}
	return fallback
}

func validateURL(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	if (u.Scheme != "http" && u.Scheme != "https") || u.Hostname() == "" || u.User != nil || u.Fragment != "" {
		return nil, fmt.Errorf("invalid URL")
	}
	return u, nil
}

func localIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("detect local IP: %w", err)
	}
	var fallback string
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addresses, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, address := range addresses {
			ip, _, err := net.ParseCIDR(address.String())
			if err != nil || !ip.IsGlobalUnicast() {
				continue
			}
			if ip.To4() != nil && ip.IsPrivate() {
				return ip.String(), nil
			}
			if fallback == "" || ip.To4() != nil {
				fallback = ip.String()
			}
		}
	}
	if fallback != "" {
		return fallback, nil
	}
	return "", fmt.Errorf("cannot detect a non-loopback local IP; connect to a network before starting the agent")
}

func readEnv(path string) (map[string]string, error) {
	values := make(map[string]string)
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return values, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read .env: %w", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	line := 0
	for scanner.Scan() {
		line++
		text := strings.TrimSpace(scanner.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		text = strings.TrimSpace(strings.TrimPrefix(text, "export "))
		key, value, ok := strings.Cut(text, "=")
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		if !ok || key == "" {
			return nil, fmt.Errorf("invalid .env entry at line %d", line)
		}
		if strings.HasPrefix(value, "\"") || strings.HasPrefix(value, "'") {
			quote := value[0]
			end := strings.IndexByte(value[1:], quote)
			if end < 0 {
				return nil, fmt.Errorf("unclosed .env quote at line %d", line)
			}
			end++
			tail := strings.TrimSpace(value[end+1:])
			if tail != "" && !strings.HasPrefix(tail, "#") {
				return nil, fmt.Errorf("invalid .env value at line %d", line)
			}
			value = value[1:end]
		} else if index := strings.Index(value, " #"); index >= 0 {
			value = strings.TrimSpace(value[:index])
		}
		values[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read .env: %w", err)
	}
	return values, nil
}
