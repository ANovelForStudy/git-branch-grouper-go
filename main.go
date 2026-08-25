package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/pelletier/go-toml/v2"
)

func main() {
	var includeLong, includeShort string
	var excludeLong, excludeShort string
	var sparseLong, sparseShort bool

	flag.StringVar(&includeLong, "include", "", "Comma-separated list of groups or sub-paths to show (e.g. feat,backup/v1)")
	flag.StringVar(&includeShort, "i", "", "Shorthand for --include")
	flag.StringVar(&excludeLong, "exclude", "", "Comma-separated list of groups or sub-paths to hide (e.g. old,backup/v2)")
	flag.StringVar(&excludeShort, "e", "", "Shorthand for --exclude")
	flag.BoolVar(&sparseLong, "sparse", false, "Override config.toml: add blank line between groups")
	flag.BoolVar(&sparseShort, "s", false, "Shorthand for --sparse")
	flag.Usage = printUsage
	flag.Parse()

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
		log.Fatalf("Failed to get working directory: %v", err)
	}

	validateRepoPath(repoPath)

	cfg := loadConfig()

	if sparseVal {
		cfg.Main.Sparse = true
	}

	includeList := parsePrefixList(includeVal)
	excludeList := parsePrefixList(excludeVal)

	gitRepo := openGitRepo(repoPath)
	branchData := collectBranches(gitRepo)

	branchData = filterBranchData(branchData, includeList, excludeList)

	printResultsWithConfig(branchData, cfg, len(includeList) > 0 || len(excludeList) > 0)
}

func loadConfig() Config {
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

func printUsage() {
	bold := color.New(color.Bold)
	green := color.New(color.FgGreen)
	cyan := color.New(color.FgCyan)
	white := color.New(color.FgWhite)
	yellow := color.New(color.FgYellow)

	bold.Println("git-branch-grouper-plugin")
	fmt.Println("  Display and filter git branches grouped by prefix.")
	fmt.Println()
	bold.Println("USAGE:")
	white.Println("  git-branch-grouper-plugin [flags]")
	fmt.Println()
	bold.Println("FLAGS:")
	fmt.Printf("  %-12s  %s\n", bold.Sprint("-i, --include"), "Show only specified groups or sub-paths")
	fmt.Printf("  %-12s  %s\n", bold.Sprint("-e, --exclude"), "Hide specified groups or sub-paths")
	fmt.Printf("  %-12s  %s\n", bold.Sprint("-s, --sparse"), "Add blank line between groups")
	fmt.Printf("  %-12s  %s\n", bold.Sprint("-h, --help"), "Show this help message")
	fmt.Println()
	bold.Println("EXAMPLES:")
	green.Println("  git-branch-grouper-plugin --include feat,fix")
	yellow.Println("    Show only 'feat' and 'fix' groups.")
	fmt.Println()
	green.Println("  git-branch-grouper-plugin -e backup/v2")
	yellow.Println("    Hide 'v2' subtree inside 'backup' group.")
	fmt.Println()
	green.Println("  git-branch-grouper-plugin -i backup/v1,feat")
	yellow.Println("    Show 'feat' group and 'v1' subtree inside 'backup'.")
	fmt.Println()
	green.Println("  git-branch-grouper-plugin -s")
	yellow.Println("    Display with blank lines between groups.")
	fmt.Println()
	bold.Println("GROUPS:")
	cyan.Println("  [default]  Standard branches (main, master, develop)")
	cyan.Println("  [other]    Branches without a group prefix")
	cyan.Println("  prefix/    Hierarchical groups (feat/, fix/, etc.)")
	cyan.Println("  sub/path   Sub-paths within groups (backup/v2, feat/login)")
	fmt.Println()
	bold.Println("CONFIGURATION:")
	white.Println("  Colors and groups can be customized in config.toml")
	white.Println("  placed in the working directory.")
}
