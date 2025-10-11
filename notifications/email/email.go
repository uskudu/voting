package email

import (
	"fmt"
	"net/smtp"
	"os"
	"voting/notifications/rmq"
)

func SendMail(n rmq.VoteNotification) error {
	auth := smtp.PlainAuth(
		"",
		os.Getenv("GMAIL_ADDRESS"),
		os.Getenv("GMAIL_PASSWORD"),
		os.Getenv("MAIL_HOST"),
	)
	subject := "subject: new vote notification\n"
	body := fmt.Sprintf("Dear %s,\n\n%s\n", n.To, n.Message)
	msg := []byte(subject + "\n" + body)
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
