package client

import (
	"consumer/logger"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

// SendMail asdfasdf
func SendMail(semdMe string, whatMessage string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", viper.GetString("emailFrom"))
	m.SetHeader("To", semdMe)
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Transaction Status!")
	m.SetBody("text/html", whatMessage)
	// m.Attach("/home/Alex/kiaraAdvani.jpg")
	d := gomail.NewDialer(viper.GetString("emailServerAddr"), viper.GetInt("emailServerPort"), viper.GetString("emailFrom"), os.Getenv("EMAIL_PASS"))
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		// panic(err)
		logger.SugarLogger.Error("Can't send Email. Error occured", err)
		return err

	}
	return nil

}
