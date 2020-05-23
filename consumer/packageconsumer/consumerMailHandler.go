package packageconsumer

import (
	"PubSub/config/getEnvVars"
	"PubSub/config/initConfig"
	// "PubSub/sender"
)

// var wg = sync.WaitGroup{}

//StartConsumerMail sadf
func StartConsumerMail() {
	getEnvVars.GetEnvVars()
	initConfig.InitConfig()
	wg.Add(2)
	go Sender("Email_Group")
	wg.Wait()
}
