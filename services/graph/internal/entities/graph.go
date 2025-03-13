package entities

type Graph struct {
	ID             string    `json:"graphID"`
	AdjencyMaxtrix [][][]int `json:"-"`
}

type GraphNode struct {
	GraphID string `json:"graphID"`
	Number  int    `json:"number"`
}

type GraphEdge struct {
	From       int    `json:"from"`
	To         int    `json:"to"`
	GraphID    string `json:"graphID"`
	Weight     int    `json:"weight"`
	Connection string `json:"connection"`
}

type Path []PathPart

type PathPart struct {
	Vertex int
	// Weight of edge which comes to this vertex
	// Example:
	//     w1       w2       w3
	// 1 -----> 2 -----> 3 -----> 4
	//
	// {1, 0}; {2, w1}; {3, w2}; {4, w3}
	Weight int
}

type GraphParams struct {
	VerticesCount    uint32
	Degrees          []uint32
	Weights          []uint32
	MaxMultipleEdges uint32
}
