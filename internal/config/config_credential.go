package config

import (
	"fmt"
	"os"
	"strings"

	z "github.com/Oudwins/zog"
)

var ConfigCredentialSchema = z.Struct(z.Shape{
	"From": z.StringLike[configCredentialSource]().OneOf([]configCredentialSource{ConfigCredentialSourceFile, ConfigCredentialSourceEnv}),
	"Name": z.String(),
})

type configCredentialSource string

const (
	ConfigCredentialSourceFile configCredentialSource = "file"
	ConfigCredentialSourceEnv  configCredentialSource = "env"
)

type ConfigCredential struct {
	From configCredentialSource `zog:"from"`
	Name string                 `zog:"name"`
}

func (c ConfigCredential) LoadCredential() (value string, err error) {
	switch c.From {
	case ConfigCredentialSourceFile:
		envValue, err := os.ReadFile(c.Name)
		if err != nil {
			return "", fmt.Errorf("failed to read credential file: %w", err)
		}

		value = strings.TrimRight(string(envValue), "\r\n")

		return value, nil

	case ConfigCredentialSourceEnv:
		envValue := os.Getenv(c.Name)
		if envValue == "" {
			return "", fmt.Errorf("environment variable %s not set", c.Name)
		}

		return envValue, nil

	default:
		return "", fmt.Errorf("unknown credential source: %s", c.From)
	}
}
