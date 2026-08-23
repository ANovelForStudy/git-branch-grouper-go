package main

import (
	"fmt"
	"sort"

	"github.com/fatih/color"
)

func printResults(data BranchData) {
	groupNames := make([]string, 0, len(data.MainGroups))
	for name := range data.MainGroups {
		groupNames = append(groupNames, name)
	}
	sort.Strings(groupNames)
	sort.Strings(data.DefaultBranches)
	sort.Strings(data.OtherBranches)

	tabulation := "    "

	color.Red("[default]")
	for _, b := range data.DefaultBranches {
		fmt.Println(tabulation, b)
	}

	for _, gName := range groupNames {
		color.Red("[%s]", gName)
		for _, b := range data.MainGroups[gName] {
			fmt.Println(tabulation, b)
		}
	}

	color.Red("[other]")
	for _, b := range data.OtherBranches {
		fmt.Println(tabulation, b)
	}
}
