package main

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

type BranchData struct {
	MainGroups      map[string][]string
	DefaultBranches []string
	OtherBranches   []string
}

func validateRepoPath(path string) {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Fatalf("Directory doesn't exist: %s", path)
		}
		log.Fatalf("Failed to open directory: %v", err)
	}
}

func openGitRepo(path string) *git.Repository {
	repo, err := git.PlainOpen(path)
	if err != nil {
		log.Fatalf("Failed to open repository: %v", err)
	}
	return repo
}

func collectBranches(repo *git.Repository) BranchData {
	branches, err := repo.Branches()
	if err != nil {
		log.Fatalf("Failed to get branches: %v", err)
	}

	activeBranch, err := repo.Head()
	if err != nil {
		log.Fatalf("Failed to get an active branch: %v", err)
	}
	activeBranchName := activeBranch.Name().Short()

	data := BranchData{
		MainGroups:      make(map[string][]string),
		DefaultBranches: []string{},
		OtherBranches:   []string{},
	}

	err = branches.ForEach(func(branch *plumbing.Reference) error {
		branchName := branch.Name().Short()
		isActive := branchName == activeBranchName

		splitted := strings.Split(branchName, "/")
		prefix := splitted[0]

		if len(splitted) > 1 {
			postfix := strings.Join(splitted[1:], "/")
			if isActive {
				postfix = "* " + postfix
			}
			data.MainGroups[prefix] = append(data.MainGroups[prefix], postfix)
		} else {
			if slices.Contains(defaultBranchNames, prefix) {
				if isActive {
					prefix = "* " + prefix
				}
				data.DefaultBranches = append(data.DefaultBranches, prefix)
			} else {
				if isActive {
					prefix = "* " + prefix
				}
				data.OtherBranches = append(data.OtherBranches, prefix)
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error iterating branches: %v", err)
	}

	return data
}
