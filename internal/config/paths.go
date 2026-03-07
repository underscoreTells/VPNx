package config

import (
	"fmt"
	"os"
)

func newDirUnavailableError(err error) error {
	return fmt.Errorf("directory unavailable: %w", err)
}

func getDir(pathEnvVar string, fallbackHome string, fallbackPath string) (string, error) {
	if v := os.Getenv(pathEnvVar); v != "" {
		return v, nil
	}

	info, err := os.Stat(os.Getenv(fallbackHome) + fallbackPath)
	if err != nil {
		return "", newDirUnavailableError(err)
	}
	if !info.IsDir() {
		return "", newDirUnavailableError(fmt.Errorf("not a directory"))
	}
	return fallbackHome + fallbackPath, nil
}

func GetConfigDir() (string, error) {
	return getDir(XDG_CONFIG_HOME, XDG_HOME, XDG_CONFIG)
}

func GetCacheDir() (string, error) {
	return getDir(XDG_CACHE_HOME, XDG_HOME, XDG_CACHE)
}

func GetDataDir() (string, error) {
	return getDir(XDG_DATA_HOME, XDG_HOME, XDG_DATA)
}

func GetStateDir() (string, error) {
	return getDir(XDG_STATE_HOME, XDG_HOME, XDG_STATE)
}
