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

func AdminChangePassword(ctx *gin.Context) {
	var form f.ChangePasswordForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.ChangeAdminPassword(tx, middleware.GetSelf(ctx).(model.Admin), form)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

