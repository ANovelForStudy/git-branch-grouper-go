package filter

import (
	"strings"

	"git-branch-grouper-plugin/internal/model"
)

func ParsePrefixList(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func Apply(data model.BranchData, includeList, excludeList []string) model.BranchData {
	if len(includeList) == 0 && len(excludeList) == 0 {
		return data
	}

	includeTree := buildFilterTree(includeList)
	excludeTree := buildFilterTree(excludeList)

	filtered := model.BranchData{
		MainGroups: make(map[string]*model.Node),
	}

	for groupName, srcNode := range data.MainGroups {
		inInclude := includeTree[groupName] != nil
		inExcludeFull := excludeTree[groupName] != nil && len(excludeTree[groupName]) == 0

		if len(includeList) > 0 && !inInclude {
			continue
		}

		if inExcludeFull && !inInclude {
			continue
		}

		hasSubInclude := inInclude && len(includeTree[groupName]) > 0
		hasSubExclude := excludeTree[groupName] != nil && len(excludeTree[groupName]) > 0

		if hasSubInclude {
			result := model.NewNode(srcNode.Name)
			for _, subFilter := range includeTree[groupName] {
				addIncludedPath(result, srcNode, subFilter)
			}
			if len(result.Children) > 0 {
				filtered.MainGroups[groupName] = result
			}
		} else if hasSubExclude {
			result := copyFilteredSubtree(srcNode, excludeTree[groupName])
			if len(result.Children) > 0 {
				filtered.MainGroups[groupName] = result
			}
		} else {
			result := model.NewNode(srcNode.Name)
			result.IsActive = srcNode.IsActive
			copySubtree(result, srcNode)
			filtered.MainGroups[groupName] = result
		}
	}

	return filtered
}

func buildFilterTree(filterList []string) map[string][]string {
	tree := make(map[string][]string)
	for _, f := range filterList {
		parts := strings.SplitN(f, "/", 2)
		groupName := parts[0]
		if len(parts) > 1 {
			tree[groupName] = append(tree[groupName], parts[1])
		} else {
			tree[groupName] = []string{}
		}
	}
	return tree
}

func addIncludedPath(dst *model.Node, src *model.Node, path string) {
	parts := strings.SplitN(path, "/", 2)
	key := parts[0]

	child, exists := src.Children[key]
	if !exists {
		return
	}

	if existing, ok := dst.Children[key]; ok {
		if len(parts) > 1 {
			addIncludedPath(existing, child, parts[1])
		}
		return
	}

	newChild := model.NewNode(child.Name)
	newChild.IsActive = child.IsActive
	if len(parts) > 1 {
		addIncludedPath(newChild, child, parts[1])
	} else {
		copySubtree(newChild, child)
	}
	dst.Children[key] = newChild
	dst.SubKeys = append(dst.SubKeys, key)
}

func copySubtree(dst *model.Node, src *model.Node) {
	for _, subKey := range src.SubKeys {
		if _, exists := dst.Children[subKey]; exists {
			continue
		}
		srcChild := src.Children[subKey]
		dstChild := model.NewNode(srcChild.Name)
		dstChild.IsActive = srcChild.IsActive
		dst.Children[subKey] = dstChild
		dst.SubKeys = append(dst.SubKeys, subKey)
		copySubtree(dstChild, srcChild)
	}
}

func copyFilteredSubtree(src *model.Node, excludePaths []string) *model.Node {
	result := model.NewNode(src.Name)

	for _, subKey := range src.SubKeys {
		child := src.Children[subKey]
		excludeAll := false
		for _, ep := range excludePaths {
			parts := strings.SplitN(ep, "/", 2)
			if parts[0] == subKey {
				if len(parts) == 1 {
					excludeAll = true
					break
				}
				subResult := copyFilteredSubtree(child, []string{parts[1]})
				result.Children[subKey] = subResult
				result.SubKeys = append(result.SubKeys, subKey)
			}
		}
		if excludeAll {
			continue
		}
		if _, exists := result.Children[subKey]; !exists {
			newChild := model.NewNode(child.Name)
			newChild.IsActive = child.IsActive
			copySubtree(newChild, child)
			result.Children[subKey] = newChild
			result.SubKeys = append(result.SubKeys, subKey)
		}
	}

	return result
}

func FindMatchingChild(node *model.Node, subPath string) *model.Node {
	parts := strings.Split(subPath, "/")
	curr := node
	for _, part := range parts {
		child, exists := curr.Children[part]
		if !exists {
			return nil
		}
		curr = child
	}
	return curr
}
