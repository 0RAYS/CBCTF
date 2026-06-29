package email

import (
	"CBCTF/internal/config"
	"fmt"
)

const VerifyEmailSubject = "Verify Your Email Address"

func SendVerifyEmail(to, token string) error {
	link := fmt.Sprintf("%s/platform/#/verify?token=%s", config.Env.Host, token)
	html := buildHTML(
		"Verify Your Email Address",
		"Thanks for signing up. Click the button below to verify your email address and activate your account.",
		"Verify Email Address",
		link,
		"If you did not create an account, you can safely ignore this email.",
	)
	return SendEmail(to, VerifyEmailSubject, html)
}
