package consumerMail

import (
	"PubSub/config/getEnvVars"
	"PubSub/config/initConfig"
	"PubSub/sender"
	"sync"
)

var wg = sync.WaitGroup{}

func StartConsumerMail() {
	getEnvVars.GetEnvVars()
	initConfig.InitConfig()
	wg.Add(2)
	go sender.Sender("Email_Group")
	wg.Wait()
}
