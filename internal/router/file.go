package router

import (
	"RayWar/internal/config"
	"RayWar/internal/db"
	"RayWar/internal/log"
	"RayWar/internal/middleware"
	"RayWar/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	p "path"
	"time"
)

func Upload(ctx *gin.Context) {
	trace := middleware.GetTraceID(ctx)
	form, err := ctx.MultipartForm()
	if err != nil {
		log.Logger.Infof("| %s | BadRequest: %s", trace, err.Error())
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		return
	}
	files := form.File["files"]
	if len(files) == 0 {
		msg := "BadRequest"
		log.Logger.Infof("| %s | BadRequest: %s", trace, msg)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
		return
	}
	basePath := config.Env.GetString("upload.path")
	allowed := []string{".png", ".jpg", ".jpeg", ".zip", ".rar", ".gz", ".tar"}
	for _, file := range files {
		random := utils.RandomString()
		suffix := p.Ext(file.Filename)
		if !utils.In(suffix, allowed) {
			msg := "FileNotAllowed"
			log.Logger.Infof("| %s | %s", trace, msg)
			ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": file.Filename})
			return
		}
		path := fmt.Sprintf("%s/%s/%s%s", basePath, time.Now().Format("2006-01-02"), random, suffix)
		if err := ctx.SaveUploadedFile(file, path); err != nil {
			log.Logger.Infof("| %s | %s", trace, err.Error())
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, "UnknownError"), "data": nil})
			return
		}
		if _, ok, msg := db.RecordFile(middleware.GetSelf(ctx).ID, random, path, file); !ok {
			log.Logger.Infof("| %s | %s", trace, msg)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
	}
	ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, "Success"), "data": len(files)})
	return
}

func Download(ctx *gin.Context) {
	trace := middleware.GetTraceID(ctx)
	type fileIDUri struct {
		FileID string `uri:"fileID" binding:"required"`
	}
	var fileID fileIDUri
	if err := ctx.ShouldBindUri(&fileID); err != nil {
		log.Logger.Infof("| %s | %s", trace, err)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		ctx.Abort()
		return
	}
	file, ok, msg := db.GetFile(fileID.FileID)
	if !ok {
		msg = "FileNotFound"
		log.Logger.Infof("| %s | %s", trace, msg)
		ctx.JSONP(http.StatusNotFound, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": fileID.FileID})
		return
	}
	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		_, _ = db.DeleteFile(fileID.FileID)
		msg = "FileNotFound"
		log.Logger.Infof("| %s | %s", trace, msg)
		ctx.JSONP(http.StatusNotFound, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": fileID.FileID})
		return
	}
	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename="+file.Name)
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	ctx.File(file.Path)
	return
}
