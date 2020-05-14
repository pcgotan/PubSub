package main

import (
	"PubSub/config"
	"PubSub/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/namsral/flag"
	kafka "github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

// var logger = log.With().Str("pkg", "main").Logger()
var (
	listenAddrAPI  string
	kafkaBrokerURL string
	kafkaVerbose   bool
	kafkaClientID  string
	kafkaTopic     string
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
	kafkaBrokerURL = viper.GetString("brokers")
	kafkaClientID = viper.GetString("clientID")
	kafkaVerbose = viper.GetBool("verbose")
	kafkaTopic = viper.GetString("topic")

	err := server(listenAddrAPI)
	if err != nil {
		logger.SugarLogger.Error("Could not stablish a server, exiting...")
	}

	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalChan:
		logger.SugarLogger.Error("got an interrupt, exiting...")
	}
}

func server(listenAddr string) error {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.POST(viper.GetString("postAddr"), postDataToKafka)
	for _, routeInfo := range router.Routes() {
		logger.SugarLogger.Debug("path: ", routeInfo.Path, "		handler: ", routeInfo.Handler, "	method: ", routeInfo.Method, ",   registered routes")

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
	kafkaProducer, _ = config.Configure(strings.Split(kafkaBrokerURL, ","), kafkaClientID, kafkaTopic, formInBytes)
	defer kafkaProducer.Close()
	parent := context.Background()
	defer parent.Done()
	err = push(parent, someData.Key, formInBytes)
	if err != nil {
		logger.SugarLogger.Error("error while pushing message into kafka: %s", err.Error())
	}
}

func push(parent context.Context, key string, value []byte) (err error) {
	message := kafka.Message{
		Key:   []byte(key),
		Value: value,
		Time:  time.Now(),
	}
	return kafkaProducer.WriteMessages(parent, message)
}
