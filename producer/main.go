package main

import (
	"PubSub/config/initConfig"
	"PubSub/logger"
	"PubSub/producer/packageproducer"
	"PubSub/routes"

	"fmt"

	"github.com/spf13/viper"
)

var (
	listenAddrAPI  string
	kafkaBrokerURL string
	kafkaVerbose   bool
	kafkaClientID  string
	kafkaTopic     string
	postAddr       string
)

//StartProducer asdf
func main() {
	initConfig.InitConfig()
	fmt.Println("Producer is running, Happy Posting!....")
	listenAddrAPI = viper.GetString("server")
	postAddr = viper.GetString("postAddr")
	router := routes.Route(listenAddrAPI, postAddr, packageproducer.ReceiveData)
	err := router.Run(listenAddrAPI)

	if err != nil {
		logger.SugarLogger.Error("Could not stablish a server, exiting...")
	}
}
