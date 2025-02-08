package utils

import (
	"CBCTF/internal/config"
	"fmt"
	"math/rand"
	"net/smtp"
)

func SendVerifyEmail(to, token, id string) error {
	sender := config.Env.Email.Senders[rand.Intn(len(config.Env.Email.Senders))]
	auth := smtp.PlainAuth("", sender.Address, sender.Password, sender.Host)

	toList := []string{to}
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Verify Email\r\n\r\n"+
		"Please click the following link to verify your email: "+
		fmt.Sprintf("%s/verify?token=%s&id=%s\r\n", config.Env.Backend, token, id), to))

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", sender.Host, sender.Port),
		auth,
		sender.Address,
		toList,
		msg,
	)
}
