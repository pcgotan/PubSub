package kafka

import (
	"context"
	"time"

	gokafka "github.com/segmentio/kafka-go"
)

//Push asdf
func Push(parent context.Context, key string, value []byte) (err error) {
	message := gokafka.Message{
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	}
	return Writer.WriteMessages(parent, message)
}
