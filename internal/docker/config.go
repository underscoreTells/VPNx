package docker

import "github.com/moby/moby/api/types/container"

type ClientConfigOption func(*container.Config)

type HostConfigOption func(*container.HostConfig)

func CreateClientConfig(options ...ClientConfigOption) *container.Config {
	config := &container.Config{}
	for _, opt := range options {
		opt(config)
	}
	return config
}

func CreateHostConfig(options ...HostConfigOption) *container.HostConfig {
	config := &container.HostConfig{}
	for _, opt := range options {
		opt(config)
	}
	return config
}
