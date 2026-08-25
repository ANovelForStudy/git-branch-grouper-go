package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

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

	execPath, err := os.Executable()
	if err != nil {
		return cfg, ""
	}

	configPath := filepath.Join(filepath.Dir(execPath), "config.toml")

	file, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, configPath
	}

	if err := toml.Unmarshal(file, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to parse %s: %v\n", configPath, err)
	}

	return cfg, configPath
}
