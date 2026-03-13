package gluetun

import "github.com/underscoreTells/vpn-exit-node/internal/config"

type GluetunProfile struct {
	Provider VPNProvider
	Protocol config.VPNProtocol
}
