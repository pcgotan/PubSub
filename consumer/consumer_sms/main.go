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
	"gopkg.in/gomail.v2"
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
		// sendMe := raw["phone"]
		sendMe := raw["email"]
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
			go sendMail(sendMe.(string), whatMessage.(string))
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

func sendMail(semdMe string, whatMessage string) {
	// temp := "<head><link href=\"//netdna.bootstrapcdn.com/bootstrap/3.0.3/css/bootstrap.min.css\" rel=\"stylesheet\" id=\"bootstrap-css\"><script src=\"//netdna.bootstrapcdn.com/bootstrap/3.0.3/js/bootstrap.min.js\"></script><script src=\"//code.jquery.com/jquery-1.11.1.min.js\"></script></head>  <div class=\"alert alert-success\"  style=\"text-align:center\">    <span class=\"glyphicon glyphicon-ok\"></span> <strong>Congratulations</strong>    <hr class=\"message-inner-separator\">    <p>      Transaction Successful.</p>  </div>"
	m := gomail.NewMessage()
	m.SetHeader("From", viper.GetString("emailFrom"))
	m.SetHeader("To", semdMe)
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Transaction Status")
	m.SetBody("text/html", whatMessage)
	// m.Attach("/home/Alex/kiaraAdvani.jpg")
	d := gomail.NewDialer(viper.GetString("emailServerAddr"), viper.GetInt("emailServerPort"), viper.GetString("emailFrom"), os.Getenv("EMAIL_PASS"))
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		// panic(err)
		logger.SugarLogger.Error("Can't send Email. Error occured")
	}
	wg.Done()
}
