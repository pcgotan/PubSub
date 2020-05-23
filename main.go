package main

import (
	"PubSub/consumer/consumerMail"
	"PubSub/consumer/consumerSms"
	"PubSub/producer"
	"fmt"
	"os"
)

func main() {
	if os.Args[1] == "producer" {
		producer.StartProducer()
	} else if os.Args[1] == "consumer" && os.Args[2] == "mail" {
		consumerMail.StartConsumerMail()
	} else if os.Args[1] == "consumer" && os.Args[2] == "sms" {
		consumerSms.StartConsumerSms()
	} else {
		fmt.Println("Command not found")
	}
}
