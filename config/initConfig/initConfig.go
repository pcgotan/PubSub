package initConfig

import (
	"PubSub/config"
	"PubSub/logger"
	"flag"
	"fmt"
	"os"
)

func InitConfig() {
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
}
