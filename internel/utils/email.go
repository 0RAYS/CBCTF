package utils

import (
	"CBCTF/internel/config"
	"fmt"
	"math/rand"
	"net/smtp"
	"regexp"
)

// IsValidEmail 邮箱格式验证
func IsValidEmail(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	if regexp.MustCompile(pattern).MatchString(email) {
		return true
	}
	return false
}

func SendVerifyEmail(to, token, id string) error {
	sender := config.Env.Email.Senders[rand.Intn(len(config.Env.Email.Senders))]
	auth := smtp.PlainAuth("", sender.Addr, sender.Pwd, sender.Host)

	toList := []string{to}
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Verify Email\r\n\r\n"+
		"Please click the following link to verify your email: "+
		fmt.Sprintf("%s/verify?token=%s&id=%s\r\n", config.Env.Backend, token, id), to))

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", sender.Host, sender.Port),
		auth,
		sender.Addr,
		toList,
		msg,
	)
}
