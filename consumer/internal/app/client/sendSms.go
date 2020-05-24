package client

import (
	"consumer/logger"
	"os"

	"github.com/sfreiberg/gotwilio"
)

//SendSms sadf
func SendSms(sendMe string, whatMessage string) error {
	accountSid := os.Getenv("smsAccountSid")
	authToken := os.Getenv("smsAuthToken")
	twilio := gotwilio.NewTwilioClient(accountSid, authToken)
	from := os.Getenv("smsFrom")
	to := sendMe
	message := whatMessage
	_, _, err := twilio.SendSMS(from, to, message, "", "")
	if err != nil {
		logger.SugarLogger.Error("Can't send Message. Error occured:", err)
		return err
	}
	return nil
}

// SendMail asdfasdf
// func SendMail(semdMe string, whatMessage string) error {
// 	m := gomail.NewMessage()
// 	m.SetHeader("From", viper.GetString("emailFrom"))
// 	m.SetHeader("To", semdMe)
// 	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
// 	m.SetHeader("Subject", "Transaction Status!")
// 	m.SetBody("text/html", whatMessage)
// 	// m.Attach("/home/Alex/kiaraAdvani.jpg")
// 	d := gomail.NewDialer(viper.GetString("emailServerAddr"), viper.GetInt("emailServerPort"), viper.GetString("emailFrom"), os.Getenv("EMAIL_PASS"))
// 	// Send the email to Bob, Cora and Dan.
// 	if err := d.DialAndSend(m); err != nil {
// 		// panic(err)
// 		logger.SugarLogger.Error("Can't send Email. Error occured", err)
// 		return err

// 	}
// 	return nil

// }
