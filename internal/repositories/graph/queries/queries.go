package queries

import (
	_ "embed"
)

var (
	//go:embed create_database.cypher
	CreateDatabase string

	//go:embed create_edge.cypher
	CreateEdge string

	//go:embed create_node.cypher
	CreateNode string

	//go:embed get_nodes.cypher
	GetNodes string

	//go:embed find_nearest_by_djikstra.cypher
	FindNearestByDjikstra string
)
