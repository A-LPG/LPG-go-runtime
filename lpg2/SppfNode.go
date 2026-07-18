package lpg2

// SppfNode is a shared packed parse forest symbol node keyed by grammar
// symbol and token extent.
type SppfNode struct {
	grammarSymbol int
	leftExtent    int
	rightExtent   int
	packs         []*SppfPackedNode
	astForest     interface{}
}

func NewSppfNode(grammarSymbol int, leftExtent int, rightExtent int) *SppfNode {
	return &SppfNode{
		grammarSymbol: grammarSymbol,
		leftExtent:    leftExtent,
		rightExtent:   rightExtent,
	}
}

func (n *SppfNode) GetGrammarSymbol() int {
	return n.grammarSymbol
}

func (n *SppfNode) GetLeftExtent() int {
	return n.leftExtent
}

func (n *SppfNode) GetRightExtent() int {
	return n.rightExtent
}

func (n *SppfNode) GetPacks() []*SppfPackedNode {
	return n.packs
}

func (n *SppfNode) GetAstForest() interface{} {
	return n.astForest
}

// SppfPackedNode is one production alternative under an SPPF symbol node.
type SppfPackedNode struct {
	rule     int
	children []*SppfNode
	semantic interface{}
}

func NewSppfPackedNode(rule int, children []*SppfNode,
	semantic interface{}) *SppfPackedNode {
	if children == nil {
		children = []*SppfNode{}
	}
	return &SppfPackedNode{
		rule:     rule,
		children: children,
		semantic: semantic,
	}
}

func (n *SppfPackedNode) GetRule() int {
	return n.rule
}

func (n *SppfPackedNode) GetChildren() []*SppfNode {
	out := make([]*SppfNode, 0, len(n.children))
	for _, child := range n.children {
		if child != nil {
			out = append(out, child)
		}
	}
	return out
}

func (n *SppfPackedNode) GetSemantic() interface{} {
	return n.semantic
}
