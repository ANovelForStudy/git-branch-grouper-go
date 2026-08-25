package git

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"
	"strings"

	gogit "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"

	"git-branch-grouper-plugin/internal/model"
)

var defaultBranchNames = []string{"develop", "main", "master"}

func ValidateRepoPath(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("directory does not exist: %s", path)
	}
	return fmt.Errorf("cannot access directory %s: %w", path, err)
}

func OpenRepo(path string) (*gogit.Repository, error) {
	repo, err := gogit.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("not a git repository (or any parent up to root): %s", path)
	}
	return repo, nil
}

func CollectBranches(repo *gogit.Repository) (model.BranchData, error) {
	branches, err := repo.Branches()
	if err != nil {
		return model.BranchData{}, fmt.Errorf("failed to list branches: %w", err)
	}

	activeBranch, err := repo.Head()
	if err != nil {
		return model.BranchData{}, fmt.Errorf("failed to resolve HEAD: %w", err)
	}
	activeBranchName := activeBranch.Name().Short()

	data := model.BranchData{
		MainGroups:      make(map[string]*model.Node),
		DefaultBranches: []string{},
		OtherBranches:   []string{},
	}

	err = branches.ForEach(func(branch *plumbing.Reference) error {
		branchName := branch.Name().Short()
		isActive := branchName == activeBranchName

		splitted := strings.Split(branchName, "/")
		prefix := splitted[0]

		if len(splitted) > 1 {
			rootNode, exists := data.MainGroups[prefix]
			if !exists {
				rootNode = model.NewNode(prefix)
				data.MainGroups[prefix] = rootNode
			}

			curr := rootNode
			for i := 1; i < len(splitted); i++ {
				part := splitted[i]
				isLast := i == len(splitted)-1

				child, childExists := curr.Children[part]
				if !childExists {
					child = model.NewNode(part)
					child.IsActive = isLast && isActive
					curr.Children[part] = child
					curr.SubKeys = append(curr.SubKeys, part)
				} else if isLast && isActive {
					child.IsActive = true
				}
				curr = child
			}
		} else {
			targetList := &data.OtherBranches
			cleanPrefix := strings.TrimPrefix(prefix, "* ")
			if slices.Contains(defaultBranchNames, cleanPrefix) {
				targetList = &data.DefaultBranches
			}

			if isActive {
				prefix = "* " + prefix
			}
			*targetList = append(*targetList, prefix)
		}

		return nil
	})

	if err != nil {
		return model.BranchData{}, fmt.Errorf("failed to iterate branches: %w", err)
	}

	return data, nil
}
