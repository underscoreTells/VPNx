package common

type Protocol string

const (
	ProtocolWireguard Protocol = "wireguard"
	ProtocolOpenVPN   Protocol = "openvpn"
)
