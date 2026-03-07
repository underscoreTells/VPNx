package config

import (
	z "github.com/Oudwins/zog"
)

// === Enums ===

const (
	LOG_LEVEL_DEBUG   = "debug"
	LOG_LEVEL_INFO    = "info"
	LOG_LEVEL_WARNING = "warning"
	LOG_LEVEL_ERROR   = "error"
)

// === Types ===

type LogConfig struct {
	Level       string
	Destination string
	FileName    string
}

type CoreConfig struct {
	Log            LogConfig
	GluetunVersion string
}

// ===  Schemas ===

var LogConfigSchema = z.Struct(z.Shape{
	"level":       z.String().Default(DEFAULT_LOG_LEVEL).OneOf([]string{LOG_LEVEL_DEBUG, LOG_LEVEL_INFO, LOG_LEVEL_WARNING, LOG_LEVEL_ERROR}),
	"destination": z.String(),
	"filename":    z.String(),
})

var CoreConfigSchema = z.Struct(z.Shape{
	"log":             LogConfigSchema,
	"gluetun_version": z.String().Default(GLUETUN_TARGET_VERSION),
})
