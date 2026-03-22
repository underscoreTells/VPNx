package docker

import (
	"github.com/moby/moby/api/types/container"
)

func WithEnvMap(env map[string]string) ClientConfigOption {
	envArray := make([]string, 0, len(env))
	for k, v := range env {
		envArray = append(envArray, k+"="+v)
	}

	return func(cfg *container.Config) {
		cfg.Env = envArray
	}
}

func WithEnvArray(env []string) ClientConfigOption {
	return func(cfg *container.Config) {
		cfg.Env = env
	}
}

func WithClientOption(fn func(*container.Config)) ClientConfigOption {
	return func(cfg *container.Config) {
		fn(cfg)
	}
}
