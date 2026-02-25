package router

import (
	"CBCTF/internal/model"
	"CBCTF/internal/resp"

	"github.com/gin-gonic/gin"
)

func GetEventTypes(ctx *gin.Context) {
	resp.JSON(ctx, model.SuccessRetVal(model.EventTypes))
}
