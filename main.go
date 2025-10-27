package main

import (
	"log"
)

func main() {
	log.Printf("Starting application...")

	config := NewConfig()

	kafka, err := NewKafkaClient(
		config.KafkaBrokers,
		config.KafkaTopic,
		config.ProducerConfig,
	)
	if err != nil {
		log.Fatalf("Failed to create Kafka client: %v", err)
	}
	defer kafka.Close()

	pool := NewWorkerPool(10000)
	for i := range 10000 {
		jobNum := i
		pool.Submit(func() {
			if err := kafka.SendMessage(); err != nil {
				log.Printf("Job %d: failed to send message: %v", jobNum, err)
			}
		})
	}
	pool.Close()
	log.Println("All jobs submitted and workers finished")
}
