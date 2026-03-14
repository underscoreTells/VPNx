package app

const (
	LOG_LEVEL_DEBUG   = "debug"
	LOG_LEVEL_INFO    = "info"
	LOG_LEVEL_WARNING = "warning"
	LOG_LEVEL_ERROR   = "error"
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
