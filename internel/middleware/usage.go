package middleware

import (
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CheckGenerated model.Team 是否初始化 model.Usage
func CheckGenerated(ctx *gin.Context) {
	usage := GetUsage(ctx)
	team := GetTeam(ctx)
	if generated, _, msg := service.IsGenerated(db.DB.WithContext(ctx), usage, team); !generated {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Next()
}
