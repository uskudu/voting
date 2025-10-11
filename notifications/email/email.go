package email

import (
	"net/smtp"
	"os"
)

func SendMail(notification map[string]string) error {

	auth := smtp.PlainAuth(
		"",
		os.Getenv("GMAIL_ADDRESS"),
		os.Getenv("GMAIL_PASSWORD"),
		os.Getenv("MAIL_HOST"),
	)
	msg := []byte("Subject: my special subject\n" + "dear " + notification["to"] + "!" + "\n" + notification["message"])
	err := smtp.SendMail(
		os.Getenv("MAIL_HOST_ADDRESS"),
		auth,
		os.Getenv("GMAIL_ADDRESS"),
		[]string{os.Getenv("GMAIL_ADDRESS")},
		msg,
	)
	if err != nil {
		return err
	}
	return nil
}
