package events

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	kafkaBrokers = "kafka:9093" // localhost if run locally
	topic        = "sarama_topic"
	groupID      = "my-group"
)

func KafkaConsumer() {
	// make a new reader that consumes from topic-A, partition 0, at offset 42
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{kafkaBrokers},
		Topic:     topic,
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
	r.SetOffset(42)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}

func KafkaProducer() {
	producer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaBrokers},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	// Produce messages
	for i := 0; i < 10; i++ {
		message := kafka.Message{
			Value: []byte(fmt.Sprintf("Message %d", i)),
		}
		err := producer.WriteMessages(context.Background(), message)
		if err != nil {
			log.Fatal("Failed to write message:", err)
		}
		fmt.Println("Message sent successfully ", string(message.Value))
		time.Sleep(time.Second) // Optional: Add a delay between message sends
	}

	// Close the producer when done
	producer.Close()

}
