// Package config provides utilities for loading and merging configuration
// data. It includes generic loaders that read YAML files and merge the
// user-supplied values over a provided default configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"go.yaml.in/yaml/v4"
)

// Schema is a type constraint that limits which configuration structs can be
// used with the generic loaders in this package. It currently permits
// CLIConfig and DaemonConfig.
type Schema interface {
	CLIConfig | DaemonConfig
}

// Loader is a generic interface for loading configuration objects of type T
// from a file path.
type Loader[T Schema] interface {
	Load(path string) (*T, error)
}

// YAMLLoader is a generic loader that reads YAML configuration files and
// merges them into a provided default configuration. The loader retains a
// pointer to the default config used when merging.
type YAMLLoader[T Schema] struct {
	DefaultConfig *T
}

// NewYAMLLoader returns a YAMLLoader initialized with the provided default
// configuration. The returned loader will use the default when merging values
// absent from the loaded file.
func NewYAMLLoader[T Schema](defaultConfig *T) *YAMLLoader[T] {
	return &YAMLLoader[T]{
		DefaultConfig: defaultConfig,
	}
}

// NewCLILoader returns a YAMLLoader preconfigured with the package's default
// CLI configuration.
func NewCLILoader() *YAMLLoader[CLIConfig] {
	return NewYAMLLoader(DefaultCLIConfig())
}

// NewDaemonLoader returns a YAMLLoader preconfigured with the package's
// default daemon configuration.
func NewDaemonLoader() *YAMLLoader[DaemonConfig] {
	return NewYAMLLoader(DefaultDaemonConfig())
}

// Load reads and parses a YAML file at the given path into an instance of T,
// merges the parsed values into the loader's default configuration, and
// returns the resulting configuration. Errors are returned if the file cannot
// be read, parsed, or merged.
func (l *YAMLLoader[T]) Load(path string) (*T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config T
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	userConfig, err := l.overwriteDefault(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to merge config: %w", err)
	}

	return userConfig, nil
}

// overwriteDefault merges the provided config over the loader's default
// configuration and returns the merged result. It requires that a
// default configuration was provided when the loader was constructed.
func (l *YAMLLoader[T]) overwriteDefault(config *T) (*T, error) {
	if l.DefaultConfig == nil {
		return nil, errors.New("no default config provided")
	}

	if config == nil {
		return nil, errors.New("no config provided")
	}

	newConfig := *l.DefaultConfig

	if err := mergeInto(reflect.ValueOf(&newConfig).Elem(), reflect.ValueOf(*config)); err != nil {
		return nil, fmt.Errorf("failed to merge config: %w", err)
	}

	return &newConfig, nil
}

// mergeInto recursively merges non-zero values from src into dest using
// reflection. The function attempts to preserve destination values when
// the corresponding source values are zero-valued, and handles nested
// structs, pointers, slices, and maps appropriately.
func mergeInto(dest, src reflect.Value) error {
	if src.Kind() != dest.Kind() {
		return nil
	}

	switch dest.Kind() {
	case reflect.Struct:

		for i := range dest.NumField() {
			destField := dest.Field(i)
			srcField := src.Field(i)

			if !destField.CanSet() {
				continue
			}

			if srcField.Kind() != reflect.Bool && srcField.IsZero() {
				continue
			}

			if srcField.Kind() == reflect.Pointer && !srcField.IsNil() {
				if srcField.Elem().IsZero() {
					continue
				}
			}

			if destField.Kind() == reflect.Struct {
				if err := mergeInto(destField, srcField); err != nil {
					return err
				}
				continue
			}

			if destField.Kind() == reflect.Pointer && srcField.Kind() == reflect.Pointer {
				if srcField.IsNil() {
					continue
				}
				destField.Set(reflect.New(destField.Type().Elem()))
				if err := mergeInto(destField.Elem(), srcField.Elem()); err != nil {
					return err
				}
				continue
			}

			destField.Set(srcField)
		}

	// TODO: Implement deep merging for slices
	case reflect.Slice:
		if src.Len() > 0 {
			dest.Set(src)
		}

	// TODO: Implement deep merging for maps
	case reflect.Map:
		if src.Len() > 0 {
			dest.Set(src)
		}
	}

	return nil
}
