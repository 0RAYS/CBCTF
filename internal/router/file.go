package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	p "path"
	"strings"
	"time"
)

func Upload(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil || len(form.File["files"]) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	basePath := config.Env.GetString("upload.path")
	allowed := []string{".png", ".jpg", ".jpeg", ".zip", ".rar", ".gz", ".tar", ".7z"}
	var records []model.File
	for _, file := range form.File["files"] {
		suffix := strings.ToLower(p.Ext(file.Filename))
		if !utils.In(suffix, allowed) {
			ctx.JSON(http.StatusForbidden, gin.H{"msg": "FileNotAllowed", "data": file.Filename})
			return
		}
		path := fmt.Sprintf("%s/%s/%s%s", basePath, time.Now().Format("2006-01-02"), utils.RandomString(), suffix)
		if err = ctx.SaveUploadedFile(file, path); err != nil {
			ctx.JSONP(http.StatusInternalServerError, gin.H{"msg": "UnknownError", "data": nil})
			return
		}
		record, ok, msg := db.RecordFile(ctx, path, middleware.GetSelfID(ctx), middleware.GetRole(ctx) == "admin", file)
		if !ok {
			ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		records = append(records, record)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"count": len(records), "records": records}})
}

func Download(ctx *gin.Context) {
	file, ok, msg := db.GetFile(ctx, middleware.GetFileID(ctx))
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": msg, "data": nil})
		return
	}
	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		_, _ = db.DeleteFile(ctx, file.ID)
		ctx.JSONP(http.StatusNotFound, gin.H{"msg": "FileNotFound", "data": file.ID})
		return
	}
	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename="+file.Filename)
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	ctx.File(file.Path)
}

func Avatar(v interface{}) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile("avatar")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
			return
		}
		basePath := fmt.Sprintf("%s/avatar", config.Env.GetString("upload.path"))
		allowed := []string{".png", ".jpg", ".jpeg"}
		suffix := strings.ToLower(p.Ext(file.Filename))
		if !utils.In(suffix, allowed) {
			ctx.JSON(http.StatusForbidden, gin.H{"msg": "FileNotAllowed", "data": file.Filename})
			return
		}
		path := fmt.Sprintf("%s/%s/%s%s", basePath, time.Now().Format("2006-01-02"), utils.RandomString(), suffix)
		if err = ctx.SaveUploadedFile(file, path); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "UnknownError", "data": nil})
			return
		}
		record, ok, msg := db.RecordFile(ctx, path, middleware.GetSelfID(ctx), middleware.GetRole(ctx) == "admin", file)
		if !ok {
			ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		path = fmt.Sprintf("%s/download/%s", config.Env.GetString("host"), record.ID)
		switch v.(type) {
		case model.Admin:
			_, msg = db.UpdateAdmin(ctx, middleware.GetSelfID(ctx), map[string]interface{}{"avatar": path})
		case model.User:
			_, msg = db.UpdateUser(ctx, middleware.GetSelfID(ctx), map[string]interface{}{"avatar": path})
		case model.Contest:
			_, msg = db.UpdateContest(ctx, middleware.GetContestID(ctx), map[string]interface{}{"avatar": path})
		case model.Team:
			_, msg = db.UpdateTeam(ctx, middleware.GetTeamID(ctx), map[string]interface{}{"avatar": path})
		}
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": path})
	}
}
