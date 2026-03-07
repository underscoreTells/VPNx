package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/parsers/zjson"
)

// LoadConfigFromFile loads a config file from disk, validates it against the provided
// schema, and returns the parsed config along with any validation errors encountered.
func LoadConfigFromFile(path string, schema *z.StructSchema) (*CoreConfig, []error) {
	ext := filepath.Ext(path)
	if ext != ".json" {
		return nil, []error{fmt.Errorf("unsupported config file extension: %s", ext)}
	}

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, []error{fmt.Errorf("failed to read config file %s: %w", path, err)}
	}

	var cfg CoreConfig
	zogErrs := schema.Parse(zjson.Decode(bytes.NewReader(data)), &cfg)

	if len(zogErrs) > 0 {
		errors := make([]error, len(zogErrs))
		for i, zogErr := range zogErrs {
			errors[i] = fmt.Errorf("config validation failed: %s: %w", zogErr.Message, zogErr.Err)
		}
		return nil, errors
	}

	return &cfg, nil
}

// LoadCoreConfigFromFile loads a core config file from disk, validates it against the
// CoreConfigSchema, and returns the parsed config along with any validation errors encountered.
func LoadCoreConfigFromFile(path string) (*CoreConfig, []error) {
	return LoadConfigFromFile(path, CoreConfigSchema)
}
