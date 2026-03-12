package config

import z "github.com/Oudwins/zog"

type SchemaVersion struct {
	SchemaVersion ConfigVersion `zog:"schema_version"`
}

var VPNSchemaVersion = z.Struct(z.Shape{
	"SchemaVersion": z.IntLike[ConfigVersion]().Default(DEFAULT_CONFIG_VERSION),
})

type ConfigVersionOne struct {
	SchemaVersion  ConfigVersion `zog:"schema_version"`
	GluetunVersion string        `zog:"gluetun_version"`
	VPNConfig      struct {
		Provider    string  `zog:"provider"`
		Type        VPNType `zog:"type"`
		Credentials struct {
			Username struct {
				From string `zog:"from"`
				Name string `zog:"name"`
			} `zog:"username"`
			Password struct {
				From string `zog:"from"`
				Name string `zog:"name"`
			} `zog:"password"`
		} `zog:"credentials"`
	} `zog:"vpn_config"`
	Log struct {
		Level       string `zog:"level"`
		Destination string `zog:"destination"`
		Filename    string `zog:"filename"`
	} `zog:"log"`
}

var VPNConfigVersionOneSchema = z.Struct(z.Shape{
	"SchemaVersion":  z.IntLike[ConfigVersion]().Default(DEFAULT_CONFIG_VERSION),
	"GluetunVersion": z.String(),
	"VPNConfig": z.Struct(z.Shape{
		"Provider": z.String(),
		"Type":     z.StringLike[VPNType]().Default(DEFAULT_VPN_TYPE).OneOf([]VPNType{VPN_TYPE_WIREGUARD, VPN_TYPE_OPENVPN}),
		"Credentials": z.Struct(z.Shape{
			"Username": z.Struct(z.Shape{
				"From": z.String(),
				"Name": z.String(),
			}),
			"Password": z.Struct(z.Shape{
				"From": z.String(),
				"Name": z.String(),
			}),
		}),
	}),
	"Log": z.Struct(z.Shape{
		"Level":       z.String().Default(DEFAULT_LOG_LEVEL).OneOf([]string{LOG_LEVEL_DEBUG, LOG_LEVEL_INFO, LOG_LEVEL_WARNING, LOG_LEVEL_ERROR}),
		"Destination": z.String(),
		"Filename":    z.String(),
	}),
})
