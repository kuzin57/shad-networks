MATCH (firstNode:Node {graphID: $graphID, number: $firstNumber}), (secondNode:Node {graphID: $graphID, number: $secondNumber})
CREATE (firstNode)-[:CONNECTED_WITH {weight: $weight}]->(secondNode)
RETURN firstNode, secondNode
