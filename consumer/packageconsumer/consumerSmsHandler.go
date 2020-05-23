package packageconsumer

import (
	"PubSub/config/getEnvVars"
	"PubSub/config/initConfig"
	// "PubSub/sender"
)

// var wg = sync.WaitGroup{}

// StartConsumerSms asdf
func StartConsumerSms() {
	getEnvVars.GetEnvVars()
	initConfig.InitConfig()
	wg.Add(2)
	go Sender("SMS_Group")
	wg.Wait()
}
