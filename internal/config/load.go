package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	z "github.com/Oudwins/zog"
	"github.com/Oudwins/zog/parsers/zjson"
)

func loadFromFile[cfg any](path string, schema *z.StructSchema) (*cfg, []error) {
	ext := filepath.Ext(path)
	if ext != ".json" {
		return nil, []error{fmt.Errorf("unsupported config file extension: %s", ext)}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, []error{fmt.Errorf("failed to read config file %s: %w", path, err)}
	}

	return loadFromBytes[cfg](data, schema)
}

func loadFromBytes[cfg any](data []byte, schema *z.StructSchema) (*cfg, []error) {
	configVersion, versionErrs := getConfigVersion(data, VPNSchemaVersion)
	if len(versionErrs) > 0 {
		errors := make([]error, len(versionErrs))
		for i, versionErr := range versionErrs {
			errors[i] = fmt.Errorf("config validation failed: %s: %w", versionErr.Message, versionErr.Err)
		}
		return nil, errors
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

func getConfigVersion(data []byte, versionSchema *z.StructSchema) (ConfigVersion, z.ZogIssueList) {
	var schemaVersion SchemaVersion
	zogErrs := versionSchema.Parse(zjson.Decode(bytes.NewReader(data)), &schemaVersion)

	if len(zogErrs) > 0 {
		return 0, zogErrs
	}

	configVersion := schemaVersion.SchemaVersion

	return configVersion, nil
}

func LoadAppConfigFromFile(path string) (*ConfigVersionOne, []error) {
	return loadFromFile[ConfigVersionOne](path, VPNConfigVersionOneSchema)
}

func LoadAppConfigFromBytes(data []byte) (*ConfigVersionOne, []error) {
	return loadFromBytes[ConfigVersionOne](data, VPNConfigVersionOneSchema)
}

func GetAppConfigVersion(data []byte) (ConfigVersion, z.ZogIssueList) {
	return getConfigVersion(data, VPNSchemaVersion)
}
