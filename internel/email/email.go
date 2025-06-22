package email

import (
	"CBCTF/internel/config"
	"CBCTF/internel/log"
	"fmt"
	"gopkg.in/gomail.v2"
	"math/rand"
	"regexp"
)

type Sender struct {
	Auth *gomail.SendCloser
	Addr string
	Host string
	Port int
}

var (
	Senders = make([]Sender, 0)
)

func Init() {
	for _, sender := range config.Env.Email.Senders {
		dialer := gomail.NewDialer(sender.Host, sender.Port, sender.Addr, sender.Pwd)
		auth, err := dialer.Dial()
		if err != nil {
			log.Logger.Warningf("Failed to connect to email server %s:%d: %v", sender.Host, sender.Port, err)
			continue
		}
		Senders = append(Senders, Sender{
			Auth: &auth,
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
	m := gomail.NewMessage()
	m.SetHeader("From", sender.Addr)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Verify Email")
	m.SetBody("text/plain",
		fmt.Sprintf(
			"Please click the following link to verify your email:\n%s/verify?token=%s&id=%s",
			config.Env.Backend, token, id,
		),
	)
	return gomail.Send(*sender.Auth, m)
}
