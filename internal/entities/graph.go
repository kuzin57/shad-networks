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
	Weight int // weight of edge which comes to this vertex
}

type GraphParams struct {
	VerticesCount uint32
	Degrees       []uint32
	Weights       []uint32
}
