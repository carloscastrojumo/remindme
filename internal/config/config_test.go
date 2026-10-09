package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/adrg/xdg"
)

func TestConfigDir(t *testing.T) {
	tests := []struct {
		name          string
		xdgConfigHome func(home string) string
		legacyConfig  bool
		want          func(home string) string
	}{
		{
			name:          "XDG_CONFIG_HOME unset",
			xdgConfigHome: func(string) string { return "" },
			want:          func(home string) string { return filepath.Join(home, ".config", "remindme") },
		},
		{
			name:          "XDG_CONFIG_HOME set",
			xdgConfigHome: func(home string) string { return filepath.Join(home, "xdg") },
			want:          func(home string) string { return filepath.Join(home, "xdg", "remindme") },
		},
		{
			name:          "existing config in legacy dir wins",
			xdgConfigHome: func(home string) string { return filepath.Join(home, "xdg") },
			legacyConfig:  true,
			want:          func(home string) string { return filepath.Join(home, ".config", "remindme") },
		},
		{
			name:          "relative XDG_CONFIG_HOME is ignored",
			xdgConfigHome: func(string) string { return "relative/xdg" },
			want:          func(home string) string { return filepath.Join(home, ".config", "remindme") },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(xdg.Reload)
			home := t.TempDir()
			t.Setenv("HOME", home)
			t.Setenv("XDG_CONFIG_HOME", tt.xdgConfigHome(home))
			xdg.Reload()

			if tt.legacyConfig {
				legacy := filepath.Join(home, ".config", "remindme")
				if err := os.MkdirAll(legacy, 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(legacy, "config.yaml"), nil, 0600); err != nil {
					t.Fatal(err)
				}
			}

			if got, want := configDir(), tt.want(home); got != want {
				t.Errorf("configDir() = %q, want %q", got, want)
			}
		})
	}
}
