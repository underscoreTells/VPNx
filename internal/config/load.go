package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/parsers/zjson"
)

// LoadConfigFromFile loads a config file from disk, validates it against the provided
// schema, and returns the parsed config along with any validation errors encountered.
func LoadConfigFromFile[cfg any](path string, schema *z.StructSchema) (*cfg, []error) {
	ext := filepath.Ext(path)
	if ext != ".json" {
		return nil, []error{fmt.Errorf("unsupported config file extension: %s", ext)}
	}

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, []error{fmt.Errorf("failed to read config file %s: %w", path, err)}
	}

	configVersion, err := GetConfigVersion(data)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to get config version: %w", err)}
	}

	rawConfig := ConfigVersions[configVersion]()
	config, ok := rawConfig.(cfg)
	if !ok {
		return nil, []error{fmt.Errorf("failed to convert config to type %T", config)}
	}

	zogErrs := schema.Parse(zjson.Decode(bytes.NewReader(data)), &config)

	if len(zogErrs) > 0 {
		errors := make([]error, len(zogErrs))
		for i, zogErr := range zogErrs {
			errors[i] = fmt.Errorf("config validation failed: %s: %w", zogErr.Message, zogErr.Err)
		}
		return nil, errors
	}

	return &config, nil
}

func GetConfigVersion(data []byte) (ConfigVersion, error) {
	var version ConfigVersion
	err := json.Unmarshal(data, &version)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal config version: %w", err)
	}
	return version, nil
}
