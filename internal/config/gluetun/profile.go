package gluetun

import (
	"github.com/underscoreTells/vpn-exit-node/internal/config/gluetun/options"
)

type GluetunProfile struct {
	VPN       options.VPN
	OpenVPN   *options.OpenVPN
	Wireguard *options.Wireguard

	ControlServer  options.ControlServer
	DNS            options.DNS
	Firewall       options.Firewall
	HealthCheck    options.HealthCheck
	HttpProxy      options.HttpProxy
	Other          options.Other
	PortForwarding options.VPNPortForwardingStatus
	ShadowSocks    options.ShadowSocks
	Storage        options.Storage
	Updater        options.ServerUpdater
}
