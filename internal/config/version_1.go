package config

import z "github.com/Oudwins/zog"

type ConfigVersionOne struct {
	SchemaVersion  ConfigVersion
	GluetunVersion string
	VPNConfig      struct {
		Provider    string
		Type        VPNType
		Credentials struct {
			Username struct {
				From string
				name string
			}
			Password struct {
				From string
				name string
			}
		}
	}
	Log struct {
		Level       string
		Destination string
		FileName    string
	}
}

var VPNConfigVersionOneSchema = z.Struct(z.Shape{
	"gluetun_version": z.String(),
	"vpn_config": z.Struct(z.Shape{
		"provider": z.String(),
		"type":     z.String().OneOf([]string{string(VPN_TYPE_WIREGUARD), string(VPN_TYPE_OPENVPN)}),
		"credentials": z.Struct(z.Shape{
			"username": z.Struct(z.Shape{
				"from": z.String(),
				"name": z.String(),
			}),
			"password": z.Struct(z.Shape{
				"from": z.String(),
				"name": z.String(),
			}),
		}),
	}),
	"log": z.Struct(z.Shape{
		"level":       z.String().Default(DEFAULT_LOG_LEVEL).OneOf([]string{LOG_LEVEL_DEBUG, LOG_LEVEL_INFO, LOG_LEVEL_WARNING, LOG_LEVEL_ERROR}),
		"destination": z.String(),
		"filename":    z.String(),
	}),
})
