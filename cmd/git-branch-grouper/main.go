package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"

	"git-branch-grouper-plugin/internal/config"
	"git-branch-grouper-plugin/internal/display"
	"git-branch-grouper-plugin/internal/filter"
	"git-branch-grouper-plugin/internal/git"
)

var version = "dev"

func main() {
	var includeLong, includeShort string
	var excludeLong, excludeShort string
	var sparseLong, sparseShort bool
	var noColorLong, noColorShort bool
	var showVersion bool
	var configPath string

	flag.StringVar(&includeLong, "include", "", "Show only specified groups or sub-paths")
	flag.StringVar(&includeShort, "i", "", "Shorthand for --include")
	flag.StringVar(&excludeLong, "exclude", "", "Hide specified groups or sub-paths")
	flag.StringVar(&excludeShort, "e", "", "Shorthand for --exclude")
	flag.BoolVar(&sparseLong, "sparse", false, "Add blank line between groups")
	flag.BoolVar(&sparseShort, "s", false, "Shorthand for --sparse")
	flag.BoolVar(&noColorLong, "no-color", false, "Disable colored output")
	flag.BoolVar(&noColorShort, "n", false, "Shorthand for --no-color")
	flag.BoolVar(&showVersion, "version", false, "Print version and exit")
	flag.BoolVar(&showVersion, "v", false, "Shorthand for --version")
	flag.StringVar(&configPath, "config", "", "Path to config file")
	flag.Usage = printUsage
	flag.Parse()

	if showVersion {
		fmt.Printf("git-branch-grouper v%s\n", version)
		os.Exit(0)
	}

	mergeFlags := func(long, short string) string {
		if short != "" {
			return short
		}
		return long
	}

	includeVal := mergeFlags(includeLong, includeShort)
	excludeVal := mergeFlags(excludeLong, excludeShort)
	sparseVal := sparseLong || sparseShort
	noColor := noColorLong || noColorShort || os.Getenv("NO_COLOR") != ""

	if noColor {
		color.NoColor = true
	}

	repoPath, err := os.Getwd()
	if err != nil {
		fatalf("cannot determine working directory: %v", err)
	}

	if err := git.ValidateRepoPath(repoPath); err != nil {
		fatalf("%v", err)
	}

	cfg, _ := config.Load(configPath)

	if sparseVal {
		cfg.Main.Sparse = true
	}

	includeList := filter.ParsePrefixList(includeVal)
	excludeList := filter.ParsePrefixList(excludeVal)

	repo, err := git.OpenRepo(repoPath)
	if err != nil {
		fatalf("%v", err)
	}

	data, err := git.CollectBranches(repo)
	if err != nil {
		fatalf("%v", err)
	}

	data = filter.Apply(data, includeList, excludeList)

	display.PrintResults(os.Stdout, data, cfg, len(includeList) > 0 || len(excludeList) > 0)
}

func fatalf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	red := color.New(color.FgRed, color.Bold)
	_, _ = red.Fprint(os.Stderr, "error: ")
	_, _ = color.New(color.FgWhite).Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func printUsage() {
	bold := color.New(color.Bold)
	dim := color.New(color.FgHiBlack)
	cyan := color.New(color.FgCyan)
	green := color.New(color.FgGreen)
	yellow := color.New(color.FgYellow)
	white := color.New(color.FgWhite)

	name := color.New(color.FgCyan, color.Bold)

	fmt.Println()
	_, _ = name.Print("  git-branch-grouper")
	_, _ = dim.Printf("  v%s", version)
	fmt.Println()
	_, _ = white.Println("  Organize and navigate git branches with prefix-based grouping.")
	fmt.Println()

	_, _ = bold.Println("  USAGE")
	_, _ = white.Println("    git-branch-grouper [flags]")
	fmt.Println()

	_, _ = bold.Println("  FLAGS")
	printFlag("-i, --include", "Show only specified groups or sub-paths")
	printFlag("-e, --exclude", "Hide specified groups or sub-paths")
	printFlag("-s, --sparse", "Add blank line between groups")
	printFlag("-n, --no-color", "Disable colored output")
	printFlag("--config", "Path to config file")
	printFlag("-v, --version", "Print version and exit")
	printFlag("-h, --help", "Show this help message")
	fmt.Println()

	_, _ = bold.Println("  EXAMPLES")
	_, _ = green.Println("    git-branch-grouper --include feat,fix")
	_, _ = yellow.Println("      Show only 'feat' and 'fix' groups.")
	fmt.Println()
	_, _ = green.Println("    git-branch-grouper -e backup/v2")
	_, _ = yellow.Println("      Hide 'v2' subtree inside 'backup' group.")
	fmt.Println()
	_, _ = green.Println("    git-branch-grouper -i backup/v1,feat")
	_, _ = yellow.Println("      Show 'feat' group and 'v1' subtree inside 'backup'.")
	fmt.Println()
	_, _ = green.Println("    git-branch-grouper -s")
	_, _ = yellow.Println("      Display with blank lines between groups.")
	fmt.Println()

	_, _ = bold.Println("  BRANCH GROUPS")
	_, _ = cyan.Println("    [default]   Standard branches (main, master, develop)")
	_, _ = cyan.Println("    [other]     Branches without a group prefix")
	_, _ = cyan.Println("    prefix/     Hierarchical groups (feat/, fix/, etc.)")
	_, _ = cyan.Println("    sub/path    Sub-paths within groups (backup/v2, feat/login)")
	fmt.Println()

	_, _ = bold.Println("  CONFIGURATION")
	_, _ = white.Println("    Use --config <path> to specify a custom config file.")
	_, _ = white.Println("    Without --config, the first found path is used:")
	_, _ = dim.Println("      1. ./config.toml")
	_, _ = dim.Println("      2. $XDG_CONFIG_HOME/git-branch-grouper/config.toml")
	_, _ = dim.Println("      3. ~/.config/git-branch-grouper/config.toml")
	_, _ = dim.Println("      4. <binary-dir>/config.toml")
	_, _ = white.Println("    See README.md for the full configuration reference.")
	fmt.Println()
}

func printFlag(flags, description string) {
	bold := color.New(color.Bold)
	dim := color.New(color.FgHiBlack)
	fmt.Printf("    %-18s", bold.Sprint(flags))
	_, _ = dim.Println(description)
}
