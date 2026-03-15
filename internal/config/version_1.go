package config

import (
	z "github.com/Oudwins/zog"
)

type SchemaVersion struct {
	SchemaVersion ConfigVersion `zog:"schema_version"`
}

var VPNSchemaVersion = z.Struct(z.Shape{
	"SchemaVersion": z.IntLike[ConfigVersion]().Default(DEFAULT_CONFIG_VERSION),
})

type ConfigCredential struct {
	From string `zog:"username"`
	Name string `zog:"password"`
}

var ConfigCredentialSchema = z.Struct(z.Shape{
	"From": z.String(),
	"Name": z.String(),
})

type ConfigVersionOne struct {
	SchemaVersion  ConfigVersion `zog:"schema_version"`
	GluetunVersion string        `zog:"gluetun_version"`
	VPNConfig      struct {
		Provider           string   `zog:"provider"`
		Protocol           Protocol `zog:"protocol"`
		OpenVPNCredentials struct {
			Username ConfigCredential `zog:"username"`
			Password ConfigCredential `zog:"password"`
		} `zog:"openvpn_credentials"`
		WireguardCredentials struct {
			PrivateKey   ConfigCredential `zog:"private_key"`
			Addresses    []string         `zog:"addresses"`
			PublicKey    ConfigCredential `zog:"public_key"`
			EndpointIP   string           `zog:"endpoint_ip"`
			EndpointPort int              `zog:"endpoint_port"`
		} `zog:"wireguard_credentials"`
	} `zog:"vpn_config"`
	Log struct {
		Level       string `zog:"level"`
		Destination string `zog:"destination"`
		Filename    string `zog:"filename"`
	} `zog:"log"`
	EnvVars []string `zog:"env_vars"`
}

var VPNConfigVersionOneSchema = z.Struct(z.Shape{
	"SchemaVersion":  z.IntLike[ConfigVersion]().Default(DEFAULT_CONFIG_VERSION),
	"GluetunVersion": z.String(),
	"VPNConfig": z.Struct(z.Shape{
		"Provider": z.String(),
		"Protocol": z.StringLike[Protocol]().Default(DEFAULT_VPN_PROTOCOL).OneOf([]Protocol{ProtocolWireguard, ProtocolOpenVPN}),
		"OpenVPNCredentials": z.Struct(z.Shape{
			"Username": ConfigCredentialSchema,
			"Password": ConfigCredentialSchema,
		}),
		"WireguardCredentials": z.Struct(z.Shape{
			"PrivateKey":   ConfigCredentialSchema,
			"Addresses":    z.Slice(z.String()),
			"PublicKey":    ConfigCredentialSchema,
			"EndpointIP":   z.String(),
			"EndpointPort": z.Int(),
		}),
	}),
	"Log": z.Struct(z.Shape{
		"Level":       z.String().Default(DEFAULT_LOG_LEVEL).OneOf([]string{LOG_LEVEL_DEBUG, LOG_LEVEL_INFO, LOG_LEVEL_WARNING, LOG_LEVEL_ERROR}),
		"Destination": z.String(),
		"Filename":    z.String(),
	}),
	"EnvVars": z.Slice(z.String()),
})
