package entities

type Graph struct {
	ID             string
	AdjencyMaxtrix [][]int
}

type GraphNode struct {
	GraphID string
	Number  int
}

type GraphEdge struct {
	StartNode int
	EndNode   int
}

type GraphParams struct {
	VerticesCount uint32
	Degrees       []uint32
	Weights       []uint32
}
