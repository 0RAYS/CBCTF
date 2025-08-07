package email

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"gopkg.in/gomail.v2"
)

type Sender struct {
	Auth       *gomail.SendCloser
	Addr       string
	Host       string
	Pwd        string
	Port       int
	CreatedAt  time.Time
	UpdateLock sync.Mutex
}

var (
	Senders = make([]*Sender, 0)
)

func Init() {
	for _, sender := range config.Env.Email.Senders {
		dialer := gomail.NewDialer(sender.Host, sender.Port, sender.Addr, sender.Pwd)
		auth, err := dialer.Dial()
		if err != nil {
			log.Logger.Warningf("Failed to connect to email server %s:%d: %s", sender.Host, sender.Port, err)
			continue
		}
		Senders = append(Senders, &Sender{
			Auth:      &auth,
			Addr:      sender.Addr,
			Host:      sender.Host,
			Port:      sender.Port,
			Pwd:       sender.Pwd,
			CreatedAt: time.Now(),
		})
	}
	if len(Senders) == 0 {
		log.Logger.Warningf("No sender configured, email sending will be failed")
	}
}

func Redial(old *Sender) error {
	dialer := gomail.NewDialer(old.Host, old.Port, old.Addr, old.Pwd)
	auth, err := dialer.Dial()
	if err != nil {
		log.Logger.Warningf("Failed to connect to email server %s:%d: %s", old.Host, old.Port, err)
		return err
	}
	old.Auth = &auth
	old.CreatedAt = time.Now()
	log.Logger.Debugf("Redialed email server %s:%d successfully", old.Host, old.Port)
	return nil
}

func SendVerifyEmail(to, token, id string) error {
	log.Logger.Debugf("Sending verify email to %s", to)
	if len(Senders) == 0 {
		return fmt.Errorf("no email sender configured")
	}
	var sender *Sender
	var count = 0
	for {
		count++
		sender = Senders[rand.Intn(len(Senders))]
		sender.UpdateLock.Lock()
		if sender.CreatedAt.Add(time.Minute).After(time.Now()) {
			sender.UpdateLock.Unlock()
			break
		}
		if Redial(sender) == nil {
			sender.UpdateLock.Unlock()
			break
		}
		sender.UpdateLock.Unlock()
		if count > 5 {
			return fmt.Errorf("failed too many times to connect smtp servers")
		}
	}
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
