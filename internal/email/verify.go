package email

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"gopkg.in/gomail.v2"
)

func SendVerifyEmail(to, token, id string) error {
	log.Logger.Debugf("Sending verify email to %s", to)
	if len(Senders) == 0 {
		return fmt.Errorf("no email sender configured")
	}
	var sender *Sender
	var count = 0
	for {
		count++
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(Senders))))
		sender = Senders[index.Int64()]
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
	m.SetHeader("Subject", "Verify Email")
	m.SetBody("text/plain",
		fmt.Sprintf(
			"Please click the following link to verify your email:\n%s/verify?token=%s&id=%s",
			config.Env.Backend, token, id,
		),
	)
	return gomail.Send(*sender.Auth, m)
}
