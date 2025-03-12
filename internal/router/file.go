package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
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

func DownloadFile(ctx *gin.Context) {
	file := middleware.GetFile(ctx)
	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		tx := db.DB.WithContext(ctx).Begin()
		if ok, _ := db.DeleteFile(tx, file.ID); !ok {
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

func DeleteFile(ctx *gin.Context) {
	var (
		form f.DeleteFileForm
		file model.File
		ok   bool
		msg  string
		DB   = db.DB.WithContext(ctx)
	)
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	filesID := form.FilesID
	var files []model.File
	if uri := middleware.GetFile(ctx); uri.ID != "" {
		files = append(files, uri)
	}
	for _, id := range filesID {
		file, ok, msg = db.GetFileByID(DB, id)
		if !ok {
			ctx.JSON(http.StatusNotFound, gin.H{"msg": msg, "data": nil})
			return
		}
		files = append(files, file)
	}
	for _, file = range files {
		tx := DB.Begin()
		if ok, msg = db.DeleteFile(tx, file.ID); !ok {
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
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	files, count, ok, msg := db.GetAvatars(db.DB.WithContext(ctx), form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "files": &files}})
}

func UploadAvatar(v string) func(ctx *gin.Context) {
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
			record model.File
			ok     bool
			msg    string
			path   string
		)
		allowed := []string{".png", ".jpg", ".jpeg"}
		suffix := strings.ToLower(p.Ext(file.Filename))
		if !utils.In(suffix, allowed) {
			ctx.JSON(http.StatusForbidden, gin.H{"msg": "FileNotAllowed", "data": file.Filename})
			return
		}
		tx := db.DB.WithContext(ctx).Begin()
		hash := hex.EncodeToString(sha256Sum.Sum(nil))
		if record, ok, _ = db.GetFileByHash(tx, hash); !ok {
			basePath := fmt.Sprintf("%s/avatar", config.Env.Gin.Upload.Path)
			path = fmt.Sprintf("%s/%s%s", basePath, utils.UUID(), suffix)
			if err = ctx.SaveUploadedFile(file, path); err != nil {
				tx.Rollback()
				ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "UnknownError", "data": nil})
				return
			}
		} else {
			path = record.Path
		}
		record, ok, msg = db.RecordFile(tx, path, middleware.GetSelfID(ctx), file, hash, model.Avatar)
		if !ok {
			tx.Rollback()
			ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		path = fmt.Sprintf("/avatars/%s", record.ID)
		switch v {
		case "self-admin":
			ok, msg = db.UpdateAdmin(tx, middleware.GetSelfID(ctx), map[string]interface{}{"avatar": path})
		case "self-user":
			ok, msg = db.UpdateUser(tx, middleware.GetSelfID(ctx), map[string]interface{}{"avatar": path})
		case "user":
			ok, msg = db.UpdateUser(tx, middleware.GetUser(ctx).ID, map[string]interface{}{"avatar": path})
		case "contest":
			ok, msg = db.UpdateContest(tx, middleware.GetContest(ctx).ID, map[string]interface{}{"avatar": path})
		case "team":
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

func UploadWriteUp(ctx *gin.Context) {
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
		record  model.File
		ok      bool
		msg     string
		contest = middleware.GetContest(ctx)
		team    = middleware.GetTeam(ctx)
	)
	hash := hex.EncodeToString(sha256Sum.Sum(nil))
	tx := db.DB.WithContext(ctx).Begin()
	if record, ok, _ = db.GetFileByHash(tx, hash); !ok {
		basePath := fmt.Sprintf("%s/writeups/%d/%d", config.Env.Gin.Upload.Path, contest.ID, team.ID)
		allowed := []string{".pdf", ".docx", ".doc"}
		suffix := strings.ToLower(p.Ext(file.Filename))
		if !utils.In(suffix, allowed) {
			tx.Rollback()
			ctx.JSON(http.StatusForbidden, gin.H{"msg": "FileNotAllowed", "data": file.Filename})
			return
		}
		path := fmt.Sprintf("%s/%s%s", basePath, time.Now().Format("2006_06_02_15_04_05"), suffix)
		if err = ctx.SaveUploadedFile(file, path); err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "UnknownError", "data": nil})
			return
		}
		record, ok, msg = db.RecordFile(tx, path, middleware.GetSelfID(ctx), file, hash, model.Avatar)
		if !ok {
			tx.Rollback()
			ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
	}
	tx.Commit()
	path := fmt.Sprintf("/writeups/%s", record.ID)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": path})
}

func GetWriteUPs(ctx *gin.Context) {
	contest := middleware.GetContest(ctx)
	team := middleware.GetTeam(ctx)
	path := fmt.Sprintf("%s/writeups/%d/%d", config.Env.Gin.Upload.Path, contest.ID, team.ID)
	dir, err := os.ReadDir(path)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": "FileNotFound", "data": nil})
		return
	}
	var files []string
	for _, file := range dir {
		files = append(files, file.Name())
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": &files})
}
