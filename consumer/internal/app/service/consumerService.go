package service

import (
	"consumer/internal/app/client"
	"consumer/logger"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/fatih/color"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

var wg = sync.WaitGroup{}

//Sender sdf
func Sender(clientID string) {
	topics1 := viper.GetString("topics")
	topics := strings.Split(topics1, ",")
	kafkaClientID := clientID
	allBrokers := viper.GetString("brokers")
	brokers := strings.Split(allBrokers, ",")
	group, err := kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
		ID:      kafkaClientID,
		Brokers: brokers,
		Topics:  topics,
	})

	if err != nil {
		logger.SugarLogger.Error("error creating consumer group:", err)

		os.Exit(1)
	}
	defer group.Close()
	for {
		gen, err := group.Next(context.TODO())
		if err != nil {
			break
		}
		for _, s := range topics {
			topic := s
			// fmt.Println(topic)
			assignments := gen.Assignments[topic]
			for _, assignment := range assignments {
				partition, offset := assignment.ID, assignment.Offset
				gen.Start(func(ctx context.Context) {
					reader := kafka.NewReader(kafka.ReaderConfig{
						Brokers:         brokers,
						Topic:           topic,
						Partition:       partition,
						MinBytes:        viper.GetInt("readerMinBytes"),
						MaxBytes:        viper.GetInt("readerMaxBytes"),
						MaxWait:         1 * time.Second,
						ReadLagInterval: -1,
					})
					defer reader.Close()
					//last committed offset for this partition + 1 (start consuming from this offset).
					reader.SetOffset(offset + 1)
					tempp := offset // this variavle will track the max committed offset
					for {
						wg.Add(1)
						m, err := reader.ReadMessage(ctx)
						switch err {
						case kafka.ErrGenerationEnded:
							gen.CommitOffsets(map[string]map[int]int64{topic: {partition: offset}})
							return
						case nil:
							value := m.Value
							var raw map[string]interface{}
							json.Unmarshal(value, &raw)
							var sendMe interface{}
							if clientID == "Email_Group" {
								sendMe = raw["email"]
							}
							if clientID == "SMS_Group" {
								sendMe = raw["phone"]
								// sendMe = raw["email"]
							}
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
								go func() {
									var err error
									if clientID == "Email_Group" {
										err = client.SendMail(sendMe.(string), whatMessage.(string))
									}
									if clientID == "SMS_Group" {
										err = client.SendSms(sendMe.(string), whatMessage.(string))
									}
									if err == nil {
										offset = m.Offset
										if offset > tempp {
											tempp = offset
											gen.CommitOffsets(map[string]map[int]int64{topic: {partition: offset}})
											// fmt.Println("----------------------------------Committed offset till--------------   ", offset)

										}
										// fmt.Println("----------------------------------Mail sent , Offset-------------->  ", offset)

									}
									wg.Done()
								}()
							} else {
								logger.SugarLogger.Error("Email Address is empty")
							}
						default:
							logger.SugarLogger.Error("error reading message: ", err)
						}
					}
					wg.Wait()
				})
			}

		}

	}

}
