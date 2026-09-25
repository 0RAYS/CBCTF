package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
)

func GetSmtpResp(smtp model.Smtp) gin.H {
	return gin.H{
		"id":           smtp.ID,
		"address":      smtp.Address,
		"host":         smtp.Host,
		"port":         smtp.Port,
		"on":           smtp.On,
		"success":      smtp.Success,
		"success_last": smtp.SuccessLast,
		"failure":      smtp.Failure,
		"failure_last": smtp.FailureLast,
	}
}
