MATCH (source:Node{number: $sourceNumber, graphID: $graphID}), (target:Node{number: $targetNumber, graphID: $graphID})
CALL gds.shortestPath.dijkstra.stream(
    $graphID,
    {
        sourceNode: source,
        targetNode: target,
        relationshipWeightProperty: 'weight'
    }
)
YIELD index, sourceNode, targetNode, totalCost, nodeIds, costs, path
RETURN 
    index,
    gds.util.asNode(sourceNode).number AS sourceNodeNumber,
    gds.util.asNode(targetNode).number AS targetNodeNumber,
    totalCost,
    [nodeId IN nodeIds | gds.util.asNode(nodeId).number] AS path,
    costs
ORDER BY index
