package visualizer

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/kuzin57/shad-networks/services/graph/internal/entities"
)

type Visualizer struct{}

func NewVisualizer() *Visualizer {
	return &Visualizer{}
}

func (v *Visualizer) Visualize(ctx context.Context, graph entities.Graph, path entities.Path) ([]byte, error) {
	canvas, err := graphviz.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("create new canvas: %w", err)
	}

	graphBuilder, err := canvas.Graph()
	if err != nil {
		return nil, fmt.Errorf("create graph builder: %w", err)
	}

	nodes := make([]*cgraph.Node, 0, len(graph.AdjencyMaxtrix))

	for i := range graph.AdjencyMaxtrix {
		newNode, err := graphBuilder.CreateNodeByName(fmt.Sprint(i))
		if err != nil {
			return nil, fmt.Errorf("create node: %w", err)
		}

		newNode.SetLabel(fmt.Sprint(i))
		nodes = append(nodes, newNode)
	}

	for i := range graph.AdjencyMaxtrix {
		for j := i + 1; j < len(graph.AdjencyMaxtrix); j++ {
			for k, weight := range graph.AdjencyMaxtrix[i][j] {
				edge, err := graphBuilder.CreateEdgeByName(v.getEdgeName(i, j, k), nodes[i], nodes[j])
				if err != nil {
					return nil, fmt.Errorf("create edge: %w", err)
				}

				edge.SetLabel(fmt.Sprint(weight))

				if v.checkEdgeInPath(i, j, weight, path) {
					edge.SetColor("red")
				}
			}
		}
	}

	var buf bytes.Buffer

	if err := canvas.Render(ctx, graphBuilder, graphviz.SVG, &buf); err != nil {
		return nil, fmt.Errorf("render image: %w", err)
	}

	return []byte(base64.StdEncoding.EncodeToString(buf.Bytes())), nil
}

func (v *Visualizer) getEdgeName(source, target, index int) string {
	return fmt.Sprintf("edge|%d|%d|%d", source, target, index)
}

func (v *Visualizer) checkEdgeInPath(source, target, weight int, path []entities.PathPart) bool {
	for i := 1; i < len(path); i++ {
		if path[i-1].Vertex == source && path[i].Vertex == target && path[i].Weight == weight {
			return true
		}
	}

	return false
}
