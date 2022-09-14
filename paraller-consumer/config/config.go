package config

type ParallelismConfig string

var Partition_ParallelismConfig ParallelismConfig = "partition"
var Key_ParallelismConfig ParallelismConfig = "key"

type Config struct {
	KafkaUrl string

	KafkaGroupId string
	KafkaTopic   string

	Parallelism         ParallelismConfig
	ParallelThreadCount int
}

func WithKafkaUrl(kafkaUrl string) func(c *Config) {
	return func(c *Config) {
		c.KafkaUrl = kafkaUrl
	}
}

func WithParallelism(parallelismConf ParallelismConfig) func(c *Config) {
	return func(c *Config) {
		c.Parallelism = parallelismConf
	}
}

func WithMaxParallelThread(tCount int) func(c *Config) {
	return func(c *Config) {
		c.ParallelThreadCount = tCount
	}
}

func WithGroupId(kGroupId string) func(c *Config) {
	return func(c *Config) {
		c.KafkaGroupId = kGroupId
	}
}

func WithTopic(topic string) func(c *Config) {
	return func(c *Config) {
		c.KafkaTopic = topic
	}
}
