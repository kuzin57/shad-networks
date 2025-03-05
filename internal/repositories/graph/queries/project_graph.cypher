MATCH (source:Node)-[r:CONNECTED_TO]->(target:Node)
RETURN gds.graph.project(
  $graphID,
  source,
  target,
  { relationshipProperties: r { .weight } }
)