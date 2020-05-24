package handler

import (
	"encoding/json"
	"net/http"
	"producer/internal/app/service"
	"producer/internal/app/structs"
	"producer/logger"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

var (
	listenAddrAPI  string
	kafkaBrokerURL string
	kafkaVerbose   bool
	kafkaClientID  string
	kafkaTopic     string
	postAddr       string
)

var someData structs.WholeData
var kafkaProducer *kafka.Writer

// ReceiveData temp
func ReceiveData(c *gin.Context) {
	c.Bind(&someData)
	c.JSON(http.StatusOK, &someData)
	formInBytes, err := json.Marshal(&someData)
	if err != nil {
		logger.SugarLogger.Error("Error occured white marshalling data", err)
	}

	service.PostDataToKafka(formInBytes, someData)
}
