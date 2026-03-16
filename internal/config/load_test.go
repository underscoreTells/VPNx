package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()

	jsonBytes, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}

	return jsonBytes
}

func validConfigInput() map[string]any {
	return map[string]any{
		"schema_version":  DEFAULT_CONFIG_VERSION,
		"gluetun_version": GLUETUN_TARGET_VERSION,
		"vpn_config": map[string]any{
			"provider": "custom",
			"protocol": ProtocolWireguard,
			"openvpn_credentials": map[string]any{
				"username": map[string]any{"from": "env", "name": "VPN_USERNAME"},
				"password": map[string]any{"from": "env", "name": "VPN_PASSWORD"},
			},
			"wireguard_credentials": map[string]any{
				"private_key":   map[string]any{"from": "env", "name": "WIREGUARD_PRIVATE_KEY"},
				"addresses":     []string{""},
				"public_key":    map[string]any{"from": "env", "name": "WIREGUARD_PUBLIC_KEY"},
				"endpoint_ip":   "",
				"endpoint_port": 3000,
			},
		},
		"log": map[string]any{
			"level":       LOG_LEVEL_INFO,
			"destination": "stdout",
			"filename":    "vpnx.log",
		},
		"env_vars": []string{},
	}
}

func configJSON(t *testing.T, mutate func(map[string]any)) []byte {
	t.Helper()

	config := validConfigInput()
	if mutate != nil {
		mutate(config)
	}

	return mustJSON(t, config)
}

func nestedMap(t *testing.T, parent map[string]any, key string) map[string]any {
	t.Helper()

	value, ok := parent[key]
	if !ok {
		t.Fatalf("missing map key %q", key)
	}

	nested, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("map key %q had type %T, want map[string]any", key, value)
	}

	return nested
}

func joinedErrors(errs []error) string {
	parts := make([]string, 0, len(errs))
	for _, err := range errs {
		parts = append(parts, err.Error())
	}

	return strings.Join(parts, "; ")
}

func requireErrorContains(t *testing.T, errs []error, want string) {
	t.Helper()

	if len(errs) == 0 {
		t.Fatalf("expected error containing %q, got none", want)
	}

	got := joinedErrors(errs)
	if !strings.Contains(got, want) {
		t.Fatalf("expected errors to contain %q, got %q", want, got)
	}
}

func TestLoadAppConfigFromBytes(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config, errs := LoadAppConfigFromBytes(configJSON(t, nil))
		if len(errs) > 0 {
			t.Fatalf("LoadAppConfigFromBytes returned errors: %s", joinedErrors(errs))
		}

		if config.SchemaVersion != DEFAULT_CONFIG_VERSION {
			t.Fatalf("SchemaVersion = %d, want %d", config.SchemaVersion, DEFAULT_CONFIG_VERSION)
		}
		if config.GluetunVersion != GLUETUN_TARGET_VERSION {
			t.Fatalf("GluetunVersion = %q, want %q", config.GluetunVersion, GLUETUN_TARGET_VERSION)
		}
		if config.VPNConfig.Provider != "custom" {
			t.Fatalf("VPNConfig.Provider = %q, want %q", config.VPNConfig.Provider, "custom")
		}
		if config.VPNConfig.Protocol != ProtocolWireguard {
			t.Fatalf("VPNConfig.Protocol = %q, want %q", config.VPNConfig.Protocol, ProtocolWireguard)
		}
		if config.VPNConfig.OpenVPNCredentials.Username.From != "env" {
			t.Fatalf("OpenVPNCredentials.Username.From = %q, want %q", config.VPNConfig.OpenVPNCredentials.Username.From, "env")
		}
		if config.VPNConfig.OpenVPNCredentials.Username.Name != "VPN_USERNAME" {
			t.Fatalf("OpenVPNCredentials.Username.Name = %q, want %q", config.VPNConfig.OpenVPNCredentials.Username.Name, "VPN_USERNAME")
		}
		if config.VPNConfig.OpenVPNCredentials.Password.From != "env" {
			t.Fatalf("OpenVPNCredentials.Password.From = %q, want %q", config.VPNConfig.OpenVPNCredentials.Password.From, "env")
		}
		if config.VPNConfig.OpenVPNCredentials.Password.Name != "VPN_PASSWORD" {
			t.Fatalf("OpenVPNCredentials.Password.Name = %q, want %q", config.VPNConfig.OpenVPNCredentials.Password.Name, "VPN_PASSWORD")
		}
		if config.VPNConfig.WireguardCredentials.PrivateKey.From != "env" {
			t.Fatalf("WireguardCredentials.PrivateKey.From = %q, want %q", config.VPNConfig.WireguardCredentials.PrivateKey.From, "env")
		}
		if config.VPNConfig.WireguardCredentials.PrivateKey.Name != "WIREGUARD_PRIVATE_KEY" {
			t.Fatalf("WireguardCredentials.PrivateKey.Name = %q, want %q", config.VPNConfig.WireguardCredentials.PrivateKey.Name, "WIREGUARD_PRIVATE_KEY")
		}
		if config.VPNConfig.WireguardCredentials.PublicKey.From != "env" {
			t.Fatalf("WireguardCredentials.PublicKey.From = %q, want %q", config.VPNConfig.WireguardCredentials.PublicKey.From, "env")
		}
		if config.VPNConfig.WireguardCredentials.PublicKey.Name != "WIREGUARD_PUBLIC_KEY" {
			t.Fatalf("WireguardCredentials.PublicKey.Name = %q, want %q", config.VPNConfig.WireguardCredentials.PublicKey.Name, "WIREGUARD_PUBLIC_KEY")
		}
		if config.VPNConfig.WireguardCredentials.EndpointPort != 3000 {
			t.Fatalf("WireguardCredentials.EndpointPort = %d, want %d", config.VPNConfig.WireguardCredentials.EndpointPort, 3000)
		}
		if config.Log.Level != LOG_LEVEL_INFO {
			t.Fatalf("Log.Level = %q, want %q", config.Log.Level, LOG_LEVEL_INFO)
		}
		if config.Log.Destination != "stdout" {
			t.Fatalf("Log.Destination = %q, want %q", config.Log.Destination, "stdout")
		}
		if config.Log.Filename != "vpnx.log" {
			t.Fatalf("Log.Filename = %q, want %q", config.Log.Filename, "vpnx.log")
		}
	})

	t.Run("applies defaults", func(t *testing.T) {
		config, errs := LoadAppConfigFromBytes(configJSON(t, func(config map[string]any) {
			delete(config, "schema_version")
			delete(nestedMap(t, config, "log"), "level")
			delete(nestedMap(t, config, "vpn_config"), "protocol")
		}))
		if len(errs) > 0 {
			t.Fatalf("LoadAppConfigFromBytes returned errors: %s", joinedErrors(errs))
		}

		if config.SchemaVersion != DEFAULT_CONFIG_VERSION {
			t.Fatalf("SchemaVersion = %d, want %d", config.SchemaVersion, DEFAULT_CONFIG_VERSION)
		}
		if config.VPNConfig.Protocol != DEFAULT_VPN_PROTOCOL {
			t.Fatalf("VPNConfig.Protocol = %q, want %q", config.VPNConfig.Protocol, DEFAULT_VPN_PROTOCOL)
		}
		if config.Log.Level != DEFAULT_LOG_LEVEL {
			t.Fatalf("Log.Level = %q, want %q", config.Log.Level, DEFAULT_LOG_LEVEL)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		config, errs := LoadAppConfigFromBytes([]byte("{"))
		if config != nil {
			t.Fatalf("config = %#v, want nil", config)
		}
		if len(errs) == 0 {
			t.Fatal("expected parse errors for invalid JSON")
		}
	})

	t.Run("invalid log level", func(t *testing.T) {
		config, errs := LoadAppConfigFromBytes(configJSON(t, func(config map[string]any) {
			nestedMap(t, config, "log")["level"] = "verbose"
		}))
		if config != nil {
			t.Fatalf("config = %#v, want nil", config)
		}
		if len(errs) == 0 {
			t.Fatal("expected validation errors for invalid log level")
		}
	})

	t.Run("invalid protocol", func(t *testing.T) {
		config, errs := LoadAppConfigFromBytes(configJSON(t, func(config map[string]any) {
			nestedMap(t, config, "vpn_config")["protocol"] = "pptp"
		}))
		if config != nil {
			t.Fatalf("config = %#v, want nil", config)
		}
		if len(errs) == 0 {
			t.Fatal("expected validation errors for invalid protocol")
		}
	})

	t.Run("unsupported schema version", func(t *testing.T) {
		config, errs := LoadAppConfigFromBytes(configJSON(t, func(config map[string]any) {
			config["schema_version"] = 999
		}))
		if config != nil {
			t.Fatalf("config = %#v, want nil", config)
		}
		requireErrorContains(t, errs, "unsupported config version")
	})
}

func TestLoadAppConfigFromFile(t *testing.T) {
	t.Run("valid JSON file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "config.json")
		if err := os.WriteFile(path, configJSON(t, nil), 0o600); err != nil {
			t.Fatalf("write config file: %v", err)
		}

		config, errs := LoadAppConfigFromFile(path)
		if len(errs) > 0 {
			t.Fatalf("LoadAppConfigFromFile returned errors: %s", joinedErrors(errs))
		}
		if config.GluetunVersion != GLUETUN_TARGET_VERSION {
			t.Fatalf("GluetunVersion = %q, want %q", config.GluetunVersion, GLUETUN_TARGET_VERSION)
		}
	})

	t.Run("unsupported file extension", func(t *testing.T) {
		config, errs := LoadAppConfigFromFile(filepath.Join(t.TempDir(), "config.yaml"))
		if config != nil {
			t.Fatalf("config = %#v, want nil", config)
		}
		requireErrorContains(t, errs, "unsupported config file extension")
	})

	t.Run("missing file", func(t *testing.T) {
		config, errs := LoadAppConfigFromFile(filepath.Join(t.TempDir(), "missing.json"))
		if config != nil {
			t.Fatalf("config = %#v, want nil", config)
		}
		requireErrorContains(t, errs, "failed to read config file")
	})
}

func TestGetAppConfigVersion(t *testing.T) {
	t.Run("explicit version", func(t *testing.T) {
		version, errs := GetAppConfigVersion(configJSON(t, nil))
		if len(errs) > 0 {
			t.Fatalf("GetAppConfigVersion returned errors: %v", errs)
		}
		if version != DEFAULT_CONFIG_VERSION {
			t.Fatalf("version = %d, want %d", version, DEFAULT_CONFIG_VERSION)
		}
	})

	t.Run("defaults when omitted", func(t *testing.T) {
		version, errs := GetAppConfigVersion(configJSON(t, func(config map[string]any) {
			delete(config, "schema_version")
		}))
		if len(errs) > 0 {
			t.Fatalf("GetAppConfigVersion returned errors: %v", errs)
		}
		if version != DEFAULT_CONFIG_VERSION {
			t.Fatalf("version = %d, want %d", version, DEFAULT_CONFIG_VERSION)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		version, errs := GetAppConfigVersion([]byte("{"))
		if len(errs) == 0 {
			t.Fatal("expected validation errors for invalid JSON")
		}
		if version != 0 {
			t.Fatalf("version = %d, want %d", version, 0)
		}
	})
}

func configWithEnvVars(version string, envVars ...string) ConfigVersionOne {
	return ConfigVersionOne{
		GluetunVersion: version,
		EnvVars:        envVars,
	}
}

func TestExtractEnvVars(t *testing.T) {
	t.Run("valid env vars", func(t *testing.T) {
		envVars, errs := ExtractEnvVars(configWithEnvVars(
			GLUETUN_TARGET_VERSION,
			"TZ=America/Toronto",
			"PUBLIC_IP_ENABLED=on",
			"FIREWALL_INPUT_PORTS=80,443",
			"HTTPPROXY_PASSWORD=secret=with=equals",
		))
		if len(errs) > 0 {
			t.Fatalf("ExtractEnvVars returned errors: %s", joinedErrors(errs))
		}

		if got := envVars["TZ"]; got != "America/Toronto" {
			t.Fatalf("envVars[TZ] = %q, want %q", got, "America/Toronto")
		}
		if got := envVars["PUBLIC_IP_ENABLED"]; got != "on" {
			t.Fatalf("envVars[PUBLIC_IP_ENABLED] = %q, want %q", got, "on")
		}
		if got := envVars["FIREWALL_INPUT_PORTS"]; got != "80,443" {
			t.Fatalf("envVars[FIREWALL_INPUT_PORTS] = %q, want %q", got, "80,443")
		}
		if got := envVars["HTTPPROXY_PASSWORD"]; got != "secret=with=equals" {
			t.Fatalf("envVars[HTTPPROXY_PASSWORD] = %q, want %q", got, "secret=with=equals")
		}
	})

	t.Run("empty env vars", func(t *testing.T) {
		envVars, errs := ExtractEnvVars(configWithEnvVars(GLUETUN_TARGET_VERSION))
		if len(errs) > 0 {
			t.Fatalf("ExtractEnvVars returned errors: %s", joinedErrors(errs))
		}
		if len(envVars) != 0 {
			t.Fatalf("len(envVars) = %d, want 0", len(envVars))
		}
	})

	t.Run("duplicate keys use last value", func(t *testing.T) {
		envVars, errs := ExtractEnvVars(configWithEnvVars(
			GLUETUN_TARGET_VERSION,
			"TZ=UTC",
			"TZ=America/Toronto",
		))
		if len(errs) > 0 {
			t.Fatalf("ExtractEnvVars returned errors: %s", joinedErrors(errs))
		}
		if got := envVars["TZ"]; got != "America/Toronto" {
			t.Fatalf("envVars[TZ] = %q, want %q", got, "America/Toronto")
		}
	})

	t.Run("invalid format", func(t *testing.T) {
		envVars, errs := ExtractEnvVars(configWithEnvVars(GLUETUN_TARGET_VERSION, "MALFORMED"))
		if envVars != nil {
			t.Fatalf("envVars = %#v, want nil", envVars)
		}
		requireErrorContains(t, errs, "invalid env var format")
	})

	t.Run("unknown env var", func(t *testing.T) {
		envVars, errs := ExtractEnvVars(configWithEnvVars(GLUETUN_TARGET_VERSION, "NOT_A_REAL_ENV=value"))
		if envVars != nil {
			t.Fatalf("envVars = %#v, want nil", envVars)
		}
		requireErrorContains(t, errs, "unknown env var")
	})

	t.Run("invalid env var value", func(t *testing.T) {
		envVars, errs := ExtractEnvVars(configWithEnvVars(GLUETUN_TARGET_VERSION, "PUBLIC_IP_ENABLED=maybe"))
		if envVars != nil {
			t.Fatalf("envVars = %#v, want nil", envVars)
		}
		requireErrorContains(t, errs, "invalid env var value: PUBLIC_IP_ENABLED")
	})

	t.Run("unsupported gluetun version", func(t *testing.T) {
		envVars, errs := ExtractEnvVars(configWithEnvVars("0.0.0", "TZ=UTC"))
		if envVars != nil {
			t.Fatalf("envVars = %#v, want nil", envVars)
		}
		requireErrorContains(t, errs, "unsupported gluetun version")
	})
}
