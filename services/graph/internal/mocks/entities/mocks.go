package entitiesmocks

import "github.com/kuzin57/shad-networks/services/graph/internal/entities"

func GetMockGraph() entities.Graph {
	return entities.Graph{
		AdjencyMaxtrix: [][][]int{
			{nil, []int{1}, []int{2}},
			{[]int{1}, nil, []int{2}},
			{[]int{2}, []int{2}, nil},
		},
	}
}
