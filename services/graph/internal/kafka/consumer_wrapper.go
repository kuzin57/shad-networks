package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/kuzin57/shad-networks/services/graph/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ConsumerWrapper struct {
	conf                     *config.KafkaConfig
	topicsPartitionConsumers map[string]map[int32]sarama.PartitionConsumer
	topicsConsumers          map[string]sarama.Consumer
	processors               map[string]Processor
	log                      *zap.Logger
}

func NewConsumerWrapper(
	lc fx.Lifecycle,
	conf *config.Config,
	log *zap.Logger,
	processorsHolder ProcessorsHolder,
) *ConsumerWrapper {
	cw := &ConsumerWrapper{
		conf:                     conf.Kafka,
		topicsPartitionConsumers: make(map[string]map[int32]sarama.PartitionConsumer),
		topicsConsumers:          make(map[string]sarama.Consumer),
		processors:               make(map[string]Processor),
		log:                      log,
	}

	for _, topicConf := range conf.Kafka.Topics {
		cw.processors[topicConf.Topic] = processorsHolder.GetProcessor(topicConf.Topic)
	}

	lc.Append(fx.Hook{
		OnStart: cw.Start,
		OnStop:  cw.Stop,
	})

	return cw
}

func (c *ConsumerWrapper) Start(ctx context.Context) error {
	for _, topicConf := range c.conf.Topics {
		consumer, err := sarama.NewConsumer([]string{fmt.Sprintf("%s:%d", c.conf.Host, c.conf.Port)}, nil)
		if err != nil {
			return fmt.Errorf("failed to build consumer: %w", err)
		}

		partitionsConsumers := make(map[int32]sarama.PartitionConsumer)
		for i := range topicConf.Partitions {
			partitionConsumer, err := consumer.ConsumePartition(topicConf.Topic, i, sarama.OffsetNewest)
			if err != nil {
				return fmt.Errorf("failed to build consumer for partition: %w", err)
			}

			partitionsConsumers[i] = partitionConsumer
		}

		c.topicsPartitionConsumers[topicConf.Topic] = partitionsConsumers
		c.topicsConsumers[topicConf.Topic] = consumer
	}

	for _, topicConf := range c.conf.Topics {
		for i := range topicConf.Partitions {
			go c.Listen(ctx, topicConf.Topic, i)
		}
	}

	return nil
}

func (c *ConsumerWrapper) Stop(ctx context.Context) error {
	for partition, consumer := range c.topicsConsumers {
		if err := consumer.Close(); err != nil {
			c.log.Sugar().Error("failed to close consumer for partition %s: error: %s", partition, err)
		}
	}

	return nil
}

func (c *ConsumerWrapper) Listen(ctx context.Context, topic string, partition int32) {
	partitionConsumer, ok := c.topicsPartitionConsumers[topic][partition]
	if !ok {
		return
	}

	wg := &sync.WaitGroup{}

ProcessMessagesLoop:
	for {
		select {
		case <-ctx.Done():
			c.log.Info("consumer is done")

			break ProcessMessagesLoop
		case msg, ok := <-partitionConsumer.Messages():
			if !ok {
				c.log.Sugar().Info("channel for topic %s, partition %s closed", topic, partition)

				break ProcessMessagesLoop
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				if err := c.processors[msg.Topic].Process(ctx, msg.Value); err != nil {
					c.log.Sugar().Error("process message from topic %s, error: %s", msg.Topic, err)
				}
			}()
		}
	}

	wg.Wait()
}
