package email

import (
	"CBCTF/internel/config"
	"CBCTF/internel/log"
	"fmt"
	"math/rand"
	"net/smtp"
	"regexp"
)

type Sender struct {
	Auth smtp.Auth
	Addr string
	Host string
	Port int
}

var (
	Senders = make([]Sender, 0)
)

func Init() {
	for _, sender := range config.Env.Email.Senders {
		auth := smtp.PlainAuth("", sender.Addr, sender.Pwd, sender.Host)
		Senders = append(Senders, Sender{
			Auth: auth,
			Addr: sender.Addr,
			Host: sender.Host,
			Port: sender.Port,
		})
	}
	if len(Senders) == 0 {
		log.Logger.Warningf("No sender configured, email sending will be failed")
	}
}

// IsValidEmail 邮箱格式验证
func IsValidEmail(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	if regexp.MustCompile(pattern).MatchString(email) {
		return true
	}
	return false
}

func SendVerifyEmail(to, token, id string) error {
	if len(Senders) == 0 {
		return fmt.Errorf("no email sender configured")
	}
	sender := Senders[rand.Intn(len(Senders))]
	toList := []string{to}
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: Verify Email\r\n\r\n"+
		"Please click the following link to verify your email:\r\n"+
		fmt.Sprintf("%s/verify?token=%s&id=%s\r\n", config.Env.Backend, token, id), sender.Addr, to))

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", sender.Host, sender.Port),
		sender.Auth,
		sender.Addr,
		toList,
		msg,
	)
}
