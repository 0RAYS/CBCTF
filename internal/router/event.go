package router

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetEventTypes(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": model.EventTypes})
}
