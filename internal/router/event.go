package router

import (
	"CBCTF/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetEventTypes(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, model.SuccessRetVal(model.EventTypes))
}
