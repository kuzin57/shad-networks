package graph

import (
	"fmt"
	"reflect"

	"github.com/kuzin57/shad-networks/services/graph/internal/entities"
	"github.com/kuzin57/shad-networks/services/graph/internal/utils/convert"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (r *Repository) getPathsFromResult(result *neo4j.EagerResult) ([]entities.Path, error) {
	if len(result.Records) == 0 {
		return nil, fmt.Errorf("invalid records len: %d", len(result.Records))
	}

	paths := make([]entities.Path, 0, len(result.Records))

	for _, record := range result.Records {
		path, err := r.parseRecord(record)
		if err != nil {
			return nil, err
		}

		paths = append(paths, path)
	}

	return paths, nil
}

func (r *Repository) parseRecord(record *neo4j.Record) (entities.Path, error) {
	if record == nil {
		return nil, fmt.Errorf("nil record")
	}

	pathValue, ok := record.Get(pathParam)
	if !ok {
		return nil, fmt.Errorf("path param not found")
	}

	costsValue, ok := record.Get(costsParam)
	if !ok {
		return nil, fmt.Errorf("costs param not found")
	}

	r.log.Sugar().Infof("pathValue: %s", reflect.TypeOf(pathValue).String())

	var (
		path, _    = pathValue.([]any)
		costs, _   = costsValue.([]any)
		resultPath = make(entities.Path, 0, len(path))
	)

	for i := range path {
		if i == 0 {
			resultPath = append(resultPath, entities.PathPart{
				Vertex: convert.FloatAnyToInt(path[i]),
				Weight: 0},
			)

			continue
		}

		resultPath = append(resultPath, entities.PathPart{
			Vertex: convert.FloatAnyToInt(path[i]),
			Weight: convert.FloatAnyToInt(costs[i]) - convert.FloatAnyToInt(costs[i-1]),
		})
	}

	r.log.Sugar().Infof("found path: %v", resultPath)

	return resultPath, nil
}
