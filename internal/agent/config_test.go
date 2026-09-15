package agent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadEnv(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("# comment\nexport A=hello # comment\nB='token#literal'\nC=\"quoted value\" # comment\n"), 0600); err != nil {
		t.Fatal(err)
	}
	values, err := readEnv(path)
	if err != nil {
		t.Fatal(err)
	}
	if values["A"] != "hello" || values["B"] != "token#literal" || values["C"] != "quoted value" {
		t.Fatalf("unexpected values: %v", values)
	}
	if err := os.WriteFile(path, []byte("TOKEN='private-unclosed"), 0600); err != nil {
		t.Fatal(err)
	}
	_, err = readEnv(path)
	if err == nil || strings.Contains(err.Error(), "private") {
		t.Fatal("expected redacted parsing error")
	}
}

func TestConfigValidationAndEnvironmentPrecedence(t *testing.T) {
	t.Setenv("TUNNEXO_SERVER_URL", "https://example.test/tunnel")
	t.Setenv("TUNNEXO_AGENT_TOKEN", "")
	t.Setenv("TUNNEXO_AGENT_NAME_PREFIX", "demo")
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("TUNNEXO_AGENT_TOKEN=file-token\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if got := configValue(map[string]string{"TUNNEXO_AGENT_TOKEN": "file-token"}, "TUNNEXO_AGENT_TOKEN"); got != "" {
		t.Fatal("empty environment must select guest mode")
	}
	t.Setenv("TUNNEXO_AGENT_TOKEN", "test-token")
	for _, target := range []string{"", "localhost:4000", "ftp://localhost", "http://user:password@localhost", "http://localhost/#fragment", "http://localhost/?q=1"} {
		if _, err := LoadConfig(path, target); err == nil {
			t.Errorf("accepted invalid target %q", target)
		}
	}
}

func TestConfigValueMigration(t *testing.T) {
	const key = "TUNNEXO_MIGRATION_TEST"
	const legacy = "PORTUNE_MIGRATION_TEST"
	cases := []struct {
		name    string
		values  map[string]string
		current *string
		old     *string
		want    string
	}{
		{name: "missing"},
		{name: "legacy file", values: map[string]string{legacy: " old "}, want: "old"},
		{name: "legacy environment overrides file", values: map[string]string{legacy: "file"}, old: stringPointer("env"), want: "env"},
		{name: "new file overrides legacy environment", values: map[string]string{key: "new"}, old: stringPointer("old"), want: "new"},
		{name: "new environment overrides file", values: map[string]string{key: "file"}, current: stringPointer("env"), want: "env"},
		{name: "empty new environment prevents fallback", values: map[string]string{key: "file", legacy: "old"}, current: stringPointer("")},
		{name: "empty new file prevents fallback", values: map[string]string{key: ""}, old: stringPointer("old")},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for _, name := range []string{key, legacy} {
				t.Setenv(name, "")
				if err := os.Unsetenv(name); err != nil {
					t.Fatal(err)
				}
			}
			if tc.current != nil {
				t.Setenv(key, *tc.current)
			}
			if tc.old != nil {
				t.Setenv(legacy, *tc.old)
			}
			if got := configValue(tc.values, key); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func stringPointer(value string) *string { return &value }

func TestPublicDefaultsAndOverrides(t *testing.T) {
	for _, key := range []string{"TUNNEXO_SERVER_URL", "PORTUNE_SERVER_URL", "TUNNEXO_AGENT_NAME_PREFIX", "PORTUNE_AGENT_NAME_PREFIX", "TUNNEXO_AGENT_TOKEN", "PORTUNE_AGENT_TOKEN"} {
		t.Setenv(key, "")
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
	}
	if got := configValueOrDefault(nil, "TUNNEXO_SERVER_URL", DefaultServerURL); got != "https://tunnexo.live/tunnel" {
		t.Fatalf("unexpected public server: %q", got)
	}
	if got := configValueOrDefault(nil, "TUNNEXO_AGENT_NAME_PREFIX", DefaultAgentNamePrefix); got != "tunnexo" {
		t.Fatalf("unexpected prefix: %q", got)
	}
	values := map[string]string{"TUNNEXO_SERVER_URL": "https://custom.example/tunnel"}
	if got := configValueOrDefault(values, "TUNNEXO_SERVER_URL", DefaultServerURL); got != values["TUNNEXO_SERVER_URL"] {
		t.Fatalf("custom server ignored: %q", got)
	}
	t.Setenv("TUNNEXO_SERVER_URL", "")
	if got := configValueOrDefault(values, "TUNNEXO_SERVER_URL", DefaultServerURL); got != DefaultServerURL {
		t.Fatalf("empty override should select public default: %q", got)
	}
}

func TestInvalidEnvironment(t *testing.T) {
	t.Setenv("ENV_ENVIRONMENT", "invalid")
	_, err := LoadConfig(filepath.Join(t.TempDir(), "missing.env"), "http://localhost:4000")
	if err == nil || !strings.Contains(err.Error(), "ENV_ENVIRONMENT") {
		t.Fatal("invalid environment accepted")
	}
}
