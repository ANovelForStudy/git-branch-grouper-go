package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

const appName = "git-branch-grouper"

type Main struct {
	Sparse bool `toml:"sparse"`
}

type Config struct {
	Main   Main              `toml:"main"`
	Colors map[string]string `toml:"colors"`
	Groups map[string]string `toml:"groups"`
}

func Load() (Config, string) {
	cfg := Config{
		Colors: map[string]string{
			"default":     "red",
			"other":       "yellow",
			"active_star": "green",
		},
		Groups: map[string]string{},
	}

	configPath := resolveConfigPath()

	file, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, configPath
	}

	if err := toml.Unmarshal(file, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to parse %s: %v\n", configPath, err)
	}

	return cfg, configPath
}

func resolveConfigPath() string {
	if p := filepath.Join(".", "config.toml"); fileExists(p) {
		return p
	}

	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		if p := filepath.Join(xdg, appName, "config.toml"); fileExists(p) {
			return p
		}
	}

	if home, err := os.UserHomeDir(); err == nil {
		if p := filepath.Join(home, ".config", appName, "config.toml"); fileExists(p) {
			return p
		}
	}

	if execPath, err := os.Executable(); err == nil {
		return filepath.Join(filepath.Dir(execPath), "config.toml")
	}

	return filepath.Join(".", "config.toml")
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
