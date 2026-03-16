package config

import (
	"os"
	"path/filepath"
	"testing"
)

type dirGetterTestCase struct {
	name           string
	xdgEnvVar      string
	fallbackSuffix string
	get            func() (string, error)
}

func clearXDGEnv(t *testing.T) {
	t.Helper()

	for _, envVar := range []string{
		XDG_CONFIG_HOME,
		XDG_CACHE_HOME,
		XDG_DATA_HOME,
		XDG_STATE_HOME,
		XDG_HOME,
	} {
		t.Setenv(envVar, "")
	}
}

func TestXDGDirectoryGetters(t *testing.T) {
	testCases := []dirGetterTestCase{
		{name: "config", xdgEnvVar: XDG_CONFIG_HOME, fallbackSuffix: XDG_CONFIG, get: GetConfigDir},
		{name: "cache", xdgEnvVar: XDG_CACHE_HOME, fallbackSuffix: XDG_CACHE, get: GetCacheDir},
		{name: "data", xdgEnvVar: XDG_DATA_HOME, fallbackSuffix: XDG_DATA, get: GetDataDir},
		{name: "state", xdgEnvVar: XDG_STATE_HOME, fallbackSuffix: XDG_STATE, get: GetStateDir},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("uses explicit xdg path", func(t *testing.T) {
				clearXDGEnv(t)

				expected := filepath.Join(t.TempDir(), "xdg")
				t.Setenv(tc.xdgEnvVar, expected)

				got, err := tc.get()
				if err != nil {
					t.Fatalf("%s returned error: %v", tc.xdgEnvVar, err)
				}
				if got != expected {
					t.Fatalf("%s returned %q, want %q", tc.xdgEnvVar, got, expected)
				}
			})

			t.Run("uses HOME fallback", func(t *testing.T) {
				clearXDGEnv(t)

				home := t.TempDir()
				expected := filepath.Join(home, tc.fallbackSuffix[1:])
				if err := os.MkdirAll(expected, 0o755); err != nil {
					t.Fatalf("MkdirAll(%q): %v", expected, err)
				}
				t.Setenv(XDG_HOME, home)

				got, err := tc.get()
				if err != nil {
					t.Fatalf("%s fallback returned error: %v", tc.xdgEnvVar, err)
				}
				if got != expected {
					t.Fatalf("%s fallback returned %q, want %q", tc.xdgEnvVar, got, expected)
				}
			})

			t.Run("errors when HOME fallback missing", func(t *testing.T) {
				clearXDGEnv(t)
				t.Setenv(XDG_HOME, t.TempDir())

				got, err := tc.get()
				if err == nil {
					t.Fatalf("%s fallback returned path %q, want error", tc.xdgEnvVar, got)
				}
			})

			t.Run("errors when HOME fallback is a file", func(t *testing.T) {
				clearXDGEnv(t)

				home := t.TempDir()
				path := filepath.Join(home, tc.fallbackSuffix[1:])
				if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
					t.Fatalf("MkdirAll(%q): %v", filepath.Dir(path), err)
				}
				if err := os.WriteFile(path, []byte("not a directory"), 0o600); err != nil {
					t.Fatalf("WriteFile(%q): %v", path, err)
				}
				t.Setenv(XDG_HOME, home)

				got, err := tc.get()
				if err == nil {
					t.Fatalf("%s fallback returned path %q, want error", tc.xdgEnvVar, got)
				}
				if got != "" {
					t.Fatalf("%s fallback returned %q, want empty path on error", tc.xdgEnvVar, got)
				}
			})
		})
	}
}
