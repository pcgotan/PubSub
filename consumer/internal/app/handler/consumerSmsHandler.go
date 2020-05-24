package handler

import (
	"consumer/internal/app/config"
	"consumer/internal/app/config/getEnvVars"
	"consumer/internal/app/service"
	"fmt"
	"sync"
)

// StartConsumerSms asdf
func StartConsumerSms() {
	getEnvVars.GetEnvVars()
	config.InitConfig()
	fmt.Println("Consumer SmsGroup is running, Happy Consuming!....")
	var wg = sync.WaitGroup{}
	wg.Add(2)
	go service.Sender("SMS_Group")
	wg.Wait()
}
