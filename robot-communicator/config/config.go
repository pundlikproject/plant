package config

import (
	"context"
	"log"

	env "github.com/Netflix/go-env"
)

type Configs struct {
	KafkaUrl   string `env:"KAFKA_URL,default=localhost:32111"`
	KafkaTopic string `env:"KAFKA_URL,default=topic_name"`
}

var conf *Configs

func GetConfig(ctx context.Context) *Configs {
	if conf == nil {
		_, err := env.UnmarshalFromEnviron(&conf)
		if err != nil {
			log.Fatal(err)
		}
	}
	return conf
}
