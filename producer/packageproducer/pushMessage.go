package packageproducer

import (
	"context"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

//Push asdf
func Push(parent context.Context, key string, value []byte) (err error) {
	message := kafka.Message{
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	}
	return kafkaProducer.WriteMessages(parent, message)
}
