package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
)

func printResultsWithConfig(data BranchData, cfg Config) {
	defaultColorStr := cfg.Colors["default"]
	otherColorStr := cfg.Colors["other"]
	starColorStr := cfg.Colors["active_star"]

	printGroup("[default]", data.DefaultBranches, defaultColorStr, starColorStr)

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
		printGroup(title, data.MainGroups[groupName], groupColorStr, starColorStr)
	}

	printGroup("[other]", data.OtherBranches, otherColorStr, starColorStr)
}

func printGroup(title string, branches []string, titleColorStr string, starColorStr string) {
	titleColor := getColorPrinter(titleColorStr)
	titleColor.Println(title)

	starColor := getColorPrinter(starColorStr)
	tabulation := "    "

	sort.Strings(branches)
	for _, branch := range branches {
		if strings.HasPrefix(branch, "* ") {
			fmt.Print(tabulation)
			starColor.Print("* ")
			fmt.Println(strings.TrimPrefix(branch, "* "))
		} else {
			fmt.Println(tabulation, branch)
		}
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
