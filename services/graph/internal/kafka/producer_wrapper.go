package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/kuzin57/shad-networks/services/graph/internal/config"
	"github.com/kuzin57/shad-networks/services/graph/internal/consts"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ProducerWrapper struct {
	producer sarama.SyncProducer
	log      *zap.Logger
	conf     *config.KafkaConfig
}

func NewProducerWrapper(lc fx.Lifecycle, conf *config.Config, log *zap.Logger) *ProducerWrapper {
	p := &ProducerWrapper{
		conf: conf.Kafka,
		log:  log,
	}

	lc.Append(fx.Hook{
		OnStart: p.Start,
		OnStop:  p.Stop,
	})

	return p
}

func (p *ProducerWrapper) Start(ctx context.Context) error {
	producer, err := sarama.NewSyncProducer([]string{"kafka:9092"}, nil)
	if err != nil {
		return fmt.Errorf("failed to build producer: %w", err)
	}

	p.producer = producer

	return nil
}

func (p *ProducerWrapper) Stop(ctx context.Context) error {
	return p.producer.Close()
}

func (p *ProducerWrapper) ProducePathsProcessingMessage(ctx context.Context, msg PathsProcessingMessage) error {
	key := uuid.NewString()

	p.log.Sugar().Info("sending paths processing message")

	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal paths processing message: %w", err)
	}

	partition, offset, err := p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: consts.PathsProcessingTopic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(messageBytes),
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	p.log.Sugar().Infof("successfully sent message to partition: %d, offset: %d", partition, offset)

	return nil
}
