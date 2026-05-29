package email

import (
	"CBCTF/internal/config"
	"fmt"
)

const VerifyEmailSubject = "Verify Your Email Address"

func SendVerifyEmail(to, token, id string) error {
	link := fmt.Sprintf("%s/platform/#/verify?token=%s&id=%s", config.Env.Host, token, id)
	html := buildHTML(
		"Verify Your Email Address",
		"Thank you for registering! Please click the button below to verify your email address and activate your account.",
		"Verify Email Address",
		link,
		"If you did not create an account, you can safely ignore this email.",
	)
	return SendEmail(to, VerifyEmailSubject, html)
}
