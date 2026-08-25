package main

type ConfigMain struct {
	Sparse bool `toml:"sparse"`
}

type Config struct {
	Main   ConfigMain       `toml:"main"`
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

func newNode(name string) *Node {
	return &Node{
		Name:     name,
		Children: make(map[string]*Node),
		SubKeys:  []string{},
	}
}
