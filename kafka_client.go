// I'm going to wrap Kafka functions here to easier to use functions
// that I can use later with the workers to push data to Kafka.
package main

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/IBM/sarama"
)

type KafkaClient struct {
	producer     sarama.AsyncProducer
	topic        string
	successCount int64
	errorCount   int64
	wg           sync.WaitGroup
}

func NewKafkaClient(brokers []string, topic string, config *sarama.Config) (*KafkaClient, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("no brokers provided")
	}
	if topic == "" {
		return nil, fmt.Errorf("topic cannot be empty")
	}

	log.Printf("Creating Kafka producer...")
	log.Printf("	Brokers: %v", brokers)
	log.Printf("	Topic: %s", topic)

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	log.Println("Kafka producer created successfully")

	client := &KafkaClient{
		producer: producer,
		topic:    topic,
	}

	client.startBackgroundHandlers()

	return client, nil
}

func (kc *KafkaClient) startBackgroundHandlers() {
	kc.wg.Add(1)
	go func() {
		defer kc.wg.Done()
		for range kc.producer.Successes() {
			atomic.AddInt64(&kc.successCount, 1)
		}
	}()

	kc.wg.Add(1)
	go func() {
		defer kc.wg.Done()
		for err := range kc.producer.Errors() {
			atomic.AddInt64(&kc.errorCount, 1)
			log.Printf("Failed to send message: %v (topic: %s, partition: %d, offset: %d", err.Err, err.Msg.Topic, err.Msg.Partition, err.Msg.Offset)
		}
	}()
}

func (kc *KafkaClient) SendMessage() error {
	msg := &sarama.ProducerMessage{
		Topic: kc.topic,
		Value: sarama.ByteEncoder("Hello, World!"),
	}

	kc.producer.Input() <- msg

	return nil
}

func (kc *KafkaClient) GetStats() (success int64, errors int64) {
	return atomic.LoadInt64(&kc.successCount), atomic.LoadInt64(&kc.errorCount)
}

func (kc *KafkaClient) Close() error {
	log.Println("Closing Kafka producer...")

	if err := kc.producer.Close(); err != nil {
		log.Printf("Error closing producer: %v", err)
		return err
	}

	kc.wg.Wait()

	success, errors := kc.GetStats()
	log.Printf("Kafka producer closed. Final stats: Successes=%d, Errors=%d", success, errors)

	return nil
}
