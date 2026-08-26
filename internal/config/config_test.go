package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func TestDefaultFormatValues(t *testing.T) {
	cfg := Config{
		Colors: map[string]string{},
		Groups: map[string]string{},
		Format: Format{
			GroupPrefix:  "[{group}]",
			Indent:       "    ",
			BranchMarker: "*",
		},
	}

	if cfg.Format.GroupPrefix != "[{group}]" {
		t.Errorf("GroupPrefix = %q, want %q", cfg.Format.GroupPrefix, "[{group}]")
	}
	if cfg.Format.Indent != "    " {
		t.Errorf("Indent = %q, want %q", cfg.Format.Indent, "    ")
	}
	if cfg.Format.BranchMarker != "*" {
		t.Errorf("BranchMarker = %q, want %q", cfg.Format.BranchMarker, "*")
	}
}

func TestTOMLParseFullFormat(t *testing.T) {
	input := []byte(`
[format]
group_prefix = "({group})"
indent = "\t"
branch_marker = ">"
`)
	var cfg Config
	if err := toml.Unmarshal(input, &cfg); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if cfg.Format.GroupPrefix != "({group})" {
		t.Errorf("GroupPrefix = %q, want %q", cfg.Format.GroupPrefix, "({group})")
	}
	if cfg.Format.Indent != "\t" {
		t.Errorf("Indent = %q, want %q", cfg.Format.Indent, "\t")
	}
	if cfg.Format.BranchMarker != ">" {
		t.Errorf("BranchMarker = %q, want %q", cfg.Format.BranchMarker, ">")
	}
}

func TestTOMLParsePartialFormat(t *testing.T) {
	input := []byte(`
[format]
branch_marker = "+"
`)
	var cfg Config
	if err := toml.Unmarshal(input, &cfg); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if cfg.Format.BranchMarker != "+" {
		t.Errorf("BranchMarker = %q, want %q", cfg.Format.BranchMarker, "+")
	}
	if cfg.Format.GroupPrefix != "" {
		t.Errorf("GroupPrefix should be empty (zero value), got %q", cfg.Format.GroupPrefix)
	}
	if cfg.Format.Indent != "" {
		t.Errorf("Indent should be empty (zero value), got %q", cfg.Format.Indent)
	}
}

func TestTOMLParseEmptyFormat(t *testing.T) {
	input := []byte(`
[main]
sparse = true
`)
	var cfg Config
	if err := toml.Unmarshal(input, &cfg); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if cfg.Format.GroupPrefix != "" {
		t.Errorf("GroupPrefix should be empty, got %q", cfg.Format.GroupPrefix)
	}
	if cfg.Format.Indent != "" {
		t.Errorf("Indent should be empty, got %q", cfg.Format.Indent)
	}
	if cfg.Format.BranchMarker != "" {
		t.Errorf("BranchMarker should be empty, got %q", cfg.Format.BranchMarker)
	}
}

func TestLoadMergesDefaultsWithFile(t *testing.T) {
	cfg := Config{
		Colors: map[string]string{},
		Groups: map[string]string{},
		Format: Format{
			GroupPrefix:  "[{group}]",
			Indent:       "    ",
			BranchMarker: "*",
		},
	}

	input := []byte(`
[format]
branch_marker = "•"
`)
	if err := toml.Unmarshal(input, &cfg); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if cfg.Format.BranchMarker != "•" {
		t.Errorf("BranchMarker = %q, want %q", cfg.Format.BranchMarker, "•")
	}
	if cfg.Format.GroupPrefix != "[{group}]" {
		t.Errorf("GroupPrefix should retain default %q, got %q", "[{group}]", cfg.Format.GroupPrefix)
	}
	if cfg.Format.Indent != "    " {
		t.Errorf("Indent should retain default %q, got %q", "    ", cfg.Format.Indent)
	}
}

func TestLoadWithExplicitPath(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "custom.toml")
	content := []byte(`
[format]
branch_marker = "+"
`)
	if err := os.WriteFile(cfgPath, content, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, returned := Load(cfgPath)

	if returned != cfgPath {
		t.Errorf("returned path = %q, want %q", returned, cfgPath)
	}
	if cfg.Format.BranchMarker != "+" {
		t.Errorf("BranchMarker = %q, want %q", cfg.Format.BranchMarker, "+")
	}
}

func TestLoadWithExplicitPathNotFound(t *testing.T) {
	cfg, returned := Load("/nonexistent/config.toml")

	if returned != "/nonexistent/config.toml" {
		t.Errorf("returned path = %q, want /nonexistent/config.toml", returned)
	}
	if cfg.Format.BranchMarker != "*" {
		t.Errorf("BranchMarker should be default '*', got %q", cfg.Format.BranchMarker)
	}
}

func TestLoadWithEmptyExplicitPath(t *testing.T) {
	cfg, _ := Load("")

	if cfg.Format.BranchMarker != "*" {
		t.Errorf("BranchMarker should be default '*', got %q", cfg.Format.BranchMarker)
	}
}
