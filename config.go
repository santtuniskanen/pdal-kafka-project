package main

import "github.com/IBM/sarama"

type Config struct {
	KafkaBrokers   []string
	KafkaTopic     string
	ProducerConfig *sarama.Config
}

func NewConfig() *Config {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true

	return &Config{
		KafkaBrokers:   []string{"localhost:9092"},
		KafkaTopic:     "test-topic",
		ProducerConfig: saramaConfig,
	}
}
