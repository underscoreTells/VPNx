package gluetun

import "github.com/underscoreTells/vpn-exit-node/internal/config"

type VPNProvider string

const (
	AirVPN                VPNProvider = "airvpn"
	Custom                VPNProvider = "Custom"
	CyberGhost            VPNProvider = "cyberghost"
	ExpressVPN            VPNProvider = "expressvpn"
	FastestVPN            VPNProvider = "fastestvpn"
	Giganews              VPNProvider = "giganews"
	Hidemyass             VPNProvider = "hidemyass"
	IPvanish              VPNProvider = "ipvanish"
	Ivpn                  VPNProvider = "ivpn"
	Mullvad               VPNProvider = "mullvad"
	NordVPN               VPNProvider = "nordvpn"
	PerfectPrivacy        VPNProvider = "perfect-privacy"
	Privado               VPNProvider = "privado"
	PrivateInternetAccess VPNProvider = "private-internet-access"
	PrivateVPN            VPNProvider = "privatevpn"
	ProtonVPN             VPNProvider = "protonvpn"
	PureVPN               VPNProvider = "purevpn"
	SlickVPN              VPNProvider = "slickvpn"
	Surfshark             VPNProvider = "surfshark"
	Torguard              VPNProvider = "torguard"
	VPNSecure             VPNProvider = "vpn-secure"
	VPNUnlimited          VPNProvider = "vpn-unlimited"
	VyprVPN               VPNProvider = "vyprvpn"
	Windscribe            VPNProvider = "windscribe"
)

type HttpControlServerLog string

const (
	HttpControlServerLogOn  HttpControlServerLog = "on"
	HttpControlServerLogOff HttpControlServerLog = "off"
)

type ControlServerOptions struct {
	Address            string               `mapstructure:"HTTP_CONTROL_SERVER_ADDRESS,omitempty"`
	Log                HttpControlServerLog `mapstructure:"HTTP_CONTROL_SERVER_LOG,omitempty"`
	AuthConfigFilepath string               `mapstructure:"HTTP_CONTROL_SERVER_AUTH_CONFIG_FILEPATH,omitempty"`
}

type DnsUpstreamResolverType string

const (
	DnsUpstreamResolverTypeDot   DnsUpstreamResolverType = "dot"
	DnsUpstreamResolverTypeDoh   DnsUpstreamResolverType = "doh"
	DnsUpstreamResolverTypePlain DnsUpstreamResolverType = "plain"
)

type DnsUpstreamResolver string

const (
	CircaFamily           DnsUpstreamResolver = "circa family"
	CircaPrivate          DnsUpstreamResolver = "circa private"
	CircaProtected        DnsUpstreamResolver = "circa protected"
	CleanBrowsing         DnsUpstreamResolver = "cleanbrowsing"
	Adult                 DnsUpstreamResolver = "adult"
	CleanBrowsingFamily   DnsUpstreamResolver = "cleanbrowsing family"
	CleanBrowsingSecurity DnsUpstreamResolver = "cleanbrowsing security"
	Cloudflare            DnsUpstreamResolver = "cloudflare"
	CloudflareFamily      DnsUpstreamResolver = "cloudflare family"
	CloudflareSecurity    DnsUpstreamResolver = "cloudflare security"
	Google                DnsUpstreamResolver = "google"
	Libredns              DnsUpstreamResolver = "libredns"
	Opendns               DnsUpstreamResolver = "opendns"
	Quad9                 DnsUpstreamResolver = "quad9"
	Quad9Secured          DnsUpstreamResolver = "quad9 secured"
	Quad9Unsecured        DnsUpstreamResolver = "quad9 unsecured"
	Quadrant              DnsUpstreamResolver = "quadrant"
)

type DnsCaching string

const (
	DnsCachingOn  DnsCaching = "on"
	DnsCachingOff DnsCaching = "off"
)

type DnsUpstreamIPV6 string

const (
	DnsUpstreamIPV6On  DnsUpstreamIPV6 = "on"
	DnsUpstreamIPV6Off DnsUpstreamIPV6 = "off"
)

type BlockMalicious string

const (
	BlockMaliciousOn  BlockMalicious = "on"
	BlockMaliciousOff BlockMalicious = "off"
)

type BlockSurveillance string

const (
	BlockSurveillanceOn  BlockSurveillance = "on"
	BlockSurveillanceOff BlockSurveillance = "off"
)

type BlockAds string

const (
	BlockAdsOn  BlockAds = "on"
	BlockAdsOff BlockAds = "off"
)

type DNS struct {
	UpstreamResolverType               DnsUpstreamResolverType `mapstructure:"DNS_UPSTREAM_RESOLVER_TYPE,omitempty"`
	UpstreamResolvers                  []DnsUpstreamResolver   `mapstructure:"DNS_UPSTREAM_RESOLVER,omitempty"`
	Caching                            DnsCaching              `mapstructure:"DNS_CACHING,omitempty"`
	UpstreamIPV6                       DnsUpstreamIPV6         `mapstructure:"DNS_UPSTREAM_IPV6,omitempty"`
	BlockIps                           []string                `mapstructure:"DNS_BLOCK_IPS,omitempty"`
	BlockIpPrefixes                    []string                `mapstructure:"DNS_BLOCK_IP_PREFIXES,omitempty"`
	RebindingProtectionExemptHostnames []string                `mapstructure:"DNS_REBINDING_PROTECTION_EXEMPT_HOSTNAMES,omitempty"`
	UpdatePeriod                       string                  `mapstructure:"DNS_UPDATE_PERIOD,omitempty"`
	BlockMalicious                     BlockMalicious          `mapstructure:"BLOCK_MALICIOUS,omitempty"`
	BlockSurveillance                  BlockSurveillance       `mapstructure:"BLOCK_SURVEILLANCE,omitempty"`
	BlockAds                           BlockAds                `mapstructure:"BLOCK_ADS,omitempty"`
	UnblockHostnames                   []string                `mapstructure:"DNS_UNBLOCK_HOSTNAMES,omitempty"`
	UpstreamPlainAddresses             []string                `mapstructure:"DNS_UPSTREAM_PLAIN_ADDRESSES,omitempty"`
}

type FirewallDebug string

const (
	FirewallDebugOn  FirewallDebug = "on"
	FirewallDebugOff FirewallDebug = "off"
)

type Firewall struct {
	VpnInputPorts   []int         `mapstructure:"FIREWALL_VPN_INPUT_PORTS,omitempty"`
	InputPorts      []int         `mapstructure:"FIREWALL_INPUT_PORTS,omitempty"`
	Debug           FirewallDebug `mapstructure:"FIREWALL_DEBUG,omitempty"`
	OutboundSubnets []string      `mapstructure:"FIREWALL_OUTBOUND_SUBNETS,omitempty"`
}

type HealthRestartVpn string

const (
	HealthRestartVpnOn  HealthRestartVpn = "on"
	HealthRestartVpnOff HealthRestartVpn = "off"
)

type HealthCheck struct {
	TargetAddresses []string         `mapstructure:"HEALTHCHECK_TARGET_ADDRESSES,omitempty"`
	IcmpTargetIps   []string         `mapstructure:"HEALTHCHECK_ICMP_TARGET_IPS,omitempty"`
	ServerAddress   string           `mapstructure:"HEALTHCHECK_SERVER_ADDRESS,omitempty"`
	RestartVpn      HealthRestartVpn `mapstructure:"HEALTHCHECK_RESTART_VPN,omitempty"`
}

type HttpProxyStatus string

const (
	HttpProxyOn  HttpProxyStatus = "on"
	HttpProxyOff HttpProxyStatus = "off"
)

type HttpProxyLog string

const (
	HttpProxyLogOn  HttpProxyLog = "on"
	HttpProxyLogOff HttpProxyLog = "off"
)

type HttpProxyStealth string

const (
	HttpProxyStealthOn  HttpProxyStealth = "on"
	HttpProxyStealthOff HttpProxyStealth = "off"
)

type HttpProxy struct {
	Status           HttpProxyStatus  `mapstructure:"HTTPPROXY_STATUS,omitempty"`
	Log              HttpProxyLog     `mapstructure:"HTTPPROXY_LOG,omitempty"`
	ListeningAddress string           `mapstructure:"HTTPPROXY_LISTENING_ADDRESS,omitempty"`
	User             string           `mapstructure:"HTTPPROXY_USER,omitempty"`
	Password         string           `mapstructure:"HTTPPROXY_PASSWORD,omitempty"`
	Stealth          HttpProxyStealth `mapstructure:"HTTPPROXY_STEALTH,omitempty"`
}

type OpenVPNProtocol string

const (
	OpenVPNProtocolTcp OpenVPNProtocol = "tcp"
	OpenVPNProtocolUdp OpenVPNProtocol = "udp"
)

type OpenVPNVersion float64

const (
	OpenVPNVersion2_5 OpenVPNVersion = 2.5
	OpenVPNVersion2_6 OpenVPNVersion = 2.6
)

type OpenVPNVerbosity int

const (
	OpenVPNVerbosity0 OpenVPNVerbosity = 0
	OpenVPNVerbosity1 OpenVPNVerbosity = 1
	OpenVPNVerbosity2 OpenVPNVerbosity = 2
	OpenVPNVerbosity3 OpenVPNVerbosity = 3
	OpenVPNVerbosity4 OpenVPNVerbosity = 4
	OpenVPNVerbosity5 OpenVPNVerbosity = 5
	OpenVPNVerbosity6 OpenVPNVerbosity = 6
)

type OpenVPNRoot string

const (
	OpenVPNRootOn  OpenVPNRoot = "on"
	OpenVPNRootOff OpenVPNRoot = "off"
)

type OpenVPN struct {
	User     string `mapstructure:"OPENVPN_USER"`
	Password string `mapstructure:"OPENVPN_PASSWORD"`

	Protocol      OpenVPNProtocol  `mapstructure:"OPENVPN_PROTOCOL,omitempty"`
	Version       OpenVPNVersion   `mapstructure:"OPENVPN_VERSION,omitempty"`
	EndpointIP    string           `mapstructure:"OPENVPN_ENDPOINT_IP,omitempty"`
	EndpointPort  int              `mapstructure:"OPENVPN_ENDPOINT_PORT,omitempty"`
	Verbosity     OpenVPNVerbosity `mapstructure:"OPENVPN_VERBOSITY,omitempty"`
	Flags         []string         `mapstructure:"OPENVPN_FLAGS,omitempty"`
	Root          OpenVPNRoot      `mapstructure:"OPENVPN_ROOT,omitempty"`
	Ciphers       string           `mapstructure:"OPENVPN_CIPHERS,omitempty"`
	Auth          string           `mapstructure:"OPENVPN_AUTH,omitempty"`
	Mssfix        int              `mapstructure:"OPENVPN_MSSFIX,omitempty"`
	Cert          string           `mapstructure:"OPENVPN_CERT,omitempty"`
	Key           string           `mapstructure:"OPENVPN_KEY,omitempty"`
	EncryptedKey  string           `mapstructure:"OPENVPN_ENCRYPTED_KEY,omitempty"`
	KeyPassphrase string           `mapstructure:"OPENVPN_KEY_PASSPHRASE,omitempty"`
	ProcessUser   string           `mapstructure:"OPENVPN_PROCESS_USER,omitempty"`
	CustomConfig  string           `mapstructure:"OPENVPN_CUSTOM_CONFIG,omitempty"`
}

type GluetunLogLevel string

const (
	GluetunLogLevelDebug   GluetunLogLevel = "debug"
	GluetunLogLevelInfo    GluetunLogLevel = "info"
	GluetunLogLevelWarning GluetunLogLevel = "warn"
	GluetunLogLevelError   GluetunLogLevel = "error"
)

type PublicIPEnabled bool

const (
	PublicIPEnabledTrue  PublicIPEnabled = true
	PublicIPEnabledFalse PublicIPEnabled = false
)

type PublicIPAPI string

const (
	PublicIPAPIIpinfo      PublicIPAPI = "ipinfo"
	PublicIPAPIIfconfico   PublicIPAPI = "ifconfigco"
	PublicIPAPIIp2location PublicIPAPI = "ip2location"
	PublicIPAPICloudflare  PublicIPAPI = "cloudflare"
)

type VersionInformation string

const (
	VersionInformationOn  VersionInformation = "on"
	VersionInformationOff VersionInformation = "off"
)

type BoringPollGluetunCom string

const (
	BoringPollGluetunComOn  BoringPollGluetunCom = "on"
	BoringPollGluetunComOff BoringPollGluetunCom = "off"
)

type Other struct {
	LogLevel             GluetunLogLevel      `mapstructure:"LOG_LEVEL,omitempty"`
	Tz                   string               `mapstructure:"TZ,omitempty"`
	Puid                 int                  `mapstructure:"PUID,omitempty"`
	Pgid                 int                  `mapstructure:"PGID,omitempty"`
	PublicIPEnabled      PublicIPEnabled      `mapstructure:"PUBLIC_IP_ENABLED,omitempty"`
	PublicIPAPI          []PublicIPAPI        `mapstructure:"PUBLIC_IP_API,omitempty"`
	PublicIPAPIToken     string               `mapstructure:"PUBLIC_IP_API_TOKEN,omitempty"`
	PublicIPFile         string               `mapstructure:"PUBLIC_IP_FILE,omitempty"`
	VersionInformation   VersionInformation   `mapstructure:"VERSION_INFORMATION,omitempty"`
	BoringPollGluetunCom BoringPollGluetunCom `mapstructure:"BORING_POLL_GLUETUN_COM,omitempty"`
}

type VPNPortForwarding string

const (
	VPNPortForwardingOn  VPNPortForwarding = "on"
	VPNPortForwardingOff VPNPortForwarding = "off"
)

type VPNPortForwardingProvider string

const (
	PortForwardingPrivate        VPNPortForwardingProvider = "private"
	PortForwardingInternetAccess VPNPortForwardingProvider = "internet access"
	PortForwardingPerfectPrivacy VPNPortForwardingProvider = "perfect privacy"
	PortForwardingPrivateVPN     VPNPortForwardingProvider = "privatevpn"
	PortForwardingProtonVPN      VPNPortForwardingProvider = "protonvpn"
)

type VPNServerPortForwardingOptions struct {
	VPNPortForwarding VPNPortForwarding         `mapstructure:"VPN_PORT_FORWARDING,omitempty"`
	Provider          VPNPortForwardingProvider `mapstructure:"VPN_PORT_FORWARDING_PROVIDER,omitempty"`
	StatusFile        string                    `mapstructure:"VPN_PORT_FORWARDING_STATUS_FILE,omitempty"`
	ListeningPort     int                       `mapstructure:"VPN_PORT_FORWARDING_LISTENING_PORT,omitempty"`
	UpCommand         string                    `mapstructure:"VPN_PORT_FORWARDING_UP_COMMAND,omitempty"`
	DownCommand       string                    `mapstructure:"VPN_PORT_FORWARDING_DOWN_COMMAND,omitempty"`
}

type ShadowSocks string

const (
	ShadowSocksOn  ShadowSocks = "on"
	ShadowSocksOff ShadowSocks = "off"
)

type ShadowSocksLog string

const (
	ShadowSocksLogOn  ShadowSocksLog = "on"
	ShadowSocksLogOff ShadowSocksLog = "off"
)

type ShadowSocksCipher string

const (
	ShadowSocksCipherChaCha20IetfPoly1305 ShadowSocksCipher = "chacha20-ietf-poly1305"
	ShadowSocksCipherAes128Gcm            ShadowSocksCipher = "aes-128-gcm"
	ShadowSocksCipherAes256Gcm            ShadowSocksCipher = "aes-256-gcm"
)

type ShadowSocksOptions struct {
	ShadowSocks      ShadowSocks       `mapstructure:"SHADOWSOCKS,omitempty"`
	Log              ShadowSocksLog    `mapstructure:"SHADOWSOCKS_LOG,omitempty"`
	ListeningAddress string            `mapstructure:"SHADOWSOCKS_LISTENING_ADDRESS,omitempty"`
	Password         string            `mapstructure:"SHADOWSOCKS_PASSWORD,omitempty"`
	Cipher           ShadowSocksCipher `mapstructure:"SHADOWSOCKS_CIPHER,omitempty"`
}

type StorageOptions struct {
	FilePath string `mapstructure:"STORAGE_FILE_PATH,omitempty"`
}

type ServerUpdaterOptions struct {
	Period              string        `mapstructure:"UPDATER_PERIOD,omitempty"`
	MinRatio            float64       `mapstructure:"UPDATER_MIN_RATIO,omitempty"`
	VPNServiceProviders []VPNProvider `mapstructure:"UPDATER_VPN_SERVICE_PROVIDERS,omitempty"`
	ProtonVPNEmail      string        `mapstructure:"PROTONVPN_EMAIL,omitempty"`
	ProtonVPNPassword   string        `mapstructure:"PROTONVPN_PASSWORD,omitempty"`
}

type VPNOptions struct {
	ServiceProvider VPNProvider        `mapstructure:"VPN_SERVICE_PROVIDER,omitempty"`
	Type            config.VPNProtocol `mapstructure:"VPN_TYPE,omitempty"`
	Interface       string             `mapstructure:"VPN_INTERFACE,omitempty"`
	UpCommand       string             `mapstructure:"VPN_UP_COMMAND,omitempty"`
	DownCommand     string             `mapstructure:"VPN_DOWN_COMMAND,omitempty"`
}

type WireguardImplementation string

const (
	Auto        WireguardImplementation = "auto"
	KernelSpace WireguardImplementation = "kernelspace"
	UserSpace   WireguardImplementation = "userspace"
)

type WireguardOptions struct {
	PrivateKey   string   `mapstructure:"WIREGUARD_PRIVATE_KEY"`
	Addresses    []string `mapstructure:"WIREGUARD_ADDRESSES"`
	PublicKey    string   `mapstructure:"WIREGUARD_PUBLIC_KEY"`
	EndpointIP   string   `mapstructure:"WIREGUARD_ENDPOINT_IP"`
	EndpointPort int      `mapstructure:"WIREGUARD_ENDPOINT_PORT"`

	PresharedKey                string                  `mapstructure:"WIREGUARD_PRESHARED_KEY,omitempty"`
	AllowedIPs                  []string                `mapstructure:"WIREGUARD_ALLOWED_IPS,omitempty"`
	Implementation              WireguardImplementation `mapstructure:"WIREGUARD_IMPLEMENTATION,omitempty"`
	MTU                         int                     `mapstructure:"WIREGUARD_MTU,omitempty"`
	PersistentKeepaliveInterval string                  `mapstructure:"WIREGUARD_PERSISTENT_KEEPALIVE_INTERVAL,omitempty"`
}
