package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/middleware"
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetFlags(ctx *gin.Context) {
	usage := middleware.GetUsage(ctx)
	repo := db.InitFlagRepo(db.DB.WithContext(ctx))
	flags, _, ok, msg := repo.GetByKeyID("usage_id", usage.ID, -1, -1, true, 3)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": &flags})
}

func UpdateFlag(ctx *gin.Context) {
	var form f.UpdateFlagForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	flag := middleware.GetFlag(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.UpdateFlag(tx, flag, form)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
