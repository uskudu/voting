package main

import (
	"fmt"
	"net/smtp"
)

func SendMail() {
	auth := smtp.PlainAuth(
		"",
		"usskudu@gmail.com",
		"yovfhlzjrtrdqrbw",
		"smtp.gmail.com",
	)
	msg := []byte("Subject: my special subject\nthis is body")
	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"usskudu@gmail.com",
		[]string{"usskudu@gmail.com"},
		msg,
	)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	SendMail()
}
