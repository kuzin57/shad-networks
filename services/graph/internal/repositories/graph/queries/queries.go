package queries

import (
	_ "embed"
)

var (
	//go:embed create_database.cypher
	CreateDatabase string

	//go:embed create_edge.cypher.tmpl
	CreateEdgeTemplate string

	//go:embed create_node.cypher
	CreateNode string

	//go:embed get_nodes_with_rels.cypher
	GetNodesWithRels string

	//go:embed find_nearest_dijkstra.cypher
	FindNearestDijkstra string

	//go:embed find_k_nearest_yens.cypher
	FindKNearestYens string

	//go:embed project_graph.cypher
	ProjectGraph string

	//go:embed check_graph_exists.cypher
	CheckGraphExists string

	//go:embed drop_graph.cypher
	DropGraph string
)
