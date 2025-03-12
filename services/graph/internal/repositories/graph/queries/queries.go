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

	//go:embed get_nodes_with_rels.cypher
	GetNodesWithRels string

	//go:embed find_nearest_between_two_dijkstra.cypher
	FindNearestBetweenTwoDijkstra string

	//go:embed project_graph.cypher
	ProjectGraph string

	//go:embed check_graph_exists.cypher
	CheckGraphExists string

	//go:embed drop_graph.cypher
	DropGraph string
)
