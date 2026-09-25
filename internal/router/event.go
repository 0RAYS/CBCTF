package router

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
	"CBCTF/internal/resp"
)

func GetEventTypes(ctx *gin.Context) {
	resp.JSON(ctx, model.SuccessRetVal(model.EventTypes))
}
