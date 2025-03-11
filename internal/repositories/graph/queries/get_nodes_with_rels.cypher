MATCH (source:Node{graphID: $graphID})-[r:CONNECTED_TO]->(target:Node{graphID: $graphID})
RETURN source.number AS sourceNumber, r.weight AS weight, target.number AS targetNumber
ORDER BY source.number DESC
