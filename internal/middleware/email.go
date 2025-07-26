package middleware

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckUnVerified(ctx *gin.Context) {
	user := GetSelf(ctx).(model.User)
	if user.Verified {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
		return
	}
	ctx.Next()
}
