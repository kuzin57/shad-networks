MATCH (source:Node{graphID: $graphID})-[r:CONNECTED_TO]->(target:Node{graphID: $graphID})
RETURN gds.graph.project(
  $graphID,
  source,
  target,
  { relationshipProperties: r { .weight } },
  { undirectedRelationshipTypes: ['*'] }
)
