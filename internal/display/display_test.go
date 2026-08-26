package display

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/fatih/color"

	"git-branch-grouper-plugin/internal/config"
	"git-branch-grouper-plugin/internal/model"
)

func makeNode(name string, isActive bool, children ...*model.Node) *model.Node {
	n := model.NewNode(name)
	n.IsActive = isActive
	for _, c := range children {
		n.Children[c.Name] = c
		n.SubKeys = append(n.SubKeys, c.Name)
	}
	return n
}

func defaultConfig() config.Config {
	return config.Config{
		Colors: map[string]string{
			"default":     "red",
			"other":       "yellow",
			"active_star": "green",
		},
		Groups: map[string]string{},
		Format: config.Format{
			GroupPrefix:  "[{group}]",
			Indent:       "    ",
			BranchMarker: "*",
		},
	}
}

func capture(fn func(io.Writer)) string {
	var buf bytes.Buffer
	fn(&buf)
	return buf.String()
}

func TestGetColorPrinterNamedColors(t *testing.T) {
	namedColors := []string{
		"black", "red", "green", "yellow", "blue", "magenta", "cyan", "white",
		"hi-black", "hi-red", "hi-green", "hi-yellow", "hi-blue", "hi-magenta", "hi-cyan", "hi-white",
	}
	for _, name := range namedColors {
		c := getColorPrinter(name)
		if c == nil {
			t.Errorf("getColorPrinter(%q) returned nil", name)
		}
	}
}

func TestGetColorPrinterHex6(t *testing.T) {
	c := getColorPrinter("#FF8800")
	if c == nil {
		t.Fatal("getColorPrinter(#FF8800) returned nil")
	}
}

func TestGetColorPrinterHex3(t *testing.T) {
	c := getColorPrinter("#F00")
	if c == nil {
		t.Fatal("getColorPrinter(#F00) returned nil")
	}
}

func TestGetColorPrinterHexInvalid(t *testing.T) {
	c := getColorPrinter("#ZZZZZZ")
	if c == nil {
		t.Fatal("getColorPrinter(#ZZZZZZ) returned nil")
	}
}

func TestGetColorPrinter256Color(t *testing.T) {
	c := getColorPrinter("196")
	if c == nil {
		t.Fatal("getColorPrinter(196) returned nil")
	}
}

func TestGetColorPrinter256ColorOutOfRange(t *testing.T) {
	c := getColorPrinter("300")
	if c == nil {
		t.Fatal("getColorPrinter(300) returned nil")
	}
}

func TestGetColorPrinterUnknownFallback(t *testing.T) {
	c := getColorPrinter("notacolor")
	if c == nil {
		t.Fatal("getColorPrinter(notacolor) returned nil")
	}
}

func TestPrintResultsWithHexColor(t *testing.T) {
	cfg := defaultConfig()
	cfg.Groups["feat"] = "#FF8800"

	data := model.BranchData{
		MainGroups: map[string]*model.Node{
			"feat": makeNode("feat", false, makeNode("auth", false)),
		},
		DefaultBranches: []string{},
		OtherBranches:   []string{},
	}

	output := capture(func(w io.Writer) {
		PrintResults(w, data, cfg, true)
	})

	if !strings.Contains(output, "feat/") {
		t.Errorf("expected feat/ in output, got:\n%s", output)
	}
}

func TestPrintResultsWith256Color(t *testing.T) {
	cfg := defaultConfig()
	cfg.Groups["fix"] = "196"

	data := model.BranchData{
		MainGroups: map[string]*model.Node{
			"fix": makeNode("fix", false, makeNode("bug", false)),
		},
		DefaultBranches: []string{},
		OtherBranches:   []string{},
	}

	output := capture(func(w io.Writer) {
		PrintResults(w, data, cfg, true)
	})

	if !strings.Contains(output, "fix/") {
		t.Errorf("expected fix/ in output, got:\n%s", output)
	}
}

func TestPrintResultsDefaultFormat(t *testing.T) {
	data := model.BranchData{
		MainGroups: map[string]*model.Node{
			"feat": makeNode("feat", false,
				makeNode("login", true),
				makeNode("signup", false),
			),
		},
		DefaultBranches: []string{"main", "* develop"},
		OtherBranches:   []string{"random-branch"},
	}

	output := capture(func(w io.Writer) {
		PrintResults(w, data, defaultConfig(), false)
	})

	if !strings.Contains(output, "[default]") {
		t.Error("expected [default] header")
	}
	if !strings.Contains(output, "feat/") {
		t.Error("expected feat/ header (no brackets)")
	}
	if !strings.Contains(output, "[other]") {
		t.Error("expected [other] header")
	}
	if !strings.Contains(output, "* develop") {
		t.Error("expected '* develop' with star marker")
	}
	if !strings.Contains(output, "* login") {
		t.Error("expected '* login' with star marker")
	}
}

func TestPrintResultsGroupPrefixOnlyForSpecialGroups(t *testing.T) {
	cfg := defaultConfig()
	cfg.Format.GroupPrefix = "({group})"

	data := model.BranchData{
		MainGroups: map[string]*model.Node{
			"fix": makeNode("fix", false, makeNode("bug", false)),
		},
		DefaultBranches: []string{"main"},
		OtherBranches:   []string{},
	}

	output := capture(func(w io.Writer) {
		PrintResults(w, data, cfg, false)
	})

	if !strings.Contains(output, "(default)") {
		t.Errorf("expected '(default)' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "(other)") {
		t.Errorf("expected '(other)' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "fix/") {
		t.Errorf("named group should be 'fix/' (no prefix), got:\n%s", output)
	}
	if strings.Contains(output, "(fix)") {
		t.Errorf("named group should NOT use group_prefix, got:\n%s", output)
	}
}

func TestPrintResultsCustomIndent(t *testing.T) {
	cfg := defaultConfig()
	cfg.Format.Indent = "\t"

	data := model.BranchData{
		MainGroups: map[string]*model.Node{
			"feat": makeNode("feat", false, makeNode("auth", false)),
		},
		DefaultBranches: []string{},
		OtherBranches:   []string{},
	}

	output := capture(func(w io.Writer) {
		PrintResults(w, data, cfg, true)
	})

	if !strings.Contains(output, "\tauth") {
		t.Errorf("expected tab-indented 'auth', got:\n%s", output)
	}
}

func TestPrintResultsCustomMarker(t *testing.T) {
	cfg := defaultConfig()
	cfg.Format.BranchMarker = ">"

	data := model.BranchData{
		MainGroups:      map[string]*model.Node{},
		DefaultBranches: []string{"* main"},
		OtherBranches:   []string{},
	}

	output := capture(func(w io.Writer) {
		PrintResults(w, data, cfg, false)
	})

	if !strings.Contains(output, "> main") {
		t.Errorf("expected '> main' with custom marker, got:\n%s", output)
	}
	if strings.Contains(output, "* main") {
		t.Errorf("should not contain default '*' marker, got:\n%s", output)
	}
}

func TestPrintResultsNoBlankLinesWithoutSparse(t *testing.T) {
	data := model.BranchData{
		MainGroups: map[string]*model.Node{
			"feat": makeNode("feat", false, makeNode("auth", false)),
			"fix":  makeNode("fix", false, makeNode("bug", false)),
		},
		DefaultBranches: []string{"main"},
		OtherBranches:   []string{"misc"},
	}

	output := capture(func(w io.Writer) {
		PrintResults(w, data, defaultConfig(), false)
	})

	output = strings.TrimRight(output, "\n")
	lines := strings.Split(output, "\n")
	blankCount := 0
	for _, line := range lines {
		if line == "" {
			blankCount++
		}
	}

	if blankCount != 0 {
		t.Errorf("expected 0 blank lines without sparse, got %d in:\n%s", blankCount, output)
	}
}

func TestPrintResultsBlankLinesWithSparse(t *testing.T) {
	cfg := defaultConfig()
	cfg.Main.Sparse = true

	data := model.BranchData{
		MainGroups: map[string]*model.Node{
			"feat": makeNode("feat", false, makeNode("auth", false)),
			"fix":  makeNode("fix", false, makeNode("bug", false)),
		},
		DefaultBranches: []string{"main"},
		OtherBranches:   []string{"misc"},
	}

	output := capture(func(w io.Writer) {
		PrintResults(w, data, cfg, false)
	})

	output = strings.TrimRight(output, "\n")
	lines := strings.Split(output, "\n")
	blankCount := 0
	for _, line := range lines {
		if line == "" {
			blankCount++
		}
	}

	if blankCount < 3 {
		t.Errorf("expected at least 3 blank lines with sparse (between 4 sections), got %d in:\n%s", blankCount, output)
	}
}

func TestPrintResultsWithFilter(t *testing.T) {
	data := model.BranchData{
		MainGroups: map[string]*model.Node{
			"feat": makeNode("feat", false, makeNode("auth", false)),
		},
		DefaultBranches: []string{"main"},
		OtherBranches:   []string{"misc"},
	}

	output := capture(func(w io.Writer) {
		PrintResults(w, data, defaultConfig(), true)
	})

	if strings.Contains(output, "[default]") {
		t.Error("should not contain [default] when filter is active")
	}
	if strings.Contains(output, "[other]") {
		t.Error("should not contain [other] when filter is active")
	}
	if !strings.Contains(output, "feat/") {
		t.Error("should contain feat/ even with filter")
	}
}

func TestPrintResultsNoColorMode(t *testing.T) {
	old := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = old }()

	cfg := defaultConfig()
	cfg.Groups["feat"] = "red"

	data := model.BranchData{
		MainGroups: map[string]*model.Node{
			"feat": makeNode("feat", false, makeNode("auth", true)),
		},
		DefaultBranches: []string{"main"},
		OtherBranches:   []string{"misc"},
	}

	output := capture(func(w io.Writer) {
		PrintResults(w, data, cfg, false)
	})

	if strings.Contains(output, "\x1b[") {
		t.Errorf("output should not contain ANSI escape codes when NoColor is set, got:\n%s", output)
	}
	if !strings.Contains(output, "[default]") {
		t.Error("expected [default] header")
	}
	if !strings.Contains(output, "feat/") {
		t.Error("expected feat/ header")
	}
	if !strings.Contains(output, "* auth") {
		t.Error("expected '* auth' with star marker")
	}
}
