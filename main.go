package main

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

const repoPath = `G:\05-Projects\01-Programming\01-Pet-Projects\06-Git-branch-grouper-plugin\99-Test-repo`

var defaultBranchNames = []string{"develop", "main", "master"}

type Config struct {
	Main   map[string]string `toml:"main"`
	Colors map[string]string `toml:"colors"`
	Groups map[string]string `toml:"groups"`
}

type Node struct {
	Name     string
	IsActive bool
	Children map[string]*Node
	SubKeys  []string
}

type BranchData struct {
	MainGroups      map[string]*Node
	DefaultBranches []string
	OtherBranches   []string
}

func main() {
	validateRepoPath(repoPath)

	cfg := loadConfig()

	gitRepo := openGitRepo(repoPath)
	branchData := collectBranches(gitRepo)

	printResultsWithConfig(branchData, cfg)
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

	_ = toml.Unmarshal(file, &cfg)
	return cfg
}
