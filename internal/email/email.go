package email

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
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
