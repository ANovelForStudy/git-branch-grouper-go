package model

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

func NewNode(name string) *Node {
	return &Node{
		Name:     name,
		Children: make(map[string]*Node),
		SubKeys:  []string{},
	}
}
