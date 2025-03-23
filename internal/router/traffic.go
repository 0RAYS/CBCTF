package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func LoadTraffic(ctx *gin.Context) {
	container := middleware.GetContainer(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := db.SaveTraffic(tx, container)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
}

func GetTraffics(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	container := middleware.GetContainer(ctx)
	traffics, count, ok, msg := db.GetTrafficByColumn(db.DB.WithContext(ctx), "container_id", container.ID, form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"traffics": &traffics, "count": count}})
}

func DownloadTraffic(ctx *gin.Context) {
	container := middleware.GetContainer(ctx)
	if _, err := os.Stat(container.TrafficPath()); err != nil {
		log.Logger.Warningf("Failed to get file: %s", err)
		if errors.Is(err, os.ErrNotExist) {
			ctx.JSON(http.StatusOK, gin.H{"msg": "FileNotFound", "data": nil})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"msg": "UnknownError", "data": nil})
		return
	}
	ctx.File(container.TrafficPath())
}
