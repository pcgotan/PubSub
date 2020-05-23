package packageproducer

import (
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
)

var writer *kafka.Writer

//Configure asdf
func Configure(kafkaBrokerUrls []string, clientID string, topic string, formInBytes []byte, balancer1 kafka.Balancer) (w *kafka.Writer, err error) {
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: clientID,
	}
	if err := json.Unmarshal(formInBytes, &someData); err != nil {
		panic(err)
	}
	config1 := kafka.WriterConfig{
		Brokers:      kafkaBrokerUrls,
		Topic:        someData.Topic_name,
		Balancer:     balancer1,
		Dialer:       dialer,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	w = kafka.NewWriter(config1)
	writer = w
	return w, nil
}
