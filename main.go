package main

import (
	"log"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
)

var (
	kafkaBrokers = []string{"localhost:9093"}
	KafkaTopic   = "sarama_topic"
)

type LOGIN struct {
	Request_id string `json:"request_id" binding:"required"`
	Topic_name string `json:"topic_name" binding:"required"`
	Message_body string `json:"message_body" binding:"required"`
	Transaction_id string `json:"transaction_id" binding:"required"` 
	Email string `json:"email" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Customer_id string `json:"customer_id" binding:"required"`
	Key string `json:"key" binding:"required"`
}

func main() {
	r := gin.Default()
	producer, err := setupProducer()
	if err != nil {
		panic(err)
	} else {
		log.Println("Kafka AsyncProducer up and running!")
	}

	
	r.POST("/foo", func(c *gin.Context) {
		requestBody := LOGIN{}
		c.Bind(&requestBody)
		b, err := json.Marshal(requestBody) // struct to json
		if err != nil {
			panic(err)
		} 
		c.JSON(200, gin.H{"status": requestBody.Email, "phone": requestBody.Phone})
		KafkaTopic = requestBody.Topic_name
		brokerAddrs := []string{"localhost:9093"}
		config := sarama.NewConfig()
		config.Version = sarama.V2_1_0_0
		admin, err := sarama.NewClusterAdmin(brokerAddrs, config)
		if err != nil {
			log.Fatal("Error while creating cluster admin: ", err.Error())
		}
		defer func() { _ = admin.Close() }()
		err = admin.CreateTopic(KafkaTopic, &sarama.TopicDetail{
			NumPartitions:     2,
			ReplicationFactor: 1,
		}, false)
		
		message := &sarama.ProducerMessage{Topic: KafkaTopic, Value: sarama.StringEncoder(b)}
		select {
		case producer.Input() <- message:
			log.Println("New Message produced")
		}
	})
	r.Run(":3000")

}


func setupProducer() (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	//config.Producer.Partitioner = sarama.NewManualPartitioner
	return sarama.NewAsyncProducer(kafkaBrokers, config)
}

