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
	var showVersion bool

	flag.StringVar(&includeLong, "include", "", "Show only specified groups or sub-paths")
	flag.StringVar(&includeShort, "i", "", "Shorthand for --include")
	flag.StringVar(&excludeLong, "exclude", "", "Hide specified groups or sub-paths")
	flag.StringVar(&excludeShort, "e", "", "Shorthand for --exclude")
	flag.BoolVar(&sparseLong, "sparse", false, "Add blank line between groups")
	flag.BoolVar(&sparseShort, "s", false, "Shorthand for --sparse")
	flag.BoolVar(&showVersion, "version", false, "Print version and exit")
	flag.BoolVar(&showVersion, "v", false, "Shorthand for --version")
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

	repoPath, err := os.Getwd()
	if err != nil {
		fatalf("cannot determine working directory: %v", err)
	}

	if err := git.ValidateRepoPath(repoPath); err != nil {
		fatalf("%v", err)
	}

	cfg, _ := config.Load()

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

	display.PrintResults(data, cfg, len(includeList) > 0 || len(excludeList) > 0)
}

func fatalf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	red := color.New(color.FgRed, color.Bold)
	red.Fprint(os.Stderr, "error: ")
	color.New(color.FgWhite).Fprintln(os.Stderr, msg)
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
	name.Print("  git-branch-grouper")
	dim.Printf("  v%s", version)
	fmt.Println()
	white.Println("  Organize and navigate git branches with prefix-based grouping.")
	fmt.Println()

	bold.Println("  USAGE")
	white.Println("    git-branch-grouper [flags]")
	fmt.Println()

	bold.Println("  FLAGS")
	printFlag("-i, --include", "Show only specified groups or sub-paths")
	printFlag("-e, --exclude", "Hide specified groups or sub-paths")
	printFlag("-s, --sparse", "Add blank line between groups")
	printFlag("-v, --version", "Print version and exit")
	printFlag("-h, --help", "Show this help message")
	fmt.Println()

	bold.Println("  EXAMPLES")
	green.Println("    git-branch-grouper --include feat,fix")
	yellow.Println("      Show only 'feat' and 'fix' groups.")
	fmt.Println()
	green.Println("    git-branch-grouper -e backup/v2")
	yellow.Println("      Hide 'v2' subtree inside 'backup' group.")
	fmt.Println()
	green.Println("    git-branch-grouper -i backup/v1,feat")
	yellow.Println("      Show 'feat' group and 'v1' subtree inside 'backup'.")
	fmt.Println()
	green.Println("    git-branch-grouper -s")
	yellow.Println("      Display with blank lines between groups.")
	fmt.Println()

	bold.Println("  BRANCH GROUPS")
	cyan.Println("    [default]   Standard branches (main, master, develop)")
	cyan.Println("    [other]     Branches without a group prefix")
	cyan.Println("    prefix/     Hierarchical groups (feat/, fix/, etc.)")
	cyan.Println("    sub/path    Sub-paths within groups (backup/v2, feat/login)")
	fmt.Println()

	bold.Println("  CONFIGURATION")
	white.Println("    Place config.toml next to the binary to customize colors and groups.")
	white.Println("    See README.md for the full configuration reference.")
	fmt.Println()
}

func printFlag(flags, description string) {
	bold := color.New(color.Bold)
	dim := color.New(color.FgHiBlack)
	fmt.Printf("    %-18s", bold.Sprint(flags))
	dim.Println(description)
}
