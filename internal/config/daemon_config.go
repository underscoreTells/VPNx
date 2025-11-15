package config

// DaemonConfig defines the configuration for the vpnx daemon.
// The config file is typically authored on the host and may be
// mounted into a container or read directly when running on the host.

type DaemonConfig struct {

	// Version of the daemon configuration schema.
	Version int `yaml:"version"`

	Daemon   DaemonSection   `yaml:"daemon"`
	Mesh     MeshSection     `yaml:"mesh"`
	Provider ProviderSection `yaml:"provider"`
	Exit     ExitSection     `yaml:"exit"`
	Logging  LoggingSection  `yaml:"logging"`
}

// DaemonSection holds daemon-local settings inside the container.
type DaemonSection struct {
	// APISocket is the Unix domain socket path used by the daemon.
	APISocket string `yaml:"api_socket"`

	// StateDir is the directory for daemon state inside the container.
	StateDir string `yaml:"state_dir"`

	// HealthAddress is an optional address for an HTTP health endpoint.
	HealthAddress string `yaml:"health_address"`
}

// MeshProviderType identifies the mesh provider.
type MeshProviderType string

const (
	MeshProviderTailscale MeshProviderType = "tailscale"
)

// MeshSection configures mesh integration (e.g. Tailscale).
type MeshSection struct {
	Provider  MeshProviderType     `yaml:"provider"`
	Tailscale TailscaleMeshSection `yaml:"tailscale"`
}

// TailscaleMeshSection contains Tailscale-specific settings.
type TailscaleMeshSection struct {
	Hostname    string `yaml:"hostname"`
	AuthKey     string `yaml:"auth_key"`
	ControlURL  string `yaml:"control_url"`
	UseMagicDNS bool   `yaml:"magic_dns"`
}

// ProviderType identifies the upstream VPN provider type.
type ProviderType string

const (
	ProviderWireGuard ProviderType = "wireguard"
	ProviderOpenVPN   ProviderType = "openvpn"
)

// ProviderSection configures the upstream VPN provider. The daemon
// only ever loads external configuration files (e.g. wg-quick configs
// for WireGuard, and OpenVPN configs in the future). ConfigPath values
// are interpreted as paths on the host and must be made available to
// the daemon by whatever launches it (for example via bind mounts when
// running inside a container).
type ProviderSection struct {
	Type ProviderType `yaml:"type"`

	WireGuard WireGuardConfigSection `yaml:"wireguard"`
	OpenVPN   OpenVPNConfigSection   `yaml:"openvpn"`
}

// WireGuardConfigSection references an external WireGuard configuration
// file (typically wg-quick compatible) on the host filesystem.
type WireGuardConfigSection struct {
	ConfigPath string `yaml:"config_path"`
}

// OpenVPNConfigSection references an external OpenVPN configuration
// file on the host filesystem.
type OpenVPNConfigSection struct {
	ConfigPath string `yaml:"config_path"`
}

// ExitDNSMode controls DNS behavior for exit traffic.
type ExitDNSMode string

const (
	ExitDNSModePassthrough ExitDNSMode = "passthrough"
	ExitDNSModeOverride    ExitDNSMode = "override"
	ExitDNSModeSplit       ExitDNSMode = "split"
)

// ExitSection defines how this node behaves as an exit node.
type ExitSection struct {
	AdvertiseSubnets []string          `yaml:"advertise_subnets"`
	ExcludeSubnets   []string          `yaml:"exclude_subnets"`
	DNS              ExitDNSSection    `yaml:"dns"`
	Killswitch       KillswitchSection `yaml:"killswitch"`
}

// ExitDNSSection configures DNS behavior for exit traffic.
type ExitDNSSection struct {
	Mode      ExitDNSMode `yaml:"mode"`
	Upstreams []string    `yaml:"upstreams"`
}

// KillswitchSection configures killswitch behavior.
type KillswitchSection struct {
	Enabled  bool `yaml:"enabled"`
	AllowLAN bool `yaml:"allow_lan"`
}

// LoggingSection configures daemon logging.
type LoggingSection struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func DefaultDaemonConfig() DaemonConfig {
	return DaemonConfig{
		Version: 1,
		Daemon: DaemonSection{
			APISocket:     "/run/vpnx/api.sock",
			StateDir:      "/var/lib/vpnx",
			HealthAddress: ":8080",
		},
		Mesh: MeshSection{
			Provider: MeshProviderTailscale,
			Tailscale: TailscaleMeshSection{
				Hostname:    "vpn.example.com",
				ControlURL:  "https://tailscale.example.com",
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
			AdvertiseSubnets: []string{},
			ExcludeSubnets:   []string{},
			DNS: ExitDNSSection{
				Mode:      ExitDNSModePassthrough,
				Upstreams: []string{},
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
