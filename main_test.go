package main

import (
	"sort"
	"testing"
)

func testNode(name string, children ...*Node) *Node {
	n := newNode(name)
	for _, c := range children {
		n.Children[c.Name] = c
		n.SubKeys = append(n.SubKeys, c.Name)
	}
	sort.Strings(n.SubKeys)
	return n
}

func leaf(name string) *Node {
	return testNode(name)
}

func branchData(groups ...*Node) BranchData {
	data := BranchData{
		MainGroups:      make(map[string]*Node),
		DefaultBranches: []string{},
		OtherBranches:   []string{},
	}
	for _, g := range groups {
		data.MainGroups[g.Name] = g
	}
	return data
}

func groupNames(data BranchData) []string {
	names := make([]string, 0, len(data.MainGroups))
	for name := range data.MainGroups {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func childNames(n *Node) []string {
	return n.SubKeys
}

func subtreeNames(n *Node) []string {
	names := []string{}
	var walk func(node *Node, prefix string)
	walk = func(node *Node, prefix string) {
		for _, sk := range node.SubKeys {
			child := node.Children[sk]
			full := prefix + sk
			names = append(names, full)
			if len(child.Children) > 0 {
				walk(child, full+"/")
			}
		}
	}
	walk(n, "")
	return names
}

func TestParsePrefixList_EmptyStringReturnsNil(t *testing.T) {
	result := parsePrefixList("")

	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestParsePrefixList_SingleValue(t *testing.T) {
	result := parsePrefixList("feat")

	if len(result) != 1 || result[0] != "feat" {
		t.Errorf("expected [feat], got %v", result)
	}
}

func TestParsePrefixList_MultipleValues(t *testing.T) {
	result := parsePrefixList("feat,fix,chore")

	expected := []string{"feat", "fix", "chore"}
	if len(result) != len(expected) {
		t.Fatalf("expected %d items, got %d", len(expected), len(result))
	}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("index %d: expected %q, got %q", i, v, result[i])
		}
	}
}

func TestParsePrefixList_TrimsSpaces(t *testing.T) {
	result := parsePrefixList(" feat , fix ")

	expected := []string{"feat", "fix"}
	if len(result) != len(expected) {
		t.Fatalf("expected %d items, got %d", len(expected), len(result))
	}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("index %d: expected %q, got %q", i, v, result[i])
		}
	}
}

func TestParsePrefixList_SkipsEmptySegments(t *testing.T) {
	result := parsePrefixList("feat,,fix,")

	expected := []string{"feat", "fix"}
	if len(result) != len(expected) {
		t.Fatalf("expected %d items, got %d", len(expected), len(result))
	}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("index %d: expected %q, got %q", i, v, result[i])
		}
	}
}

func TestBuildFilterTree_GroupLevelOnly(t *testing.T) {
	tree := buildFilterTree([]string{"feat", "fix"})

	if len(tree) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(tree))
	}
	if len(tree["feat"]) != 0 {
		t.Errorf("expected empty sub-paths for feat, got %v", tree["feat"])
	}
	if len(tree["fix"]) != 0 {
		t.Errorf("expected empty sub-paths for fix, got %v", tree["fix"])
	}
}

func TestBuildFilterTree_SubPathsGroupedUnderParent(t *testing.T) {
	tree := buildFilterTree([]string{"backup/v1", "backup/v2", "feat"})

	if len(tree) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(tree))
	}
	if len(tree["backup"]) != 2 {
		t.Errorf("expected 2 sub-paths for backup, got %d", len(tree["backup"]))
	}
	if len(tree["feat"]) != 0 {
		t.Errorf("expected empty sub-paths for feat, got %v", tree["feat"])
	}
}

func TestBuildFilterTree_DeepSubPaths(t *testing.T) {
	tree := buildFilterTree([]string{"backup/v2/refactor"})

	subs := tree["backup"]
	if len(subs) != 1 || subs[0] != "v2/refactor" {
		t.Errorf("expected [v2/refactor], got %v", subs)
	}
}

func TestFindMatchingChild_LeafNode(t *testing.T) {
	root := testNode("backup",
		testNode("v1",
			leaf("smth"),
		),
	)

	result := findMatchingChild(root, "v1")

	if result == nil || result.Name != "v1" {
		t.Errorf("expected v1 node, got %v", result)
	}
}

func TestFindMatchingChild_DeepPath(t *testing.T) {
	root := testNode("backup",
		testNode("v2",
			testNode("refactor",
				leaf("smth-useful"),
			),
		),
	)

	result := findMatchingChild(root, "v2/refactor")

	if result == nil || result.Name != "refactor" {
		t.Errorf("expected refactor node, got %v", result)
	}
}

func TestFindMatchingChild_NonExistentPath(t *testing.T) {
	root := testNode("backup",
		leaf("v1"),
	)

	result := findMatchingChild(root, "v99")

	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestFindMatchingChild_PartialPathReturnsIntermediate(t *testing.T) {
	root := testNode("backup",
		testNode("v1",
			testNode("refactor",
				leaf("smth"),
			),
		),
	)

	result := findMatchingChild(root, "v1")

	if result == nil || result.Name != "v1" {
		t.Errorf("expected v1 node, got %v", result)
	}
	if _, ok := result.Children["refactor"]; !ok {
		t.Error("expected v1 to have refactor child")
	}
}

func TestCopySubtree_CreatesDeepCopy(t *testing.T) {
	src := testNode("a",
		testNode("b",
			leaf("c"),
		),
		leaf("d"),
	)
	dst := newNode("root")

	copySubtree(dst, src)

	if len(dst.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(dst.Children))
	}
	if dst.Children["b"] == src.Children["b"] {
		t.Error("expected deep copy, got shared pointer")
	}
	if dst.Children["b"].Children["c"] == src.Children["b"].Children["c"] {
		t.Error("expected deep copy at nested level, got shared pointer")
	}
}

func TestCopySubtree_PreservesOrder(t *testing.T) {
	src := testNode("root",
		leaf("c"),
		leaf("a"),
		leaf("b"),
	)
	dst := newNode("dst")

	copySubtree(dst, src)

	if len(dst.SubKeys) != 3 {
		t.Fatalf("expected 3 sub-keys, got %d", len(dst.SubKeys))
	}
	expected := []string{"a", "b", "c"}
	for i, k := range expected {
		if dst.SubKeys[i] != k {
			t.Errorf("SubKeys[%d]: expected %q, got %q", i, k, dst.SubKeys[i])
		}
	}
}

func TestCopyFilteredSubtree_ExcludesTopLevelChild(t *testing.T) {
	src := testNode("backup",
		leaf("v1"),
		leaf("v2"),
		leaf("feat"),
	)

	result := copyFilteredSubtree(src, []string{"v2"})

	if len(result.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(result.Children))
	}
	if _, ok := result.Children["v2"]; ok {
		t.Error("v2 should be excluded")
	}
	if _, ok := result.Children["v1"]; !ok {
		t.Error("v1 should be present")
	}
	if _, ok := result.Children["feat"]; !ok {
		t.Error("feat should be present")
	}
}

func TestCopyFilteredSubtree_ExcludesNestedSubtree(t *testing.T) {
	src := testNode("backup",
		testNode("v2",
			testNode("refactor",
				leaf("smth"),
			),
		),
		leaf("feat"),
	)

	result := copyFilteredSubtree(src, []string{"v2/refactor"})

	v2 := result.Children["v2"]
	if v2 == nil {
		t.Fatal("v2 should be present")
	}
	if _, ok := v2.Children["refactor"]; ok {
		t.Error("v2/refactor should be excluded")
	}
}

func TestCopyFilteredSubtree_ParentBecomesLeafWhenAllChildrenExcluded(t *testing.T) {
	src := testNode("backup",
		testNode("v2",
			leaf("only-child"),
		),
	)

	result := copyFilteredSubtree(src, []string{"v2/only-child"})

	v2 := result.Children["v2"]
	if v2 == nil {
		t.Fatal("v2 should be present")
	}
	if len(v2.Children) != 0 {
		t.Errorf("v2 should have no children, got %d", len(v2.Children))
	}
}

func TestCopyFilteredSubtree_MultipleExcludePaths(t *testing.T) {
	src := testNode("backup",
		leaf("v1"),
		leaf("v2"),
		leaf("feat"),
	)

	result := copyFilteredSubtree(src, []string{"v2", "feat"})

	if len(result.Children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(result.Children))
	}
	if _, ok := result.Children["v1"]; !ok {
		t.Error("v1 should be the only remaining child")
	}
}

func TestCopyFilteredSubtree_NoExcludesReturnsFullCopy(t *testing.T) {
	src := testNode("backup",
		leaf("v1"),
		leaf("v2"),
	)

	result := copyFilteredSubtree(src, nil)

	if len(result.Children) != 2 {
		t.Fatalf("expected 2 children, got %d", len(result.Children))
	}
	if result.Children["v1"] == src.Children["v1"] {
		t.Error("expected deep copy")
	}
}

func TestFilterBranchData_EmptyFiltersReturnsOriginal(t *testing.T) {
	data := branchData(
		testNode("feat", leaf("f-1")),
		testNode("fix", leaf("fix-1")),
	)

	result := filterBranchData(data, nil, nil)

	if len(result.MainGroups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(result.MainGroups))
	}
}

func TestFilterBranchData_GroupLevelInclude(t *testing.T) {
	data := branchData(
		testNode("feat", leaf("f-1")),
		testNode("fix", leaf("fix-1")),
		testNode("chore", leaf("c-1")),
	)

	result := filterBranchData(data, []string{"feat", "fix"}, nil)

	names := groupNames(result)
	if len(names) != 2 || names[0] != "feat" || names[1] != "fix" {
		t.Errorf("expected [feat fix], got %v", names)
	}
}

func TestFilterBranchData_GroupLevelExclude(t *testing.T) {
	data := branchData(
		testNode("feat", leaf("f-1")),
		testNode("fix", leaf("fix-1")),
		testNode("chore", leaf("c-1")),
	)

	result := filterBranchData(data, nil, []string{"fix"})

	names := groupNames(result)
	if len(names) != 2 || names[0] != "chore" || names[1] != "feat" {
		t.Errorf("expected [chore feat], got %v", names)
	}
}

func TestFilterBranchData_ExcludeAllGroupsResultEmpty(t *testing.T) {
	data := branchData(
		testNode("feat", leaf("f-1")),
		testNode("fix", leaf("fix-1")),
	)

	result := filterBranchData(data, nil, []string{"feat", "fix"})

	if len(result.MainGroups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(result.MainGroups))
	}
}

func TestFilterBranchData_NonExistentGroupIgnored(t *testing.T) {
	data := branchData(
		testNode("feat", leaf("f-1")),
	)

	result := filterBranchData(data, []string{"nonexistent"}, nil)

	if len(result.MainGroups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(result.MainGroups))
	}
}

func TestFilterBranchData_SubPathInclude(t *testing.T) {
	data := branchData(
		testNode("backup",
			testNode("v1", leaf("smth")),
			testNode("v2", leaf("smth")),
			testNode("feat", leaf("some-feature")),
		),
	)

	result := filterBranchData(data, []string{"backup/v1"}, nil)

	backup := result.MainGroups["backup"]
	if backup == nil {
		t.Fatal("backup group should be present")
	}
	if len(backup.Children) != 1 {
		t.Fatalf("expected 1 child under backup, got %d", len(backup.Children))
	}
	if _, ok := backup.Children["v1"]; !ok {
		t.Error("v1 should be present under backup")
	}
}

func TestFilterBranchData_SubPathExclude(t *testing.T) {
	data := branchData(
		testNode("backup",
			testNode("v1", leaf("smth")),
			testNode("v2", leaf("smth")),
			testNode("feat", leaf("some-feature")),
		),
	)

	result := filterBranchData(data, nil, []string{"backup/v2"})

	backup := result.MainGroups["backup"]
	if backup == nil {
		t.Fatal("backup group should be present")
	}
	if _, ok := backup.Children["v2"]; ok {
		t.Error("v2 should be excluded")
	}
	if _, ok := backup.Children["v1"]; !ok {
		t.Error("v1 should remain")
	}
	if _, ok := backup.Children["feat"]; !ok {
		t.Error("feat should remain")
	}
}

func TestFilterBranchData_MultipleSubPathIncludes(t *testing.T) {
	data := branchData(
		testNode("backup",
			testNode("v1", leaf("smth")),
			testNode("v2", leaf("smth")),
			testNode("feat", leaf("some-feature")),
		),
		testNode("feat", leaf("f-1")),
	)

	result := filterBranchData(data, []string{"backup/v1", "feat"}, nil)

	names := groupNames(result)
	if len(names) != 2 {
		t.Fatalf("expected 2 groups, got %d: %v", len(names), names)
	}
	backup := result.MainGroups["backup"]
	if backup == nil || len(backup.Children) != 1 {
		t.Error("backup should have exactly 1 child (v1)")
	}
	if _, ok := result.MainGroups["feat"]; !ok {
		t.Error("feat group should be present")
	}
}

func TestFilterBranchData_MixedGroupAndSubPathInclude(t *testing.T) {
	data := branchData(
		testNode("backup",
			testNode("v1", leaf("smth")),
			testNode("v2", leaf("smth")),
		),
		testNode("feat", leaf("f-1")),
		testNode("fix", leaf("fix-1")),
	)

	result := filterBranchData(data, []string{"backup/v1", "feat"}, nil)

	names := groupNames(result)
	if len(names) != 2 || names[0] != "backup" || names[1] != "feat" {
		t.Errorf("expected [backup feat], got %v", names)
	}
	backup := result.MainGroups["backup"]
	if backup == nil || len(backup.Children) != 1 {
		t.Error("backup should have exactly 1 child (v1)")
	}
}

func TestFilterBranchData_DeepSubPathExclude(t *testing.T) {
	data := branchData(
		testNode("backup",
			testNode("v2",
				testNode("refactor",
					leaf("smth-useful"),
				),
			),
			testNode("v1",
				testNode("refactor",
					leaf("smth-useful"),
				),
			),
			testNode("feat", leaf("some-feature")),
		),
	)

	result := filterBranchData(data, nil, []string{"backup/v2/refactor"})

	backup := result.MainGroups["backup"]
	if backup == nil {
		t.Fatal("backup group should be present")
	}
	v2 := backup.Children["v2"]
	if v2 == nil {
		t.Fatal("v2 should be present")
	}
	if _, ok := v2.Children["refactor"]; ok {
		t.Error("v2/refactor should be excluded")
	}
	v1 := backup.Children["v1"]
	if v1 == nil {
		t.Fatal("v1 should be present")
	}
	if _, ok := v1.Children["refactor"]; !ok {
		t.Error("v1/refactor should remain")
	}
}

func TestFilterBranchData_DeepSubPathInclude(t *testing.T) {
	data := branchData(
		testNode("backup",
			testNode("v2",
				testNode("refactor",
					leaf("smth-useful"),
				),
			),
			testNode("v1",
				testNode("refactor",
					leaf("smth-useful"),
				),
			),
		),
	)

	result := filterBranchData(data, []string{"backup/v2/refactor"}, nil)

	backup := result.MainGroups["backup"]
	if backup == nil {
		t.Fatal("backup group should be present")
	}
	names := childNames(backup)
	if len(names) != 1 || names[0] != "v2" {
		t.Errorf("expected [v2], got %v", names)
	}
	v2 := backup.Children["v2"]
	if v2 == nil {
		t.Fatal("v2 should be present")
	}
	refNames := childNames(v2)
	if len(refNames) != 1 || refNames[0] != "refactor" {
		t.Errorf("expected [refactor], got %v", refNames)
	}
}

func TestFilterBranchData_ExcludeAndIncludeSameGroup_IncludeWins(t *testing.T) {
	data := branchData(
		testNode("feat", leaf("f-1")),
		testNode("fix", leaf("fix-1")),
	)

	result := filterBranchData(data, []string{"feat"}, []string{"feat"})

	names := groupNames(result)
	if len(names) != 1 || names[0] != "feat" {
		t.Errorf("include should take priority, expected [feat], got %v", names)
	}
}

func TestFilterBranchData_ExcludeMultipleSubPaths(t *testing.T) {
	data := branchData(
		testNode("backup",
			testNode("v1", leaf("smth")),
			testNode("v2", leaf("smth")),
			testNode("feat", leaf("some-feature")),
		),
	)

	result := filterBranchData(data, nil, []string{"backup/v2", "backup/feat"})

	backup := result.MainGroups["backup"]
	if backup == nil {
		t.Fatal("backup group should be present")
	}
	names := childNames(backup)
	if len(names) != 1 || names[0] != "v1" {
		t.Errorf("expected [v1], got %v", names)
	}
}

func TestFilterBranchData_GroupExcludeAndSubPathExclude(t *testing.T) {
	data := branchData(
		testNode("backup",
			testNode("v1", leaf("smth")),
			testNode("v2", leaf("smth")),
		),
		testNode("old", leaf("o-1")),
	)

	result := filterBranchData(data, nil, []string{"old", "backup/v2"})

	names := groupNames(result)
	if len(names) != 1 || names[0] != "backup" {
		t.Errorf("expected [backup], got %v", names)
	}
	backup := result.MainGroups["backup"]
	if _, ok := backup.Children["v2"]; ok {
		t.Error("backup/v2 should be excluded")
	}
	if _, ok := backup.Children["v1"]; !ok {
		t.Error("backup/v1 should remain")
	}
}

func TestFilterBranchData_IsolatedFromOriginal(t *testing.T) {
	data := branchData(
		testNode("feat", leaf("f-1")),
	)

	_ = filterBranchData(data, []string{"feat"}, nil)

	feat := data.MainGroups["feat"]
	if feat == nil {
		t.Fatal("original data should be unchanged")
	}
	if _, ok := feat.Children["f-1"]; !ok {
		t.Error("original feat should still have f-1")
	}
}

func TestFilterBranchData_ExcludeDoesNotMutateOriginal(t *testing.T) {
	data := branchData(
		testNode("backup",
			leaf("v1"),
			leaf("v2"),
		),
	)

	_ = filterBranchData(data, nil, []string{"backup/v2"})

	backup := data.MainGroups["backup"]
	if backup == nil {
		t.Fatal("original backup should exist")
	}
	if len(backup.Children) != 2 {
		t.Errorf("original backup should have 2 children, got %d", len(backup.Children))
	}
}

func TestFilterBranchData_WholeGroupIncludeDeepCopies(t *testing.T) {
	data := branchData(
		testNode("feat", leaf("f-1")),
	)

	result := filterBranchData(data, []string{"feat"}, nil)

	feat := result.MainGroups["feat"]
	if feat == nil {
		t.Fatal("feat should be present")
	}
	if feat == data.MainGroups["feat"] {
		t.Error("expected deep copy, got shared pointer")
	}
}

func TestFilterBranchData_WholeGroupExcludeDeepCopies(t *testing.T) {
	data := branchData(
		testNode("backup",
			leaf("v1"),
			leaf("v2"),
		),
	)

	result := filterBranchData(data, nil, []string{"backup/v1"})

	backup := result.MainGroups["backup"]
	if backup == nil {
		t.Fatal("backup should be present")
	}
	if backup == data.MainGroups["backup"] {
		t.Error("expected deep copy, got shared pointer")
	}
}
