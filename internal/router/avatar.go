package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/constants"
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
)

func DownloadAvatar(ctx *gin.Context) {
	file := middleware.GetAvatar(ctx)
	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		tx := db.DB.WithContext(ctx).Begin()
		if ok, _ := db.DeleteAvatar(tx, file.ID); !ok {
			tx.Rollback()
			ctx.JSONP(http.StatusNotFound, gin.H{"msg": "FileNotFound", "data": file.ID})
			return
		}
		tx.Commit()
	}
	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename="+file.Filename)
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	ctx.File(file.Path)
}

func DeleteAvatar(ctx *gin.Context) {
	var (
		form constants.DeleteFileForm
		file model.Avatar
		ok   bool
		msg  string
	)
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	filesID := form.FilesID
	var files []model.Avatar
	if uri := middleware.GetAvatar(ctx); uri.ID != "" {
		files = append(files, uri)
	}
	for _, id := range filesID {
		file, ok, msg = db.GetAvatarByID(ctx, id)
		if !ok {
			ctx.JSON(http.StatusNotFound, gin.H{"msg": msg, "data": nil})
			return
		}
		files = append(files, file)
	}
	for _, file = range files {
		tx := db.DB.WithContext(ctx).Begin()
		if ok, msg = db.DeleteAvatar(tx, file.ID); !ok {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": msg, "data": nil})
			return
		}
		tx.Commit()
		if form.Force && os.Remove(file.Path) != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "UnknownError", "data": nil})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
}

func GetAvatars(ctx *gin.Context) {
	var form constants.GetModelsForm
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
		tx := db.DB.WithContext(ctx).Begin()
		if record, ok, _ = db.GetAvatarByHash(ctx, hash); !ok {
			basePath := fmt.Sprintf("%s/avatar", config.Env.Gin.Upload.Path)
			allowed := []string{".png", ".jpg", ".jpeg"}
			suffix := strings.ToLower(p.Ext(file.Filename))
			if !utils.In(suffix, allowed) {
				tx.Rollback()
				ctx.JSON(http.StatusForbidden, gin.H{"msg": "FileNotAllowed", "data": file.Filename})
				return
			}
			path := fmt.Sprintf("%s/%s%s", basePath, utils.RandomString(), suffix)
			if err = ctx.SaveUploadedFile(file, path); err != nil {
				tx.Rollback()
				ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "UnknownError", "data": nil})
				return
			}
			record, ok, msg = db.RecordAvatar(tx, path, middleware.GetSelfID(ctx), file, hash)
			if !ok {
				tx.Rollback()
				ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
				return
			}
		}
		path := fmt.Sprintf("/avatar/%s", record.ID)
		switch v.(type) {
		case model.Admin:
			ok, msg = db.UpdateAdmin(tx, middleware.GetSelfID(ctx), map[string]interface{}{"avatar": path})
		case model.User:
			ok, msg = db.UpdateUser(tx, middleware.GetSelfID(ctx), map[string]interface{}{"avatar": path})
		case model.Contest:
			ok, msg = db.UpdateContest(tx, middleware.GetContest(ctx).ID, map[string]interface{}{"avatar": path})
		case model.Team:
			ok, msg = db.UpdateTeam(tx, middleware.GetTeam(ctx).ID, map[string]interface{}{"avatar": path})
		}
		if !ok {
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		tx.Commit()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": path})
	}
}
