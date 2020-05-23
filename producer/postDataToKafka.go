package producer

import (
	"PubSub/config"
	"PubSub/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	kafka "github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

// PostDataToKafka temp
func PostDataToKafka(c *gin.Context) {
	c.Bind(&someData)
	c.JSON(http.StatusOK, &someData)
	formInBytes, err := json.Marshal(&someData)
	fmt.Printf("%T", formInBytes)
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
	// leaderAddr := strings.Split(kafkaBrokerURL, ":")[0] + ":" + strconv.Itoa(leaderBroker.Port)
	dialForTopicCreation, _ = kafka.Dial("tcp", leaderAddr)
	newTopicConfig := kafka.TopicConfig{Topic: kafkaTopic, NumPartitions: 10, ReplicationFactor: 3}
	err = dialForTopicCreation.CreateTopics(newTopicConfig)
	if err != nil {
		logger.SugarLogger.Error("Error while creating topic, dial to the leader", err)
	}

	kafkaProducer, _ = config.Configure(strings.Split(kafkaBrokerURL, ","), kafkaClientID, kafkaTopic, formInBytes)
	defer kafkaProducer.Close()
	parent := context.Background()
	defer parent.Done()
	err = Push(parent, someData.Key, formInBytes)
	if err != nil {
		logger.SugarLogger.Error("error while pushing message into kafka: %s", err.Error())
	}
}
