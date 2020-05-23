package main

import (
	"PubSub/consumer/packageconsumer"
	// "PubSub/consumer/consumerSms"
	"fmt"
	"os"
)

func main() {
	if os.Args[1] == "mail" {
		packageconsumer.StartConsumerMail()
	} else if os.Args[1] == "sms" {
		packageconsumer.StartConsumerSms()
	} else {
		fmt.Println("Command not found")
	}
}
