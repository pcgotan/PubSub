package main

import (
	"PubSub/config"
	"PubSub/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/namsral/flag"
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

func main() {
	service := "PubSub"
	environment := os.Getenv("BOOT_CUR_ENV")
	if environment == "" {
		environment = "test"
	}
	flag.Usage = func() {
		fmt.Println("Usage: server -s {service_name} -e {environment}")
		os.Exit(1)
	}
	flag.Parse()
	configURL := "" // Put the configuration url of spring cloud config
	config.Init(configURL, service, environment)
	logger.InitLogger()
	flag.Parse()
	fmt.Println("Producer is running, Happy Posting!....")
	listenAddrAPI = viper.GetString("server")
	postAddr = viper.GetString("postAddr")

	err := server(listenAddrAPI, postAddr)
	if err != nil {
		logger.SugarLogger.Error("Could not stablish a server, exiting...")
	}

}
func server(listenAddr, postAddr string) error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.POST(postAddr, postDataToKafka)
	for _, routeInfo := range router.Routes() {
		logger.SugarLogger.Debug("path: ", routeInfo.Path, "\thandler: ", routeInfo.Handler, "\tmethod: ", routeInfo.Method, ",\tregistered routes")

	}
	return router.Run(listenAddr)
}

func postDataToKafka(c *gin.Context) {
	c.Bind(&someData)
	c.JSON(http.StatusOK, &someData)
	formInBytes, err := json.Marshal(&someData)
	if err != nil {
		logger.SugarLogger.Error("Error occured white marshalling data", err)
	}

	kafkaBrokerURL = viper.GetString("brokers")
	kafkaClientID = viper.GetString("clientID")
	kafkaVerbose = viper.GetBool("verbose")
	kafkaTopic = viper.GetString("topic")
	if (someData.Topic_name) != "" {
		viper.Set("topic", (someData.Topic_name))
		kafkaTopic = viper.GetString("topic")
	}

	dialForTopicCreation, _ := kafka.Dial("tcp", strings.Split(kafkaBrokerURL, ",")[0])
	leaderBroker, _ := dialForTopicCreation.Controller()
	leaderAddr := "localhost:" + strconv.Itoa(leaderBroker.Port)
	dialForTopicCreation, _ = kafka.Dial("tcp", leaderAddr)
	newTopicConfig := kafka.TopicConfig{Topic: kafkaTopic, NumPartitions: 10, ReplicationFactor: 3}
	err = dialForTopicCreation.CreateTopics(newTopicConfig)
	if err != nil {
		logger.SugarLogger.Error("Error white creating topic, dial to the leader", err)
	}

	kafkaProducer, _ = config.Configure(strings.Split(kafkaBrokerURL, ","), kafkaClientID, kafkaTopic, formInBytes)
	defer kafkaProducer.Close()
	parent := context.Background()
	defer parent.Done()
	err = push(parent, someData.Key, formInBytes)
	if err != nil {
		logger.SugarLogger.Error("error while pushing message into kafka: %s", err.Error())
	}
	// fmt.Println(viper.GetInt("temppp"))
	// viper.Set("temppp", viper.GetInt("temppp")+1)
}

func push(parent context.Context, key string, value []byte) (err error) {
	message := kafka.Message{
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	}
	return kafkaProducer.WriteMessages(parent, message)

}
