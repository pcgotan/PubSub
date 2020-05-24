package main

import (
	"consumer/internal/app/handler"
	"fmt"
	"os"
)

func main() {
	if os.Args[1] == "mail" {
		handler.StartConsumerMail()
	} else if os.Args[1] == "sms" {
		handler.StartConsumerSms()
	} else {
		fmt.Println("Command not found")
	}
}
