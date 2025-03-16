package processors

import (
	"github.com/kuzin57/shad-networks/services/graph/internal/config"
	"github.com/kuzin57/shad-networks/services/graph/internal/consts"
	"github.com/kuzin57/shad-networks/services/graph/internal/kafka"
	"github.com/kuzin57/shad-networks/services/graph/internal/pkg/processors/paths"
)

type Holder struct {
	processors map[string]kafka.Processor
}

func NewProcessorsHolder(conf *config.Config, pathProcessor *paths.PathsProcessor) *Holder {
	var (
		holder = &Holder{
			processors: make(map[string]kafka.Processor),
		}
		processorsMap = map[string]kafka.Processor{
			consts.PathsProcessingTopic: pathProcessor,
		}
	)

	for _, topicConf := range conf.Kafka.Topics {
		holder.processors[topicConf.Topic] = processorsMap[topicConf.Topic]
	}

	return holder
}

func (h *Holder) GetProcessor(name string) kafka.Processor {
	return h.processors[name]
}
