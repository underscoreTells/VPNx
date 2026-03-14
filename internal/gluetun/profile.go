package gluetun

import "github.com/underscoreTells/vpn-exit-node/config"

type GluetunProfile struct {
	Provider VPNProvider
	Protocol config.VPNProtocol
}
