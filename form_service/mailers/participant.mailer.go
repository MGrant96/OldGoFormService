package mailers

import (
	"fmt"
	"github.com/joho/godotenv"
	// "net/smtp"
	"os"
	"gopkg.in/gomail.v2"
)

func SendEmail(participantID string, recipient string, incomingMessage string, subject string) {

	error := godotenv.Load()
	fmt.Println("Error", error)

	// Sender data.
	from := os.Getenv("USER_NAME")
	password := os.Getenv("PASSWORD")

	// smtp server configuration.
	smtpHost := os.Getenv("HOST")
	smtpPort := os.Getenv("SMPTP_PORT")

	message := (
		"Subject: " + subject + "\r\n" +
		"\r\n" + incomingMessage +
		"\r\nPlease follow the below link to provide your feedback:" +
		// "\r\n" + "http://staging.polyloop.io.s3-website.eu-west-1.amazonaws.com/feedback/" + participantID +
		"\r\n" + "http://localhost:3001/feedback/" + participantID +
		"\r\n")

	mail := gomail.NewMessage()
	mail.SetHeader("From", from)
	mail.SetHeader("To", recipient)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/html", message)

	fmt.Println("PORTS: ", smtpPort)

	dispatch := gomail.NewDialer(smtpHost, 587, from, password)

	// Message.

	// Authentication.
	// auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	// err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, recipients, message)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	if err := dispatch.DialAndSend(mail); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Email Sent Successfully!")
	}
}
