package visualizer

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/kuzin57/shad-networks/services/graph/internal/entities"
)

const (
	defaultPathEdgeWidth     = 5
	defaultPathLabelFontSize = 20
	defaultPathEdgeColor     = "red"
	defaultPathVertexColor   = "blue"
)

type Visualizer struct{}

func NewVisualizer() *Visualizer {
	return &Visualizer{}
}

func (v *Visualizer) Visualize(ctx context.Context, graph entities.Graph, paths []entities.Path) ([]byte, error) {
	canvas := graphviz.New()

	graphBuilder, err := canvas.Graph()
	if err != nil {
		return nil, fmt.Errorf("create graph builder: %w", err)
	}
	defer func() {
		if err := graphBuilder.Close(); err != nil {
			log.Fatal(err)
		}
		canvas.Close()
	}()

	nodes := make([]*cgraph.Node, 0, len(graph.AdjencyMaxtrix))

	for i := range graph.AdjencyMaxtrix {
		newNode, err := graphBuilder.CreateNode(fmt.Sprint(i))
		if err != nil {
			return nil, fmt.Errorf("create node: %w", err)
		}

		newNode.SetLabel(fmt.Sprint(i))
		nodes = append(nodes, newNode)
	}

	for i := range graph.AdjencyMaxtrix {
		for j := i + 1; j < len(graph.AdjencyMaxtrix); j++ {
			for k, weight := range graph.AdjencyMaxtrix[i][j] {
				edge, err := graphBuilder.CreateEdge(v.getEdgeName(i, j, k), nodes[i], nodes[j])
				if err != nil {
					return nil, fmt.Errorf("create edge: %w", err)
				}

				edge.SetLabel(fmt.Sprint(weight)).
					SetLabelDistance(0).
					SetDir(cgraph.NoneDir)

				if ok, dir := v.checkEdgeInPaths(i, j, weight, paths); ok {
					edge.SetPenWidth(defaultPathEdgeWidth).
						SetColor(defaultPathEdgeColor).
						SetLabelFontColor(defaultPathEdgeColor).
						SetLabelFontSize(defaultPathLabelFontSize).
						SetDir(dir)

					nodes[i].SetFillColor(defaultPathVertexColor)
				}
			}
		}
	}

	var buf bytes.Buffer

	if err := canvas.Render(graphBuilder, graphviz.SVG, &buf); err != nil {
		return nil, fmt.Errorf("render image: %w", err)
	}

	return []byte(base64.StdEncoding.EncodeToString(buf.Bytes())), nil
}

func (v *Visualizer) getEdgeName(source, target, index int) string {
	return fmt.Sprintf("edge|%d|%d|%d", source, target, index)
}

func (v *Visualizer) checkEdgeInPaths(source, target, weight int, paths []entities.Path) (bool, cgraph.DirType) {
	for _, path := range paths {
		for i := 1; i < len(path); i++ {
			if path[i-1].Vertex == source && path[i].Vertex == target && path[i].Weight == weight {
				return true, cgraph.ForwardDir
			}

			if path[i-1].Vertex == target && path[i].Vertex == source && path[i].Weight == weight {
				return true, cgraph.BackDir
			}
		}
	}

	return false, ""
}
