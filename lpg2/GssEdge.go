package lpg2

// GssEdge labels a predecessor link with the recognized grammar symbol,
// semantic value, token location, and optional SPPF node.
type GssEdge struct {
	predecessor *GssNode
	symbol      int
	location    int
	semantic    interface{}
	sppf        *SppfNode
}

func NewGssEdge(predecessor *GssNode, symbol int, location int,
	semantic interface{}, sppf *SppfNode) *GssEdge {
	return &GssEdge{
		predecessor: predecessor,
		symbol:      symbol,
		location:    location,
		semantic:    semantic,
		sppf:        sppf,
	}
}

func (e *GssEdge) GetPredecessor() *GssNode {
	return e.predecessor
}

func (e *GssEdge) GetSymbol() int {
	return e.symbol
}

func (e *GssEdge) GetLocation() int {
	return e.location
}

func (e *GssEdge) GetSemantic() interface{} {
	return e.semantic
}

func (e *GssEdge) GetSppf() *SppfNode {
	return e.sppf
}
