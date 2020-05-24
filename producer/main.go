package main

import (
	"producer/internal/app/config"
	"producer/internal/app/handler"
	"producer/logger"
	"producer/routes"

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
	config.InitConfig()
	fmt.Println("Producer is running, Happy Posting!....")
	listenAddrAPI = viper.GetString("server")
	postAddr = viper.GetString("postAddr")
	router := routes.Route(listenAddrAPI, postAddr, handler.ReceiveData)
	err := router.Run(listenAddrAPI)

	if err != nil {
		logger.SugarLogger.Error("Could not stablish a server, exiting...")
	}
}
