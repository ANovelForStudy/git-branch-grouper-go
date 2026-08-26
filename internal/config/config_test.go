package config

import (
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
			Separator:    "\n",
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
	if cfg.Format.Separator != "\n" {
		t.Errorf("Separator = %q, want %q", cfg.Format.Separator, "\n")
	}
}

func TestTOMLParseFullFormat(t *testing.T) {
	input := []byte(`
[format]
group_prefix = "({group})"
indent = "\t"
branch_marker = ">"
separator = "\n\n"
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
	if cfg.Format.Separator != "\n\n" {
		t.Errorf("Separator = %q, want %q", cfg.Format.Separator, "\n\n")
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
	if cfg.Format.Separator != "" {
		t.Errorf("Separator should be empty (zero value), got %q", cfg.Format.Separator)
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
	if cfg.Format.Separator != "" {
		t.Errorf("Separator should be empty, got %q", cfg.Format.Separator)
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
			Separator:    "\n",
		},
	}

	input := []byte(`
[format]
branch_marker = "•"
separator = ""
`)
	if err := toml.Unmarshal(input, &cfg); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if cfg.Format.BranchMarker != "•" {
		t.Errorf("BranchMarker = %q, want %q", cfg.Format.BranchMarker, "•")
	}
	if cfg.Format.Separator != "" {
		t.Errorf("Separator = %q, want empty", cfg.Format.Separator)
	}
	if cfg.Format.GroupPrefix != "[{group}]" {
		t.Errorf("GroupPrefix should retain default %q, got %q", "[{group}]", cfg.Format.GroupPrefix)
	}
	if cfg.Format.Indent != "    " {
		t.Errorf("Indent should retain default %q, got %q", "    ", cfg.Format.Indent)
	}
}
