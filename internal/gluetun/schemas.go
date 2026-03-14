package gluetun

import (
	z "github.com/Oudwins/zog"
	"github.com/underscoreTells/vpn-exit-node/internal/config"
)

var vpnProviderValues = []VPNProvider{
	AirVPN,
	Custom,
	CyberGhost,
	ExpressVPN,
	FastestVPN,
	Giganews,
	Hidemyass,
	IPvanish,
	Ivpn,
	Mullvad,
	NordVPN,
	PerfectPrivacy,
	Privado,
	PrivateInternetAccess,
	PrivateVPN,
	ProtonVPN,
	PureVPN,
	SlickVPN,
	Surfshark,
	Torguard,
	VPNSecure,
	VPNUnlimited,
	VyprVPN,
	Windscribe,
}

var dnsUpstreamResolverValues = []DnsUpstreamResolver{
	CircaFamily,
	CircaPrivate,
	CircaProtected,
	CleanBrowsing,
	Adult,
	CleanBrowsingFamily,
	CleanBrowsingSecurity,
	Cloudflare,
	CloudflareFamily,
	CloudflareSecurity,
	Google,
	Libredns,
	Opendns,
	Quad9,
	Quad9Secured,
	Quad9Unsecured,
	Quadrant,
}

var publicIPAPIValues = []PublicIPAPI{
	PublicIPAPIIpinfo,
	PublicIPAPIIfconfico,
	PublicIPAPIIp2location,
	PublicIPAPICloudflare,
}

var VPNProviderSchema = stringEnumSchema(vpnProviderValues)

var HttpControlServerLogSchema = stringEnumSchema([]HttpControlServerLog{
	HttpControlServerLogOn,
	HttpControlServerLogOff,
})

var ControlServerOptionsSchema = z.Struct(z.Shape{
	"Address":            z.String().TestFunc(isValidListeningAddress),
	"Log":                HttpControlServerLogSchema,
	"AuthConfigFilepath": z.String().TestFunc(isValidFilePath),
})

var DnsUpstreamResolverTypeSchema = stringEnumSchema([]DnsUpstreamResolverType{
	DnsUpstreamResolverTypeDot,
	DnsUpstreamResolverTypeDoh,
	DnsUpstreamResolverTypePlain,
})

var DnsUpstreamResolverSchema = stringEnumSchema(dnsUpstreamResolverValues)

var DnsCachingSchema = stringEnumSchema([]DnsCaching{
	DnsCachingOn,
	DnsCachingOff,
})

var DnsUpstreamIPV6Schema = stringEnumSchema([]DnsUpstreamIPV6{
	DnsUpstreamIPV6On,
	DnsUpstreamIPV6Off,
})

var BlockMaliciousSchema = stringEnumSchema([]BlockMalicious{
	BlockMaliciousOn,
	BlockMaliciousOff,
})

var BlockSurveillanceSchema = stringEnumSchema([]BlockSurveillance{
	BlockSurveillanceOn,
	BlockSurveillanceOff,
})

var BlockAdsSchema = stringEnumSchema([]BlockAds{
	BlockAdsOn,
	BlockAdsOff,
})

var DNSSchema = z.Struct(z.Shape{
	"UpstreamResolverType":               DnsUpstreamResolverTypeSchema,
	"UpstreamResolvers":                  z.Slice(stringEnumSchema(dnsUpstreamResolverValues).Required()),
	"Caching":                            DnsCachingSchema,
	"UpstreamIPV6":                       DnsUpstreamIPV6Schema,
	"BlockIps":                           z.Slice(z.String().IP()),
	"BlockIpPrefixes":                    z.Slice(z.String()),
	"RebindingProtectionExemptHostnames": z.Slice(z.String().TestFunc(isValidHostname)),
	"UpdatePeriod":                       z.String().TestFunc(isValidTimePeriod),
	"BlockMalicious":                     BlockMaliciousSchema,
	"BlockSurveillance":                  BlockSurveillanceSchema,
	"BlockAds":                           BlockAdsSchema,
	"UnblockHostnames":                   z.Slice(z.String().TestFunc(isValidHostname)),
	"UpstreamPlainAddresses":             z.Slice(z.String().TestFunc(isValidPlainAddress)),
})

var FirewallDebugSchema = stringEnumSchema([]FirewallDebug{
	FirewallDebugOn,
	FirewallDebugOff,
})

var FirewallSchema = z.Struct(z.Shape{
	"VpnInputPorts":   z.Slice(z.Int().GTE(0).LTE(65535)),
	"InputPorts":      z.Slice(z.Int().GTE(0).LTE(65535)),
	"Debug":           FirewallDebugSchema,
	"OutboundSubnets": z.Slice(z.String().TestFunc(isValidSubnet)),
})

var HealthRestartVpnSchema = stringEnumSchema([]HealthRestartVpn{
	HealthRestartVpnOn,
	HealthRestartVpnOff,
})

var HealthCheckSchema = z.Struct(z.Shape{
	"TargetAddresses": z.Slice(z.String()),
	"IcmpTargetIps":   z.Slice(z.String().IPv4()),
	"ServerAddress":   z.String().TestFunc(isValidPlainAddress),
	"RestartVpn":      HealthRestartVpnSchema,
})

var HttpProxyStatusSchema = stringEnumSchema([]HttpProxyStatus{
	HttpProxyOn,
	HttpProxyOff,
})

var HttpProxyLogSchema = stringEnumSchema([]HttpProxyLog{
	HttpProxyLogOn,
	HttpProxyLogOff,
})

var HttpProxyStealthSchema = stringEnumSchema([]HttpProxyStealth{
	HttpProxyStealthOn,
	HttpProxyStealthOff,
})

var HttpProxySchema = z.Struct(z.Shape{
	"Status":           HttpProxyStatusSchema,
	"Log":              HttpProxyLogSchema,
	"ListeningAddress": z.String().TestFunc(isValidListeningAddress),
	"User":             z.String(),
	"Password":         z.String(),
	"Stealth":          HttpProxyStealthSchema,
})

var OpenVPNProtocolSchema = stringEnumSchema([]OpenVPNProtocol{
	OpenVPNProtocolTcp,
	OpenVPNProtocolUdp,
})

var OpenVPNVersionSchema = floatEnumSchema([]OpenVPNVersion{
	OpenVPNVersion2_5,
	OpenVPNVersion2_6,
})

var OpenVPNVerbositySchema = intEnumSchema([]OpenVPNVerbosity{
	OpenVPNVerbosity0,
	OpenVPNVerbosity1,
	OpenVPNVerbosity2,
	OpenVPNVerbosity3,
	OpenVPNVerbosity4,
	OpenVPNVerbosity5,
	OpenVPNVerbosity6,
})

var OpenVPNRootSchema = stringEnumSchema([]OpenVPNRoot{
	OpenVPNRootOn,
	OpenVPNRootOff,
})

var OpenVPNSchema = z.Struct(z.Shape{
	"User":          z.String().Required(),
	"Password":      z.String().Required(),
	"Protocol":      OpenVPNProtocolSchema,
	"Version":       OpenVPNVersionSchema,
	"EndpointIP":    z.String().IP(),
	"EndpointPort":  z.Int().GTE(0).LTE(65535),
	"Verbosity":     OpenVPNVerbositySchema,
	"Flags":         z.Slice(z.String()),
	"Root":          OpenVPNRootSchema,
	"Ciphers":       z.String(),
	"Auth":          z.String(),
	"Mssfix":        z.Int().GTE(0).LTE(9999),
	"Cert":          z.String().TestFunc(isValidBase64PEM),
	"Key":           z.String().TestFunc(isValidBase64PEM),
	"EncryptedKey":  z.String().TestFunc(isValidBase64PEM),
	"KeyPassphrase": z.String(),
	"ProcessUser":   z.String(),
	"CustomConfig":  z.String().TestFunc(isValidFilePath),
})

var GluetunLogLevelSchema = stringEnumSchema([]GluetunLogLevel{
	GluetunLogLevelDebug,
	GluetunLogLevelInfo,
	GluetunLogLevelWarning,
	GluetunLogLevelError,
})

var PublicIPEnabledSchema = z.BoolLike[PublicIPEnabled]()

var PublicIPAPISchema = stringEnumSchema(publicIPAPIValues)

var VersionInformationSchema = stringEnumSchema([]VersionInformation{
	VersionInformationOn,
	VersionInformationOff,
})

var BoringPollGluetunComSchema = stringEnumSchema([]BoringPollGluetunCom{
	BoringPollGluetunComOn,
	BoringPollGluetunComOff,
})

var OtherSchema = z.Struct(z.Shape{
	"LogLevel":             GluetunLogLevelSchema,
	"Tz":                   z.String(),
	"Puid":                 z.Int(),
	"Pgid":                 z.Int(),
	"PublicIPEnabled":      PublicIPEnabledSchema,
	"PublicIPAPI":          z.Slice(stringEnumSchema(publicIPAPIValues).Required()),
	"PublicIPAPIToken":     z.String(),
	"PublicIPFile":         z.String().TestFunc(isValidFilePath),
	"VersionInformation":   VersionInformationSchema,
	"BoringPollGluetunCom": BoringPollGluetunComSchema,
})

var VPNPortForwardingSchema = stringEnumSchema([]VPNPortForwarding{
	VPNPortForwardingOn,
	VPNPortForwardingOff,
})

var VPNPortForwardingProviderSchema = stringEnumSchema([]VPNPortForwardingProvider{
	PortForwardingPrivate,
	PortForwardingInternetAccess,
	PortForwardingPerfectPrivacy,
	PortForwardingPrivateVPN,
	PortForwardingProtonVPN,
})

var VPNServerPortForwardingOptionsSchema = z.Struct(z.Shape{
	"VPNPortForwarding": VPNPortForwardingSchema,
	"Provider":          VPNPortForwardingProviderSchema,
	"StatusFile":        z.String().TestFunc(isValidFilePath),
	"ListeningPort":     z.Int().GTE(0).LTE(65535),
	"UpCommand":         z.String(),
	"DownCommand":       z.String(),
})

var ShadowSocksSchema = stringEnumSchema([]ShadowSocks{
	ShadowSocksOn,
	ShadowSocksOff,
})

var ShadowSocksLogSchema = stringEnumSchema([]ShadowSocksLog{
	ShadowSocksLogOn,
	ShadowSocksLogOff,
})

var ShadowSocksCipherSchema = stringEnumSchema([]ShadowSocksCipher{
	ShadowSocksCipherChaCha20IetfPoly1305,
	ShadowSocksCipherAes128Gcm,
	ShadowSocksCipherAes256Gcm,
})

var ShadowSocksOptionsSchema = z.Struct(z.Shape{
	"ShadowSocks":      ShadowSocksSchema,
	"Log":              ShadowSocksLogSchema,
	"ListeningAddress": z.String().TestFunc(isValidListeningAddress),
	"Password":         z.String(),
	"Cipher":           ShadowSocksCipherSchema,
})

var StorageOptionsSchema = z.Struct(z.Shape{
	"FilePath": z.String().TestFunc(isValidFilePath),
})

var ServerUpdaterOptionsSchema = z.Struct(z.Shape{
	"Period":              z.String().TestFunc(isValidTimePeriod),
	"MinRatio":            z.Float64().GT(0).LTE(1),
	"VPNServiceProviders": z.Slice(stringEnumSchema(vpnProviderValues).Required()),
	"ProtonVPNEmail":      z.String().Email(),
	"ProtonVPNPassword":   z.String(),
})

var VPNOptionsSchema = z.Struct(z.Shape{
	"ServiceProvider": VPNProviderSchema,
	"Type": z.StringLike[config.VPNProtocol]().OneOf([]config.VPNProtocol{
		config.VPN_PROTOCOL_WIREGUARD,
		config.VPN_PROTOCOL_OPENVPN,
	}),
	"Interface":   z.String(),
	"UpCommand":   z.String(),
	"DownCommand": z.String(),
})

var WireguardImplementationSchema = stringEnumSchema([]WireguardImplementation{
	Auto,
	KernelSpace,
	UserSpace,
})

var WireguardOptionsSchema = z.Struct(z.Shape{
	"PrivateKey":                  z.String().Required(),
	"Addresses":                   z.Slice(z.String().TestFunc(isValidSubnet)).Required(),
	"PublicKey":                   z.String().Required(),
	"EndpointIP":                  z.String().Required().IP(),
	"EndpointPort":                z.Int().Required().GTE(0).LTE(65535),
	"PresharedKey":                z.String(),
	"AllowedIPs":                  z.Slice(z.String()),
	"Implementation":              WireguardImplementationSchema,
	"MTU":                         z.Int().GTE(0).LTE(1440),
	"PersistentKeepaliveInterval": z.String().TestFunc(isValidTimePeriod),
})
