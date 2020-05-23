package packageproducer

import (
	"PubSub/logger"
	"encoding/json"
	"net/http"

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

// PostDataToKafka temp
func ReceiveData(c *gin.Context) {
	c.Bind(&someData)
	c.JSON(http.StatusOK, &someData)
	formInBytes, err := json.Marshal(&someData)
	if err != nil {
		logger.SugarLogger.Error("Error occured white marshalling data", err)
	}

	PostDataToKafka(formInBytes)
}
