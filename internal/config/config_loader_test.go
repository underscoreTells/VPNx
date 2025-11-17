package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// mockDefaultCLIConfig returns a default-like CLIConfig tailored for tests.
func mockDefaultCLIConfig() *CLIConfig {
	return &CLIConfig{
		SocketPath:     "/default.sock",
		RequestTimeout: 10 * time.Second,
		LogLevel:       "info",
		Container: CLIContainerConfig{
			Name:    "default-container",
			Runtime: "docker",
			Manage:  false,
		},
	}
}

// mockUserCLIConfig returns a CLIConfig that overrides some, but not all, fields.
func mockUserCLIConfig() *CLIConfig {
	return &CLIConfig{
		SocketPath: "/custom.sock",
		// RequestTimeout left as zero, should keep default
		LogLevel: "debug",
		Container: CLIContainerConfig{
			Name:   "custom-container",
			Manage: true,
			// Runtime left as empty string -> zero, should keep default
		},
	}
}

// mockDefaultDaemonConfig returns a default-like DaemonConfig tailored for tests.
func mockDefaultDaemonConfig() *DaemonConfig {
	return &DaemonConfig{
		Version: 1,
		Daemon: DaemonSection{
			APISocket:     "/run/vpnx/api.sock",
			StateDir:      "/var/lib/vpnx",
			HealthAddress: ":8080",
		},
		Mesh: MeshSection{
			Provider: MeshProviderTailscale,
			Tailscale: TailscaleMeshSection{
				Hostname:    "default-host",
				AuthKey:     "default-key",
				ControlURL:  "https://default.example.com",
				UseMagicDNS: true,
			},
		},
		Provider: ProviderSection{
			Type: ProviderWireGuard,
			WireGuard: WireGuardConfigSection{
				ConfigPath: "/etc/vpnx/provider.conf",
			},
			OpenVPN: OpenVPNConfigSection{
				ConfigPath: "",
			},
		},
		Exit: ExitSection{
			AdvertiseSubnets: []string{"0.0.0.0/0"},
			ExcludeSubnets:   []string{"10.0.0.0/8"},
			DNS: ExitDNSSection{
				Mode:      ExitDNSModePassthrough,
				Upstreams: []string{"1.1.1.1"},
			},
			Killswitch: KillswitchSection{
				Enabled:  true,
				AllowLAN: true,
			},
		},
		Logging: LoggingSection{
			Level:  "info",
			Format: "text",
		},
	}
}

// mockUserDaemonConfig returns a DaemonConfig that overrides some fields.
func mockUserDaemonConfig() *DaemonConfig {
	return &DaemonConfig{
		// Version left as zero -> default should be kept
		Daemon: DaemonSection{
			APISocket: "/custom/api.sock",
			// StateDir zero, should keep default
			HealthAddress: ":9090",
		},
		Mesh: MeshSection{
			// Provider left as zero -> keep default
			Tailscale: TailscaleMeshSection{
				Hostname: "custom-host",
				// AuthKey zero, should keep default
				ControlURL:  "",
				UseMagicDNS: false,
			},
		},
		Provider: ProviderSection{
			// Type zero, should keep default
			WireGuard: WireGuardConfigSection{
				ConfigPath: "/custom/provider.conf", // non-zero
			},
		},
		Exit: ExitSection{
			AdvertiseSubnets: []string{"192.0.2.0/24"},
			// ExcludeSubnets left nil -> keep default
			DNS: ExitDNSSection{
				Mode: ExitDNSModeOverride,
				// Upstreams left nil -> keep default
			},
			Killswitch: KillswitchSection{
				Enabled:  true,
				AllowLAN: false,
			},
		},
		Logging: LoggingSection{
			Level:  "debug",
			Format: "",
		},
	}
}

// writeTempYAML writes YAML content to a temp file and returns its path.
func writeTempYAML(t *testing.T, name, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, name)

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp YAML file: %v", err)
	}

	return path
}

func TestYAMLLoader_overwriteDefault_MergesNonZeroFields_CLI(t *testing.T) {
	defaultCfg := mockDefaultCLIConfig()
	userCfg := mockUserCLIConfig()

	loader := NewYAMLLoader(defaultCfg)

	got, err := loader.overwriteDefault(userCfg)
	if err != nil {
		t.Fatalf("overwriteDefault returned error: %v", err)
	}

	// Fields that should be overridden
	if got.SocketPath != "/custom.sock" {
		t.Fatalf("SocketPath = %q, want %q", got.SocketPath, "/custom.sock")
	}
	if got.LogLevel != "debug" {
		t.Fatalf("LogLevel = %q, want %q", got.LogLevel, "debug")
	}

	// Fields that should stay at default because user value is zero
	if got.RequestTimeout != defaultCfg.RequestTimeout {
		t.Fatalf("RequestTimeout = %v, want %v", got.RequestTimeout, defaultCfg.RequestTimeout)
	}

	// Nested struct merge behavior
	if got.Container.Name != "custom-container" {
		t.Fatalf("Container.Name = %q, want %q", got.Container.Name, "custom-container")
	}
	// Runtime was zero in user config, should keep default
	if got.Container.Runtime != defaultCfg.Container.Runtime {
		t.Fatalf("Container.Runtime = %q, want %q", got.Container.Runtime, defaultCfg.Container.Runtime)
	}
	// Manage was non-zero in user config, should override
	if got.Container.Manage != true {
		t.Fatalf("Container.Manage = %v, want %v", got.Container.Manage, true)
	}
}

func TestYAMLLoader_overwriteDefault_IgnoresEmptyContainerStruct(t *testing.T) {
	defaultCfg := mockDefaultCLIConfig()
	userCfg := &CLIConfig{
		// override some other field so we know merge is happening
		LogLevel:  "debug",
		Container: CLIContainerConfig{}, // all-zero container: should NOT override defaults
	}

	loader := NewYAMLLoader(defaultCfg)

	got, err := loader.overwriteDefault(userCfg)
	if err != nil {
		t.Fatalf("overwriteDefault returned error: %v", err)
	}

	// LogLevel was non-zero in user config, should override
	if got.LogLevel != "debug" {
		t.Fatalf("LogLevel = %q, want %q", got.LogLevel, "debug")
	}

	// Container was zero in user config, should keep default container fields
	if got.Container.Name != defaultCfg.Container.Name {
		t.Fatalf("Container.Name = %q, want default %q", got.Container.Name, defaultCfg.Container.Name)
	}
	if got.Container.Runtime != defaultCfg.Container.Runtime {
		t.Fatalf("Container.Runtime = %q, want default %q", got.Container.Runtime, defaultCfg.Container.Runtime)
	}
	if got.Container.Manage != defaultCfg.Container.Manage {
		t.Fatalf("Container.Manage = %v, want default %v", got.Container.Manage, defaultCfg.Container.Manage)
	}
}

func TestYAMLLoader_overwriteDefault_MergesNonZeroFields_Daemon(t *testing.T) {
	defaultCfg := mockDefaultDaemonConfig()
	userCfg := mockUserDaemonConfig()

	loader := NewYAMLLoader(defaultCfg)

	got, err := loader.overwriteDefault(userCfg)
	if err != nil {
		t.Fatalf("overwriteDefault returned error: %v", err)
	}

	// Top-level scalar: Version should remain default (user zero)
	if got.Version != defaultCfg.Version {
		t.Fatalf("Version = %d, want %d", got.Version, defaultCfg.Version)
	}

	// Nested Daemon section
	if got.Daemon.APISocket != "/custom/api.sock" {
		t.Fatalf("Daemon.APISocket = %q, want %q", got.Daemon.APISocket, "/custom/api.sock")
	}
	if got.Daemon.StateDir != defaultCfg.Daemon.StateDir {
		t.Fatalf("Daemon.StateDir = %q, want default %q", got.Daemon.StateDir, defaultCfg.Daemon.StateDir)
	}
	if got.Daemon.HealthAddress != ":9090" {
		t.Fatalf("Daemon.HealthAddress = %q, want %q", got.Daemon.HealthAddress, ":9090")
	}

	// Mesh section
	if got.Mesh.Provider != defaultCfg.Mesh.Provider {
		t.Fatalf("Mesh.Provider = %q, want default %q", got.Mesh.Provider, defaultCfg.Mesh.Provider)
	}
	if got.Mesh.Tailscale.Hostname != "custom-host" {
		t.Fatalf("Mesh.Tailscale.Hostname = %q, want %q", got.Mesh.Tailscale.Hostname, "custom-host")
	}
	if got.Mesh.Tailscale.AuthKey != defaultCfg.Mesh.Tailscale.AuthKey {
		t.Fatalf("Mesh.Tailscale.AuthKey = %q, want default %q", got.Mesh.Tailscale.AuthKey, defaultCfg.Mesh.Tailscale.AuthKey)
	}
	// ControlURL zero in user, should keep default
	if got.Mesh.Tailscale.ControlURL != defaultCfg.Mesh.Tailscale.ControlURL {
		t.Fatalf("Mesh.Tailscale.ControlURL = %q, want default %q", got.Mesh.Tailscale.ControlURL, defaultCfg.Mesh.Tailscale.ControlURL)
	}
	// UseMagicDNS non-zero in user, should override
	if got.Mesh.Tailscale.UseMagicDNS != false {
		t.Fatalf("Mesh.Tailscale.UseMagicDNS = %v, want %v", got.Mesh.Tailscale.UseMagicDNS, false)
	}

	// Provider section
	if got.Provider.Type != defaultCfg.Provider.Type {
		t.Fatalf("Provider.Type = %q, want default %q", got.Provider.Type, defaultCfg.Provider.Type)
	}
	if got.Provider.WireGuard.ConfigPath != "/custom/provider.conf" {
		t.Fatalf("Provider.WireGuard.ConfigPath = %q, want %q", got.Provider.WireGuard.ConfigPath, "/custom/provider.conf")
	}

	// Exit section: slices and nested structs
	if len(got.Exit.AdvertiseSubnets) != 1 || got.Exit.AdvertiseSubnets[0] != "192.0.2.0/24" {
		t.Fatalf("Exit.AdvertiseSubnets = %#v, want [\"192.0.2.0/24\"]", got.Exit.AdvertiseSubnets)
	}
	if len(got.Exit.ExcludeSubnets) != len(defaultCfg.Exit.ExcludeSubnets) {
		t.Fatalf("Exit.ExcludeSubnets length = %d, want %d", len(got.Exit.ExcludeSubnets), len(defaultCfg.Exit.ExcludeSubnets))
	}
	if got.Exit.DNS.Mode != ExitDNSModeOverride {
		t.Fatalf("Exit.DNS.Mode = %q, want %q", got.Exit.DNS.Mode, ExitDNSModeOverride)
	}
	// Upstreams should stay default (user nil)
	if len(got.Exit.DNS.Upstreams) != len(defaultCfg.Exit.DNS.Upstreams) {
		t.Fatalf("Exit.DNS.Upstreams length = %d, want %d", len(got.Exit.DNS.Upstreams), len(defaultCfg.Exit.DNS.Upstreams))
	}
	if got.Exit.Killswitch.Enabled != true {
		t.Fatalf("Exit.Killswitch.Enabled = %v, want %v", got.Exit.Killswitch.Enabled, true)
	}
	if got.Exit.Killswitch.AllowLAN != false {
		t.Fatalf("Exit.Killswitch.AllowLAN = %v, want %v", got.Exit.Killswitch.AllowLAN, false)
	}

	// Logging section
	if got.Logging.Level != "debug" {
		t.Fatalf("Logging.Level = %q, want %q", got.Logging.Level, "debug")
	}
	// Format zero in user, should keep default
	if got.Logging.Format != defaultCfg.Logging.Format {
		t.Fatalf("Logging.Format = %q, want default %q", got.Logging.Format, defaultCfg.Logging.Format)
	}
}

func TestYAMLLoader_overwriteDefault_ErrorsOnNilDefault(t *testing.T) {
	var nilDefault *CLIConfig
	loader := NewYAMLLoader(nilDefault)

	cfg := mockUserCLIConfig()
	_, err := loader.overwriteDefault(cfg)
	if err == nil {
		t.Fatalf("expected error when DefaultConfig is nil, got nil")
	}
}

func TestYAMLLoader_overwriteDefault_ErrorsOnNilConfig(t *testing.T) {
	defaultCfg := mockDefaultCLIConfig()
	loader := NewYAMLLoader(defaultCfg)

	_, err := loader.overwriteDefault(nil)
	if err == nil {
		t.Fatalf("expected error when config is nil, got nil")
	}
}

func TestYAMLLoader_Load_MergesWithDefaults_CLI(t *testing.T) {
	defaultCfg := mockDefaultCLIConfig()
	loader := NewYAMLLoader(defaultCfg)

	yamlContent := `
socket_path: "/from-file.sock"
log_level: "warn"
container:
  name: "file-container"
  manage: true
`
	path := writeTempYAML(t, "cli.yaml", yamlContent)

	got, err := loader.Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	// Values from file
	if got.SocketPath != "/from-file.sock" {
		t.Fatalf("SocketPath = %q, want %q", got.SocketPath, "/from-file.sock")
	}
	if got.LogLevel != "warn" {
		t.Fatalf("LogLevel = %q, want %q", got.LogLevel, "warn")
	}
	if got.Container.Name != "file-container" {
		t.Fatalf("Container.Name = %q, want %q", got.Container.Name, "file-container")
	}
	if got.Container.Manage != true {
		t.Fatalf("Container.Manage = %v, want %v", got.Container.Manage, true)
	}

	// Values that should remain at defaults (not set in YAML)
	if got.RequestTimeout != defaultCfg.RequestTimeout {
		t.Fatalf("RequestTimeout = %v, want default %v", got.RequestTimeout, defaultCfg.RequestTimeout)
	}
	if got.Container.Runtime != defaultCfg.Container.Runtime {
		t.Fatalf("Container.Runtime = %q, want default %q", got.Container.Runtime, defaultCfg.Container.Runtime)
	}
}

func TestYAMLLoader_Load_IgnoresEmptyContainerStruct(t *testing.T) {
	defaultCfg := mockDefaultCLIConfig()
	loader := NewYAMLLoader(defaultCfg)

	yamlContent := `
socket_path: "/from-file.sock"
container: {}
`
	path := writeTempYAML(t, "cli_empty_container.yaml", yamlContent)

	got, err := loader.Load(path)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	// SocketPath overridden from YAML
	if got.SocketPath != "/from-file.sock" {
		t.Fatalf("SocketPath = %q, want %q", got.SocketPath, "/from-file.sock")
	}

	// Container was empty in YAML; keep defaults
	if got.Container.Name != defaultCfg.Container.Name {
		t.Fatalf("Container.Name = %q, want default %q", got.Container.Name, defaultCfg.Container.Name)
	}
	if got.Container.Runtime != defaultCfg.Container.Runtime {
		t.Fatalf("Container.Runtime = %q, want default %q", got.Container.Runtime, defaultCfg.Container.Runtime)
	}
	if got.Container.Manage != defaultCfg.Container.Manage {
		t.Fatalf("Container.Manage = %v, want default %v", got.Container.Manage, defaultCfg.Container.Manage)
	}
}

func TestYAMLLoader_Load_FileNotFound(t *testing.T) {
	loader := NewCLILoader()

	_, err := loader.Load("does-not-exist.yaml")
	if err == nil {
		t.Fatalf("expected error for missing file, got nil")
	}
}

func TestYAMLLoader_Load_InvalidYAML(t *testing.T) {
	defaultCfg := mockDefaultCLIConfig()
	loader := NewYAMLLoader(defaultCfg)

	invalid := "::: this is not valid yaml :::"
	path := writeTempYAML(t, "invalid.yaml", invalid)

	_, err := loader.Load(path)
	if err == nil {
		t.Fatalf("expected error for invalid YAML, got nil")
	}
}

func TestNewCLILoader_UsesDefaultCLIConfig(t *testing.T) {
	loader := NewCLILoader()
	if loader == nil || loader.DefaultConfig == nil {
		t.Fatalf("NewCLILoader returned nil or had nil DefaultConfig")
	}

	got := loader.DefaultConfig
	def := DefaultCLIConfig()

	if got.SocketPath != def.SocketPath ||
		got.LogLevel != def.LogLevel ||
		got.RequestTimeout != def.RequestTimeout {
		t.Fatalf("NewCLILoader default config = %#v, want %#v", got, def)
	}
}

func TestNewDaemonLoader_UsesDefaultDaemonConfig(t *testing.T) {
	loader := NewDaemonLoader()
	if loader == nil || loader.DefaultConfig == nil {
		t.Fatalf("NewDaemonLoader returned nil or had nil DefaultConfig")
	}

	got := loader.DefaultConfig
	def := DefaultDaemonConfig()

	if got.Version != def.Version ||
		got.Daemon.APISocket != def.Daemon.APISocket ||
		got.Logging.Level != def.Logging.Level {
		t.Fatalf("NewDaemonLoader default config = %#v, want %#v", got, def)
	}
}
