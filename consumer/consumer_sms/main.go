package main

import (
	"PubSub/config"
	"PubSub/logger"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	kafka "github.com/segmentio/kafka-go"
	"github.com/sfreiberg/gotwilio"
	"github.com/spf13/viper"
)

func getEnvVars() {
	err := godotenv.Load("./src/PubSub/credentials.env")
	if err != nil {
		// logger.SugarLogger.Error("Error laoding .env file")
		fmt.Println("Error laoding .env file")
	}
}

var wg = sync.WaitGroup{}

func main() {
	getEnvVars()
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

	wg.Add(2)
	go smsWale()
	wg.Wait()
}

func smsWale() {
	topic := viper.GetString("topic")
	kafkaClientID := "SMS_Group"
	allBrokers := viper.GetString("brokers")
	brokers := strings.Split(allBrokers, ",")
	config := kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		// Partition:       partition,
		GroupID:         kafkaClientID,
		MinBytes:        viper.GetInt("readerMinBytes"),
		MaxBytes:        viper.GetInt("readerMaxBytes"),
		MaxWait:         1 * time.Second,
		ReadLagInterval: -1,
	}
	reader := kafka.NewReader(config)
	defer reader.Close()
	for {
		wg.Add(1)
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			logger.SugarLogger.Error("Error Reading Messages", err)
			continue
		}
		value := m.Value
		var raw map[string]interface{}
		json.Unmarshal(value, &raw)
		sendMe := raw["phone"]
		// sendMe := raw["email"]
		whatMessage := raw["message_body"]
		f := colorjson.NewFormatter()
		f.Indent = 4
		s, _ := f.Marshal(raw)
		green := color.New(color.FgGreen).SprintFunc()

		fmt.Println(color.YellowString("\nMessage Consumed"))
		fmt.Println(color.WhiteString("Consumer_Group:"), green(kafkaClientID))
		fmt.Println(color.WhiteString("Topic:"), green(m.Topic))
		fmt.Println(color.WhiteString("Partition:"), green(m.Partition))
		fmt.Println(color.WhiteString("Offset:"), green(m.Offset))
		fmt.Println(string(s))
		fmt.Println("_______________________________________________________")
		if sendMe != nil {
			go sendSms(sendMe.(string), whatMessage.(string))
		} else {
			fmt.Println("Empty")
		}
	}
	wg.Wait()
}

func sendSms(sendMe string, whatMessage string) {
	accountSid := os.Getenv("smsAccountSid")
	authToken := os.Getenv("smsAuthToken")
	twilio := gotwilio.NewTwilioClient(accountSid, authToken)
	from := os.Getenv("smsFrom")
	to := sendMe
	message := whatMessage
	_, _, err := twilio.SendSMS(from, to, message, "", "")
	if err != nil {
		logger.SugarLogger.Error("Can't send Message. Error occured:", err)
	}
	wg.Done()
}
