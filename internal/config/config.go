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

type Format struct {
	GroupPrefix  string `toml:"group_prefix"`
	Indent       string `toml:"indent"`
	BranchMarker string `toml:"branch_marker"`
}

type Config struct {
	Main   Main              `toml:"main"`
	Colors map[string]string `toml:"colors"`
	Groups map[string]string `toml:"groups"`
	Format Format            `toml:"format"`
}

func Load(explicitPath ...string) (Config, string) {
	cfg := Config{
		Colors: map[string]string{
			"default":     "#FF3333",
			"other":       "#FBBF24",
			"active_star": "#b7ff32",
		},
		Groups: map[string]string{
			"backup":   "#64748B",
			"default":  "#1E293B",
			"feat":     "#10B981",
			"refactor": "#34D399",
			"fix":      "#EF4444",
			"hotfix":   "#DC2626",
			"build":    "#3B82F6",
			"ci":       "#22D3EE",
			"release":  "#0EA5E9",
			"docs":     "#818CF8",
			"test":     "#C084FC",
			"chore":    "#9CA3AF",
			"old":      "#F59E0B",
			"exp":      "#F97316",
		},
		Format: Format{
			GroupPrefix:  "[{group}]",
			Indent:       "    ",
			BranchMarker: "*",
		},
	}

	var configPath string

	if len(explicitPath) > 0 && explicitPath[0] != "" {
		configPath = explicitPath[0]
	} else {
		configPath = resolveConfigPath()
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		if len(explicitPath) > 0 && explicitPath[0] != "" {
			fmt.Fprintf(os.Stderr, "Warning: specified config %s not found\n", configPath)
		}
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
