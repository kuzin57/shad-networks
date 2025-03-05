package entitiesmocks

import "github.com/kuzin57/shad-networks/internal/entities"

func GetMockGraph() entities.Graph {
	return entities.Graph{
		AdjencyMaxtrix: [][]int{
			{0, 1, 2},
			{1, 0, 2},
			{2, 2, 0},
		},
	}
}
