package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"go.yaml.in/yaml/v4"
)

type Schema interface {
	CLIConfig | DaemonConfig
}

type Loader[T Schema] interface {
	Load(path string) (*T, error)
}

type YAMLLoader[T Schema] struct {
	DefaultConfig *T
}

func NewYAMLLoader[T Schema](defaultConfig *T) *YAMLLoader[T] {
	return &YAMLLoader[T]{
		DefaultConfig: defaultConfig,
	}
}

func NewCLILoader() *YAMLLoader[CLIConfig] {
	return NewYAMLLoader(DefaultCLIConfig())
}

func NewDaemonLoader() *YAMLLoader[DaemonConfig] {
	return NewYAMLLoader(DefaultDaemonConfig())
}

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

func mergeInto(dest, src reflect.Value) error {
	if src.Kind() != dest.Kind() {
		return nil
	}

	switch dest.Kind() {
	case reflect.Struct:

		for i := range dest.NumField() {
			destField := dest.Field(i)
			srcField := src.Field(i)

			// Skip unexported fields
			if !destField.CanSet() {
				continue
			}

			// If source field is a zero value, keep destination value
			if srcField.IsZero() {
				continue
			}

			// If source field is a pointer to zero value, keep destination value
			if srcField.Kind() == reflect.Pointer && !srcField.IsNil() {
				if srcField.Elem().IsZero() {
					continue
				}
			}

			// Handle nested structs recursively
			if destField.Kind() == reflect.Struct {
				if err := mergeInto(destField, srcField); err != nil {
					return err
				}
				continue
			}

			// Handle pointer types
			if destField.Kind() == reflect.Pointer && srcField.Kind() == reflect.Pointer {
				if srcField.IsNil() {
					continue // Keep dest if src is nil
				}
				// Create new pointer and copy value
				destField.Set(reflect.New(destField.Type().Elem()))
				if err := mergeInto(destField.Elem(), srcField.Elem()); err != nil {
					return err
				}
				continue
			}

			// For simple types, just copy the value
			destField.Set(srcField)
		}

	case reflect.Slice:
		// If source slice has elements, replace destination slice
		if src.Len() > 0 {
			dest.Set(src)
		}

	case reflect.Map:
		// If source map has entries, replace destination map
		if src.Len() > 0 {
			dest.Set(src)
		}
	}

	return nil
}
