package main

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"net/smtp"
	"encoding/json"
)
type smtpServer struct {
	host string
	port string
   }
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
   }
var (
	kafkaBrokers = []string{"localhost:9093"}
	kafkaTopics = []string{"sarama_topic"}
	consumerGroupID =  "sarama_consumer"
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
type ConsumerGroupHandler struct{}


func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }


func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }


func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
			m := LOGIN{}
	json.Unmarshal(msg.Value, &m)
		fmt.Printf("Message topic:%q partition:%d offset:%d message: %v\n",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		sess.MarkMessage(msg, "")
		from := "yourmail@gmail.com"
    password := "yourpassword"
    
    to := []string{
        string(m.Email),
    }
   
    smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}
    
    message := []byte("successful")
    
    auth := smtp.PlainAuth("", from, password, smtpServer.host)
    
    err := smtp.SendMail(smtpServer.Address(), auth, from, to, message)
    if err != nil {
        fmt.Println(err)
        return nil
    }
    fmt.Println("Email Sent!")

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
