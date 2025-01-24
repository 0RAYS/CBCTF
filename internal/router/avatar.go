package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	p "path"
	"strings"
	"time"
)

func DownloadAvatar(ctx *gin.Context) {
	file, ok, msg := db.GetAvatarByID(ctx, middleware.GetFileID(ctx))
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": msg, "data": nil})
		return
	}
	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		_, _ = db.DeleteAvatar(ctx, file.ID)
		ctx.JSONP(http.StatusNotFound, gin.H{"msg": "FileNotFound", "data": file.ID})
		return
	}
	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename="+file.Filename)
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	ctx.File(file.Path)
}

func DeleteAvatar(ctx *gin.Context) {
	var form DeleteFileForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	files := form.Files
	if id := middleware.GetFileID(ctx); id != "" {
		files = append(files, id)
	}
	for _, id := range files {
		file, ok, msg := db.GetAvatarByID(ctx, id)
		if !ok {
			ctx.JSON(http.StatusNotFound, gin.H{"msg": msg, "data": nil})
			return
		}
		if ok, msg = db.DeleteAvatar(ctx, file.ID); !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": msg, "data": nil})
			return
		}
		if form.Force && os.Remove(file.Path) != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "UnknownError", "data": nil})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
}

func GetAvatars(ctx *gin.Context) {
	var form GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	files, count, ok, msg := db.GetAvatars(ctx, form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "files": files}})
}

func UploadAvatar(v interface{}) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile("avatar")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
			return
		}
		src, err := file.Open()
		if err != nil {
			log.Logger.Warningf("Failed to open file: %v", err)
			ctx.JSON(http.StatusOK, gin.H{"msg": "BadRequest", "data": nil})
			return
		}
		defer func(src multipart.File) {
			err := src.Close()
			if err != nil {
				log.Logger.Warningf("Failed to close file: %v", err)
			}
		}(src)
		sha256Sum := sha256.New()
		if _, err := io.Copy(sha256Sum, src); err != nil {
			log.Logger.Warningf("Failed to hash file: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "UnknownError", "data": nil})
			return
		}
		var (
			record model.Avatar
			ok     bool
			msg    string
		)
		hash := hex.EncodeToString(sha256Sum.Sum(nil))
		if record, ok, _ = db.GetAvatarByHash(ctx, hash); !ok {
			basePath := fmt.Sprintf("%s/avatar", config.Env.Gin.Upload.Path)
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
			record, ok, msg = db.RecordAvatar(ctx, path, middleware.GetSelfID(ctx), file, hash)
			if !ok {
				ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
				return
			}
		}
		path := fmt.Sprintf("%s/download/%s", config.Env.Backend, record.ID)
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
