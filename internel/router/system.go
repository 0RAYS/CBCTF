package router

import (
	"CBCTF/internel/config"
	"CBCTF/internel/i18n"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SystemConfig(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": config.Env})
}
