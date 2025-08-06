package middleware

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckUnVerified(ctx *gin.Context) {
	user := GetSelf(ctx).(model.User)
	if user.Verified {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.AlreadyVerifiedEmail, "data": nil})
		return
	}
	ctx.Next()
}
