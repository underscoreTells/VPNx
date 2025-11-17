package config

import "time"

// CLIConfig configures the host-side vpnx CLI.
type CLIConfig struct {
	// SocketPath is the Unix domain socket the CLI uses to talk to vpnxd.
	SocketPath string `yaml:"socket_path"`

	// RequestTimeout is the default timeout for CLI requests to the daemon.
	RequestTimeout time.Duration `yaml:"request_timeout,omitempty"`

	// LogLevel controls verbosity of CLI output (e.g. "debug", "info", "warn", "error").
	LogLevel string `yaml:"log_level,omitempty"`

	// Container holds optional container-integration settings for the CLI.
	Container CLIContainerConfig `yaml:"container,omitempty"`
}

// CLIContainerConfig holds optional container-integration settings for the CLI.
type CLIContainerConfig struct {
	// Name is the container name or ID running the vpnxd daemon.
	Name string `yaml:"name"`

	// Runtime is the container runtime to use (e.g. "docker", "podman").
	Runtime string `yaml:"runtime,omitempty"`

	// Manage indicates whether the CLI is allowed to manage the container lifecycle.
	Manage bool `yaml:"manage,omitempty"`
}

// DefaultCLIConfig returns a CLIConfig populated with reasonable defaults.
func DefaultCLIConfig() *CLIConfig {
	return &CLIConfig{
		SocketPath:     "/var/run/vpnx.sock",
		RequestTimeout: 5 * time.Second,
		LogLevel:       "info",
	}
}
