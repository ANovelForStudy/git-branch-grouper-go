package display

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"

	"git-branch-grouper-plugin/internal/config"
	"git-branch-grouper-plugin/internal/model"
)

func PrintResults(data model.BranchData, cfg config.Config, hasFilter bool) {
	starColorStr := cfg.Colors["active_star"]

	if !hasFilter {
		defaultColorStr := cfg.Colors["default"]
		printSimpleGroup("[default]", data.DefaultBranches, defaultColorStr, starColorStr, cfg.Main.Sparse)
	}

	groupNames := make([]string, 0, len(data.MainGroups))
	for name := range data.MainGroups {
		groupNames = append(groupNames, name)
	}
	sort.Strings(groupNames)

	for _, groupName := range groupNames {
		groupColorStr := cfg.Groups[groupName]
		if groupColorStr == "" {
			groupColorStr = "cyan"
		}

		title := fmt.Sprintf("%s/", groupName)
		titleColor := getColorPrinter(groupColorStr)
		_, _ = titleColor.Println(title)

		rootNode := data.MainGroups[groupName]
		sort.Strings(rootNode.SubKeys)
		for _, subKey := range rootNode.SubKeys {
			printNode(rootNode.Children[subKey], 1, groupColorStr, starColorStr, cfg)
		}

		if cfg.Main.Sparse {
			fmt.Println()
		}
	}

	if !hasFilter {
		otherColorStr := cfg.Colors["other"]
		printSimpleGroup("[other]", data.OtherBranches, otherColorStr, starColorStr, cfg.Main.Sparse)
	}
}

func printNode(node *model.Node, level int, parentColorStr string, starColorStr string, cfg config.Config) {
	tabulation := strings.Repeat("    ", level)
	starColor := getColorPrinter(starColorStr)

	nodeColorStr := cfg.Groups[node.Name]
	if nodeColorStr == "" {
		nodeColorStr = parentColorStr
	}
	nodeColor := getColorPrinter(nodeColorStr)

	if len(node.Children) == 0 {
		if node.IsActive {
			fmt.Print(tabulation)
			_, _ = starColor.Print("* ")
			_, _ = nodeColor.Println(node.Name)
		} else {
			_, _ = nodeColor.Println(tabulation + node.Name)
		}
		return
	}

	if node.IsActive {
		fmt.Print(tabulation)
		_, _ = starColor.Print("* ")
		_, _ = nodeColor.Println(node.Name + "/")
	} else {
		_, _ = nodeColor.Println(tabulation + node.Name + "/")
	}

	sort.Strings(node.SubKeys)
	for _, subKey := range node.SubKeys {
		printNode(node.Children[subKey], level+1, nodeColorStr, starColorStr, cfg)
	}
}

func printSimpleGroup(title string, branches []string, titleColorStr string, starColorStr string, sparse bool) {
	titleColor := getColorPrinter(titleColorStr)
	_, _ = titleColor.Println(title)

	starColor := getColorPrinter(starColorStr)
	tabulation := "    "

	sort.Strings(branches)
	for _, branch := range branches {
		if strings.HasPrefix(branch, "* ") {
			fmt.Print(tabulation)
			_, _ = starColor.Print("* ")
			fmt.Println(strings.TrimPrefix(branch, "* "))
		} else {
			fmt.Println(tabulation, branch)
		}
	}

	if sparse {
		fmt.Println()
	}
}

func getColorPrinter(colorName string) *color.Color {
	switch strings.ToLower(colorName) {
	case "red":
		return color.New(color.FgRed)
	case "green":
		return color.New(color.FgGreen)
	case "yellow":
		return color.New(color.FgYellow)
	case "blue":
		return color.New(color.FgBlue)
	case "magenta":
		return color.New(color.FgMagenta)
	case "cyan":
		return color.New(color.FgCyan)
	case "white":
		return color.New(color.FgWhite)
	case "hi-red":
		return color.New(color.FgHiRed)
	case "hi-green":
		return color.New(color.FgHiGreen)
	case "hi-yellow":
		return color.New(color.FgHiYellow)
	case "hi-blue":
		return color.New(color.FgHiBlue)
	case "hi-magenta":
		return color.New(color.FgHiMagenta)
	case "hi-cyan":
		return color.New(color.FgHiCyan)
	case "hi-white":
		return color.New(color.FgHiWhite)
	default:
		return color.New(color.FgWhite)
	}
}
