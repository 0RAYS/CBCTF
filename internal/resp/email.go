package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
)

func GetEmailResp(email model.Email) gin.H {
	return gin.H{
		"id":      email.ID,
		"smtp_id": email.SmtpID,
		"from":    email.From,
		"to":      email.To,
		"subject": email.Subject,
		"content": email.Content,
		"time":    email.CreatedAt,
		"success": email.Success,
	}
}
