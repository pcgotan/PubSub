package producer

import (
	"PubSub/config/initConfig"
	"PubSub/logger"

	"PubSub/server"
	"fmt"

	kafka "github.com/segmentio/kafka-go"
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

type wholeData struct {
	Request_Id     string `form:"request_id" json:"request_id"`
	Topic_name     string `form:"topic_name" json:"topic_name"`
	Message_body   string `form:"message_body" json:"message_body"`
	Transaction_id string `form:"transaction_id" json:"transaction_id"`
	Email          string `form:"email" json:"email"`
	Phone          string `form:"phone" json:"phone"`
	Customer_id    string `form:"customer_id" json:"customer_id"`
	Key            string `form:"key" json:"key"`
}

var someData wholeData
var kafkaProducer *kafka.Writer

//StartProducer asdf
func StartProducer() {
	initConfig.InitConfig()

	fmt.Println("Producer is running, Happy Posting!....")
	listenAddrAPI = viper.GetString("server")
	postAddr = viper.GetString("postAddr")
	router := server.Server(listenAddrAPI, postAddr, PostDataToKafka)
	err := router.Run(listenAddrAPI)

	if err != nil {
		logger.SugarLogger.Error("Could not stablish a server, exiting...")
	}
}
