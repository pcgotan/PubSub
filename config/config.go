package config

import (
	"PubSub/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
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

type springCloudConfig struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           string           `json:"label"`
	Version         string           `json:"version"`
	PropertySources []propertySource `json:"propertySources"`
}
type propertySource struct {
	Name   string                 `json:"name"`
	Source map[string]interface{} `json:"source"`
}

var config *viper.Viper
var writer *kafka.Writer
var someData wholeData

// INIT
func Init(configurationURL, service, env string) {
	url := configurationURL + service + "/" + env
	fmt.Println("url is : ", url)
	fmt.Println("Loading config from \n", url)
	body, err := fetchConfiguration(url)
	if err != nil {
		fmt.Println("Couldn't load configuration, cannot start. Terminating. Error: " + err.Error())
	}
	parseConfiguration(body)
}

func fetchConfiguration(url string) ([]byte, error) {
	resp, err := http.Get(url)
	var bodyBytes []byte
	if err != nil {
		//panic("Couldn't load configuration, cannot start. Terminating. Error: " + err.Error())
		bodyBytes, err = ioutil.ReadFile("src/PubSub/config/config.json")
		if err != nil {
			fmt.Println("Couldn't read local configuration file.", err)
		} else {
			log.Print("using local config.")
		}
	} else {
		if resp != nil {
			bodyBytes, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading configuration response body.")
			}
		}
	}
	return bodyBytes, err
}

func parseConfiguration(body []byte) {
	var cloudConfig springCloudConfig
	err := json.Unmarshal(body, &cloudConfig)
	if err != nil {
		fmt.Println("Cannot parse configuration, message: " + err.Error())
	}
	for key, value := range cloudConfig.PropertySources[0].Source {
		viper.Set(key, value)
		// fmt.Println("Loading config property\n", key, value)
	}
	fmt.Println("Successfully loaded all configurations\n")
	if viper.IsSet("server_name") {
		fmt.Println("Successfully loaded configuration for service\n", viper.GetString("server_name"))
	}
}

func Configure(kafkaBrokerUrls []string, clientId string, topic string, formInBytes []byte) (w *kafka.Writer, err error) {
	dialer := &kafka.Dialer{
		Timeout:  10 * time.Second,
		ClientID: clientId,
	}
	if err := json.Unmarshal(formInBytes, &someData); err != nil {
		panic(err)
	}

	config_partition := kafka.WriterConfig{
		Brokers: kafkaBrokerUrls,
		Topic:   someData.Topic_name,
		Balancer: kafka.BalancerFunc(func(msg kafka.Message, partitions ...int) int {
			i, _ := strconv.ParseInt(someData.Request_Id, 10, 32)
			if int(i) >= len(partitions) {
				logger.SugarLogger.Info("Specified partition is greater than total number of partition, Writing it to (specified partition) mod (total partition)")
			}
			return int(int(i) % (len(partitions)))
		}),
		Dialer:       dialer,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	config_key := kafka.WriterConfig{
		Brokers:      kafkaBrokerUrls,
		Topic:        someData.Topic_name,
		Balancer:     &kafka.Hash{},
		Dialer:       dialer,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	config_leastbytes := kafka.WriterConfig{
		Brokers:      kafkaBrokerUrls,
		Topic:        someData.Topic_name,
		Balancer:     &kafka.LeastBytes{},
		Dialer:       dialer,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	if someData.Key == "" && someData.Request_Id != "" {
		logger.SugarLogger.Info("Using partition specific balancer")
		w = kafka.NewWriter(config_partition)
	}
	if someData.Key != "" {
		logger.SugarLogger.Info("Using Hash Key balancer")
		w = kafka.NewWriter(config_key)
	}
	if someData.Key == "" && someData.Request_Id == "" {
		logger.SugarLogger.Info("Using LeastBytes balancer")
		w = kafka.NewWriter(config_leastbytes)
	}

	writer = w
	return w, nil
}
