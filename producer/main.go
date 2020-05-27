package main

import (
	"producer/internal/app/config"
	"producer/internal/app/handler"
	"producer/internal/app/service"
	"producer/logger"
	"producer/routes"
	"strings"

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
	kafkaBrokerURL := viper.GetString("brokers")
	topics1 := viper.GetString("topics")
	topics := strings.Split(topics1, ",")
	for _, topic := range topics {
		service.CreateTopic(kafkaBrokerURL, topic)
	}

	postAddr = viper.GetString("postAddr")
	router := routes.Route(listenAddrAPI, postAddr, handler.ReceiveData)
	err := router.Run(listenAddrAPI)

	if err != nil {
		logger.SugarLogger.Error("Could not stablish a server, exiting...")
	}
}
