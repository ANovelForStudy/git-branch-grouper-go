package config

import (
	"log"
	"os"

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

func Load() Config {
	cfg := Config{
		Colors: map[string]string{
			"default":     "red",
			"other":       "yellow",
			"active_star": "green",
		},
		Groups: map[string]string{},
	}

	file, err := os.ReadFile("config.toml")
	if err != nil {
		return cfg
	}

	if err := toml.Unmarshal(file, &cfg); err != nil {
		log.Printf("Failed to parse config.toml: %v", err)
	}

	return cfg
}
