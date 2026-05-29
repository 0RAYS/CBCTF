package email

import (
	"CBCTF/internal/config"
	"fmt"
)

const ResetPasswordEmailSubject = "Reset Your Password"

func SendResetPasswordEmail(to, token, id string) error {
	link := fmt.Sprintf("%s/platform/#/reset-password?token=%s&id=%s", config.Env.Host, token, id)
	html := buildHTML(
		"Reset Your Password",
		"We received a request to reset the password for your account. Click the button below to choose a new password.",
		"Reset Password",
		link,
		"This link expires in <strong style=\"color:#8a8a8a;font-weight:500;\">30 minutes</strong>. If you did not request this, no action is needed — your password remains unchanged.",
	)
	return SendEmail(to, ResetPasswordEmailSubject, html)
}
