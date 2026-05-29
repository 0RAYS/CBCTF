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
		"We received a request to reset the password for your account. Click the button below to set a new password.",
		"Reset Password",
		link,
		"This link is valid for <strong style=\"color:#c8c8c8;\">30 minutes</strong>. If you did not request a password reset, please ignore this email — your password will remain unchanged.",
	)
	return SendEmail(to, ResetPasswordEmailSubject, html)
}
