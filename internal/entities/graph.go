package entities

type Graph struct {
	ID             string
	AdjencyMaxtrix [][]int
}

type GraphParams struct {
	Name          string
	VerticesCount uint32
	Degrees       []uint32
	Weights       []uint32
}
