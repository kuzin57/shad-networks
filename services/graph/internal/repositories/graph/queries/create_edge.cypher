MATCH (firstNode:Node {graphID: $graphID, number: $from}), (secondNode:Node {graphID: $graphID, number: $to})
CREATE (firstNode)-[:CONNECTED_TO {weight: $weight}]->(secondNode)
RETURN firstNode, secondNode
