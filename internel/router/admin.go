package router

import (
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	"CBCTF/internel/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAdmin(ctx *gin.Context) {
	admin := middleware.GetSelf(ctx).(model.Admin)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": resp.GetAdminResp(admin)})
}
