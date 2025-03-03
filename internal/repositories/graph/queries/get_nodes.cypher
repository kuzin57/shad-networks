MATCH (n:Node)
WHERE n.graphID = $graphID
RETURN n.number