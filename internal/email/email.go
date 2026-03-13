package email

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"crypto/rand"
	"fmt"
	"math/big"
	"slices"
	"sync"
	"time"

	"gopkg.in/gomail.v2"
)

type Sender struct {
	Auth       *gomail.SendCloser
	Smtp       model.Smtp
	CreatedAt  time.Time
	UpdateLock sync.Mutex
}

var (
	Senders []*Sender
	lock    sync.RWMutex
)

func Init() {
	Senders = make([]*Sender, 0)
	smtpL, _, ret := db.InitSmtpRepo(db.DB).List(-1, -1, db.GetOptions{Conditions: map[string]any{"on": true}})
	if !ret.OK {
		log.Logger.Warningf("No smtp configured, email sending will be failed")
		return
	}
	for _, smtp := range smtpL {
		go AddSenders(smtp)
	}
}

func AddSenders(smtp model.Smtp) {
	dialer := gomail.NewDialer(smtp.Host, smtp.Port, smtp.Address, smtp.Pwd)
	auth, err := dialer.Dial()
	if err != nil {
		log.Logger.Warningf("Failed to connect to email server %s:%d: %s", smtp.Host, smtp.Port, err)
		return
	}
	lock.Lock()
	Senders = append(Senders, &Sender{
		Auth:      &auth,
		Smtp:      smtp,
		CreatedAt: time.Now(),
	})
	lock.Unlock()
}

func DelSenders(smtp model.Smtp) {
	lock.Lock()
	Senders = slices.DeleteFunc(Senders, func(s *Sender) bool {
		return s != nil && s.Smtp.ID == smtp.ID
	})
	lock.Unlock()
}

func Redial(old *Sender) error {
	dialer := gomail.NewDialer(old.Smtp.Host, old.Smtp.Port, old.Smtp.Address, old.Smtp.Pwd)
	auth, err := dialer.Dial()
	if err != nil {
		log.Logger.Warningf("Failed to connect to email server %s:%d: %s", old.Smtp.Host, old.Smtp.Port, err)
		return err
	}
	old.Auth = &auth
	old.CreatedAt = time.Now()
	log.Logger.Debugf("Redialed email server %s:%d successfully", old.Smtp.Host, old.Smtp.Port)
	return nil
}

func SendEmail(to, subject, content string) error {
	log.Logger.Debugf("Sending verify email to %s", to)
	if len(Senders) == 0 {
		return fmt.Errorf("no email sender configured")
	}
	var sender *Sender
	var count = 0
	for {
		count++
		lock.RLock()
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(Senders))))
		sender = Senders[index.Int64()]
		lock.RUnlock()
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
	m.SetHeader("From", sender.Smtp.Address)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", content)
	options := db.CreateEmailOptions{
		SmtpID:  sender.Smtp.ID,
		From:    sender.Smtp.Address,
		To:      to,
		Subject: subject,
		Content: content,
		Success: true,
	}
	err := gomail.Send(*sender.Auth, m)
	if err != nil {
		options.Success = false
	}
	db.InitEmailRepo(db.DB).Create(options)
	return err
}
