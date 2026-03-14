package config

// === Enums ===

const (
	LOG_LEVEL_DEBUG   = "debug"
	LOG_LEVEL_INFO    = "info"
	LOG_LEVEL_WARNING = "warning"
	LOG_LEVEL_ERROR   = "error"
)

type VPNProtocol string

const (
	VPN_PROTOCOL_WIREGUARD VPNProtocol = "wireguard"
	VPN_PROTOCOL_OPENVPN   VPNProtocol = "openvpn"
)

type VPNCredentialSource string

const (
	CREDENTIAL_SOURCE_FILE VPNCredentialSource = "file"
	CREDENTIAL_SOURCE_ENV  VPNCredentialSource = "env"
)

type ConfigVersion int

const (
	CONFIG_VERSION_ONE ConfigVersion = 1
)
