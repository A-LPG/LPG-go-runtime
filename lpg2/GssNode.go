package lpg2

// GssNode is a graph-structured stack node: an LR state at an input index.
type GssNode struct {
	state int
	index int
	edges []*GssEdge
}

func NewGssNode(state int, index int) *GssNode {
	return &GssNode{state: state, index: index}
}

func (n *GssNode) GetState() int {
	return n.state
}

func (n *GssNode) GetIndex() int {
	return n.index
}

func (n *GssNode) GetEdges() []*GssEdge {
	return n.edges
}
