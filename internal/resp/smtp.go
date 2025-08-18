package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetSmtpResp(smtp model.Smtp) gin.H {
	return gin.H{
		"id":           smtp.ID,
		"address":      smtp.Address,
		"host":         smtp.Host,
		"port":         smtp.Port,
		"pwd":          smtp.Pwd,
		"on":           smtp.On,
		"success":      smtp.Success,
		"success_last": smtp.SuccessLast,
		"failure":      smtp.Failure,
		"failure_last": smtp.FailureLast,
	}
}
