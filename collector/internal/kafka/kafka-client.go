package kafkaclient

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/terracegarden/collector/config"
	"github.com/terracegarden/framework/kafka/connector"
)

func PostMessage(ctx context.Context, key string, message []byte) error {
	conf := config.GetConfig(ctx)
	producer, err := connector.GetKafkaProducer(ctx, "")
	if err != nil {
		return err
	}
	producer.Produce(&kafka.Message{TopicPartition: kafka.TopicPartition{
		Topic: &conf.KafkaTopic,
	}, Value: message, Key: []byte(key)}, nil)

	return nil
}
