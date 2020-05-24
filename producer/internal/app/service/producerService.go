package service

import (
	"context"
	"producer/internal/app/kafka"
	"producer/internal/app/structs"
	"producer/logger"
	"strconv"
	"strings"

	gokafka "github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

// PostDataToKafka temp
func PostDataToKafka(formInBytes []byte, someData structs.WholeData) {

	kafkaBrokerURL := viper.GetString("brokers")
	kafkaClientID := viper.GetString("clientID")
	//	kafkaVerbose := viper.GetBool("verbose")
	kafkaTopic := viper.GetString("topic")
	if (someData.Topic_name) != "" {
		viper.Set("topic", (someData.Topic_name))
		kafkaTopic = viper.GetString("topic")
	}

	dialForTopicCreation, _ := gokafka.Dial("tcp", strings.Split(kafkaBrokerURL, ",")[0])
	leaderBroker, _ := dialForTopicCreation.Controller()
	leaderAddr := "localhost:" + strconv.Itoa(leaderBroker.Port)
	// leaderAddr := strings.Split(kafkaBrokerURL, ":")[0] + ":" + strconv.Itoa(leaderBroker.Port)
	dialForTopicCreation, _ = gokafka.Dial("tcp", leaderAddr)
	newTopicConfig := gokafka.TopicConfig{Topic: kafkaTopic, NumPartitions: 10, ReplicationFactor: 3}
	err := dialForTopicCreation.CreateTopics(newTopicConfig)
	if err != nil {
		logger.SugarLogger.Error("Error while creating topic, dial to the leader", err)
	}

	var balancer1 gokafka.Balancer

	if someData.Key == "" && someData.Request_Id != "" {

		balancer1 = gokafka.BalancerFunc(func(msg gokafka.Message, partitions ...int) int {
			i, _ := strconv.ParseInt(someData.Request_Id, 10, 32)
			if int(i) >= len(partitions) {
				logger.SugarLogger.Info("Specified partition is greater than total number of partition, Writing it to (specified partition) mod (total partition)")
			}
			return int(int(i) % (len(partitions)))
		})
	}

	if someData.Key != "" {
		balancer1 = &gokafka.Hash{}
	}

	if someData.Key == "" && someData.Request_Id == "" {
		balancer1 = &gokafka.LeastBytes{}
	}

	kafkaProducer, _ := kafka.Configure(strings.Split(kafkaBrokerURL, ","), kafkaClientID, kafkaTopic, formInBytes, balancer1, someData)
	defer kafkaProducer.Close()
	parent := context.Background()
	defer parent.Done()
	err = kafka.Push(parent, someData.Key, formInBytes)
	if err != nil {
		logger.SugarLogger.Error("error while pushing message into kafka: %s", err.Error())
	}
}
