package util

import (
	gomail "gopkg.in/mail.v2"
)

func SendEmail(receiver string, subject string, body string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", "wattflow12@gmail.com")
	message.SetHeader("To", receiver)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	dialer := gomail.NewDialer("smtp.gmail.com", 587, "wattflow12@gmail.com", "qcwb rfyi hiyi ehxo")

	if err := dialer.DialAndSend(message); err != nil {
		return err
	} else {
		return nil
	}
}
