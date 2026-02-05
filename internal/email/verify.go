package email

import (
	"CBCTF/internal/config"
	"fmt"
)

const VerifyEmailSubject = "Verify Email"

func SendVerifyEmail(to, token, id string) error {
	content := fmt.Sprintf(
		"Please click the following link to verify your email:\n%s/verify?token=%s&id=%s",
		config.Env.Host, token, id,
	)
	return SendEmail(to, VerifyEmailSubject, content)
}
