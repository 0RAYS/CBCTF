package email

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
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

var Senders = make([]*Sender, 0)

func Init() {
	smtpL, _, ok, _ := db.InitSmtpRepo(db.DB).List(-1, -1, db.GetOptions{Conditions: map[string]any{"on": true}})
	if !ok {
		log.Logger.Warningf("No smtp configured, email sending will be failed")
		return
	}
	for _, smtp := range smtpL {
		dialer := gomail.NewDialer(smtp.Host, smtp.Port, smtp.Address, smtp.Pwd)
		auth, err := dialer.Dial()
		if err != nil {
			log.Logger.Warningf("Failed to connect to email server %s:%d: %s", smtp.Host, smtp.Port, err)
			continue
		}
		Senders = append(Senders, &Sender{
			Auth:      &auth,
			Smtp:      smtp,
			CreatedAt: time.Now(),
		})
	}
	if len(Senders) == 0 {
		log.Logger.Warningf("No smtp configured, email sending will be failed")
	}
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
