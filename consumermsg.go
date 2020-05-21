package main

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"encoding/json"
	"net/http"
    "github.com/nexmo-community/nexmo-go"
)

var (
	kafkaBrokers = []string{"localhost:9093"}
	kafkaTopics = []string{"sarama_topic"}
	consumerGroupID =  "sarama_consumer2"
)
type LOGIN struct {
	Request_id string `json:"rid" binding:"required"`
	Topic_name string `json:"topic_name" binding:"required"`
	Message_body string `json:"message_body" binding:"required"`
	Transaction_id string `json:"transction_id" binding:"required"` 
	Email string `json:"email" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Customer_id string `json:"customer_id" binding:"required"`
	Key string `json:"key" binding:"required"`
}
type ConsumerGroupHandler struct{}


func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }


func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }


func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		m := LOGIN{}
json.Unmarshal(msg.Value, &m) // json object to struct
		fmt.Printf("Message topic:%q partition:%d offset:%d message: %v\n",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		sess.MarkMessage(msg, "")

	
	auth := nexmo.NewAuthSet()
    auth.SetAPISecret("aaf3f540", "ANhXN0UaGiTaDlMZ")
 
   
    client := nexmo.NewClient(http.DefaultClient, auth)
 
    
    smsContent := nexmo.SendSMSRequest{
    From: "Vonage SMS API",
    To:   m.Phone,
    Text: "Transaction Successful",
  }
 
    smsResponse, _, err := client.SMS.SendSMS(smsContent)
 
    if err != nil {
        log.Fatal(err)
    }
	if smsResponse.Messages[0].Status =="0" {
	fmt.Println("Status:", "Successful")
	}
	
	
}
return nil
}
func startConsumer() {
	
	config := sarama.NewConfig()
	sarama.Logger = log.New(os.Stderr, "[sarama_logger]", log.LstdFlags)
	config.Version = sarama.V2_1_0_0

	
	client, err := sarama.NewClient(kafkaBrokers, config)
	if err != nil {
		panic(err)
	}
	defer func() { _ = client.Close() }()

	
	group, err := sarama.NewConsumerGroupFromClient(consumerGroupID, client)
	if err != nil {
		panic(err)
	}
	defer func() { _ = group.Close() }()
	log.Println("Consumer up and running")

	
	go func() {
		for err := range group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	
	ctx := context.Background()
	for {
		handler := ConsumerGroupHandler{}

		err := group.Consume(ctx, kafkaTopics, handler)
		if err != nil {
			panic(err)
		}
	}
}
func main() {
	startConsumer()
}
