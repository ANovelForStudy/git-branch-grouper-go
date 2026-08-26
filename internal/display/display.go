package display

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"

	"git-branch-grouper-plugin/internal/config"
	"git-branch-grouper-plugin/internal/model"
)

const ansiForeground = 38

func PrintResults(w io.Writer, data model.BranchData, cfg config.Config, hasFilter bool) {
	starColorStr := cfg.Colors["active_star"]
	marker := cfg.Format.BranchMarker + " "

	if !hasFilter {
		defaultColorStr := cfg.Colors["default"]
		title := strings.ReplaceAll(cfg.Format.GroupPrefix, "{group}", "default")
		printSimpleGroup(w, title, data.DefaultBranches, defaultColorStr, starColorStr, marker, cfg.Format.Indent, cfg.Main.Sparse)
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

		title := groupName + "/"
		titleColor := getColorPrinter(groupColorStr)
		_, _ = titleColor.Fprintln(w, title)

		rootNode := data.MainGroups[groupName]
		sort.Strings(rootNode.SubKeys)
		for _, subKey := range rootNode.SubKeys {
			printNode(w, rootNode.Children[subKey], 1, groupColorStr, starColorStr, marker, cfg.Format.Indent, cfg)
		}

		if cfg.Main.Sparse {
			_, _ = fmt.Fprintln(w)
		}
	}

	if !hasFilter {
		otherColorStr := cfg.Colors["other"]
		title := strings.ReplaceAll(cfg.Format.GroupPrefix, "{group}", "other")
		printSimpleGroup(w, title, data.OtherBranches, otherColorStr, starColorStr, marker, cfg.Format.Indent, cfg.Main.Sparse)
	}
}

func printNode(w io.Writer, node *model.Node, level int, parentColorStr string, starColorStr string, marker string, indent string, cfg config.Config) {
	tabulation := strings.Repeat(indent, level)
	starColor := getColorPrinter(starColorStr)

	nodeColorStr := cfg.Groups[node.Name]
	if nodeColorStr == "" {
		nodeColorStr = parentColorStr
	}
	nodeColor := getColorPrinter(nodeColorStr)

	if len(node.Children) == 0 {
		if node.IsActive {
			_, _ = fmt.Fprint(w, tabulation)
			_, _ = starColor.Fprint(w, marker)
			_, _ = nodeColor.Fprintln(w, node.Name)
		} else {
			_, _ = nodeColor.Fprintln(w, tabulation+node.Name)
		}
		return
	}

	if node.IsActive {
		_, _ = fmt.Fprint(w, tabulation)
		_, _ = starColor.Fprint(w, marker)
		_, _ = nodeColor.Fprintln(w, node.Name+"/")
	} else {
		_, _ = nodeColor.Fprintln(w, tabulation+node.Name+"/")
	}

	sort.Strings(node.SubKeys)
	for _, subKey := range node.SubKeys {
		printNode(w, node.Children[subKey], level+1, nodeColorStr, starColorStr, marker, indent, cfg)
	}
}

func printSimpleGroup(w io.Writer, title string, branches []string, titleColorStr string, starColorStr string, marker string, indent string, sparse bool) {
	titleColor := getColorPrinter(titleColorStr)
	_, _ = titleColor.Fprintln(w, title)

	starColor := getColorPrinter(starColorStr)

	sort.Strings(branches)
	for _, branch := range branches {
		if strings.HasPrefix(branch, "* ") {
			_, _ = fmt.Fprint(w, indent)
			_, _ = starColor.Fprint(w, marker)
			_, _ = fmt.Fprintln(w, strings.TrimPrefix(branch, "* "))
		} else {
			_, _ = fmt.Fprintln(w, indent, branch)
		}
	}

	if sparse {
		_, _ = fmt.Fprintln(w)
	}
}

func getColorPrinter(colorName string) *color.Color {
	switch strings.ToLower(colorName) {
	case "black":
		return color.New(color.FgBlack)
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
	case "hi-black":
		return color.New(color.FgHiBlack)
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
		return parseExtendedColor(colorName)
	}
}

func parseExtendedColor(s string) *color.Color {
	s = strings.TrimSpace(s)

	if strings.HasPrefix(s, "#") {
		if c, ok := parseHexColor(s); ok {
			return c
		}
		return color.New(color.FgWhite)
	}

	if code, err := strconv.Atoi(s); err == nil && code >= 0 && code <= 255 {
		return color.New(color.Attribute(ansiForeground), color.Attribute(5), color.Attribute(code))
	}

	return color.New(color.FgWhite)
}

func parseHexColor(hex string) (*color.Color, bool) {
	hex = strings.TrimPrefix(hex, "#")

	switch len(hex) {
	case 6:
		r, err1 := strconv.ParseInt(hex[0:2], 16, 32)
		g, err2 := strconv.ParseInt(hex[2:4], 16, 32)
		b, err3 := strconv.ParseInt(hex[4:6], 16, 32)
		if err1 != nil || err2 != nil || err3 != nil {
			return nil, false
		}
		return color.RGB(int(r), int(g), int(b)), true
	case 3:
		r, err1 := strconv.ParseInt(string(hex[0])+string(hex[0]), 16, 32)
		g, err2 := strconv.ParseInt(string(hex[1])+string(hex[1]), 16, 32)
		b, err3 := strconv.ParseInt(string(hex[2])+string(hex[2]), 16, 32)
		if err1 != nil || err2 != nil || err3 != nil {
			return nil, false
		}
		return color.RGB(int(r), int(g), int(b)), true
	default:
		return nil, false
	}
}
