package email

import (
	"CBCTF/internal/config"
	"fmt"
)

const ResetPasswordEmailSubject = "Reset Password"

func SendResetPasswordEmail(to, token, id string) error {
	content := fmt.Sprintf(
		"Please click the following link to reset your password (valid for 30 minutes):\n%s/platform/#/reset-password?token=%s&id=%s",
		config.Env.Host, token, id,
	)
	return SendEmail(to, ResetPasswordEmailSubject, content)
}
