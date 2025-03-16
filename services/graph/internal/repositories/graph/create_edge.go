package graph

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/kuzin57/shad-networks/services/graph/internal/entities"
	"github.com/kuzin57/shad-networks/services/graph/internal/repositories/graph/queries"
	jsonutils "github.com/kuzin57/shad-networks/services/graph/internal/utils/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (r *Repository) СreateEdge(
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	edge *entities.GraphEdge,
) (any, error) {
	tmpl, err := template.New("create_edge_query").Parse(queries.CreateEdgeTemplate)
	if err != nil {
		return nil, fmt.Errorf("create template: %w", err)
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, edge); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}

	records, err := tx.Run(
		ctx,
		buf.String(),
		jsonutils.Serialize(edge),
	)
	if err != nil {
		return nil, fmt.Errorf("create edge: %w", err)
	}

	_, err = records.Single(ctx)
	if err != nil {
		return nil, fmt.Errorf("got single record: %w", err)
	}

	return nil, nil
}
