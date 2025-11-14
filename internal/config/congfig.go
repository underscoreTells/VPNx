package config

type Type string //Network type

const (
	TypeTailscale Type = "tailscale"
	TypeWireguard Type = "wireguard"
)

type ControlPlane struct {
	Enabled     bool     `yaml:"enabled"`
	Endpoints   []string `yaml:"endpoints"`
	IPAddresses []string `yaml:"ip_addresses"`
}

type Tailscale struct {
	MagicDNS bool `yaml:"magic_dns"`
}

type Coordinator struct {
	Enabld    bool     `yaml:"enabled"`
	URL       string   `yaml:"url"`
	Endpoints []string `yaml:"endpoints"`
}

type Wireguard struct {
	Subnet      string   `yaml:"subnet"`
	Peers       []string `yaml:"peers"`
	Coordinator Coordinator
}

type Network struct {
	Type            Type     `yaml:"type"`
	Interface       string   `yaml:"interface"`
	ExcludedSubnets []string `yaml:"excluded_subnets"`
	ControlPlane    ControlPlane
	Tailscale       Tailscale
	Wireguard       Wireguard
}

type VPNType string

const (
	VPNTypeWireguard VPNType = "wireguard"
)

type VPN struct {
	Type       VPNType `yaml:"type"`
	Interface  string  `yaml:"interface"`
	ConfigPath string  `yaml:"config_path"`
}

type Strategy string

const (
	PolicyBased Strategy = "policy_based"
	TableBased  Strategy = "talbe_based"
)

type Policy struct {
	Mark      int    `yaml:"mark"`
	TableID   int    `yaml:"table_id"`
	TableName string `yaml:"table_name"`
}

type Priotity struct {
	NetworkControl  int `yaml:"network_control"`
	NetworkInternal int `yaml:"network_internal"`
	VPNForwarded    int `yaml:"vpn_forwarded"`
	Default         int `yaml:"default"`
}

type Routing struct {
	Strategy Strategy `yaml:"strategy"`
	Policy   Policy   `yaml:"policy"`
	Priotity Priotity `yaml:"priotity"`
}

type Backend string

const (
	Auto           Backend = "auto"
	IpTables       Backend = "iptables"
	NfTables       Backend = "nftables"
	Firewalld      Backend = "firewalld"
	Ufw            Backend = "ufw"
	SusSEfirewall2 Backend = "susSEfirewall2"
	Shorewall      Backend = "shorewall"
)

type Config struct {
	Version int `yaml:"version"`
	Network Network
}
