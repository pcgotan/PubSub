package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"producer/logger"

	"github.com/spf13/viper"
)

// type wholeData struct {
// 	Request_Id     string `form:"request_id" json:"request_id"`
// 	Topic_name     string `form:"topic_name" json:"topic_name"`
// 	Message_body   string `form:"message_body" json:"message_body"`
// 	Transaction_id string `form:"transaction_id" json:"transaction_id"`
// 	Email          string `form:"email" json:"email"`
// 	Phone          string `form:"phone" json:"phone"`
// 	Customer_id    string `form:"customer_id" json:"customer_id"`
// 	Key            string `form:"key" json:"key"`
// }

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

// var writer *kafka.Writer
// var someData wholeData

// Init :
func Init(configurationURL, service, env string) {
	url := configurationURL + service + "/" + env
	// fmt.Println("url is : ", url)
	// fmt.Println("Loading config from \n", url)
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
		bodyBytes, err = ioutil.ReadFile("./internal/app/config/config.json")
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
	fmt.Println("Successfully loaded all configurations")
	if viper.IsSet("server_name") {
		fmt.Println("Successfully loaded configuration for service\n", viper.GetString("server_name"))
	}
}

// InitConfig :
func InitConfig() {
	service := "producer"
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
	Init(configURL, service, environment)
	logger.InitLogger()
}
