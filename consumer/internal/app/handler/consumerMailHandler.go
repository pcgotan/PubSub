package handler

import (
	"consumer/internal/app/config"
	"consumer/internal/app/config/getEnvVars"
	"consumer/internal/app/service"
	"fmt"
	"sync"
)

//StartConsumerMail sadf
func StartConsumerMail() {
	getEnvVars.GetEnvVars()
	config.InitConfig()
	fmt.Println("Consumer MailGroup is running, Happy Consuming!....")
	var wg = sync.WaitGroup{}
	wg.Add(2)
	go service.Sender("Email_Group")
	wg.Wait()
}
