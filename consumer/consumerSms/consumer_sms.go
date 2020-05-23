package consumerSms

import (
	"PubSub/config/getEnvVars"
	"PubSub/config/initConfig"
	"PubSub/sender"

	"sync"
)

var wg = sync.WaitGroup{}

func StartConsumerSms() {
	getEnvVars.GetEnvVars()
	initConfig.InitConfig()
	wg.Add(2)
	go sender.Sender("SMS_Group")
	wg.Wait()
}
