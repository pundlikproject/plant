package consumer

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/pundliksarafar/kafka-parallel-consumer/config"
)

func GetConsumer(ctx context.Context, conf *config.Config) (*kafka.Consumer, error) {

	confMap := &kafka.ConfigMap{
		"bootstrap.servers": conf.KafkaUrl,
		"group.id":          conf.KafkaGroupId,
		"auto.offset.reset": "earliest",
	}
	return kafka.NewConsumer(confMap)
}
