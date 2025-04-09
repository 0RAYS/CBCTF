package router

import (
	"CBCTF/internel/config"
	f "CBCTF/internel/form"
	"CBCTF/internel/log"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"CBCTF/internel/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func DownloadFile(ctx *gin.Context) {
	file := middleware.GetFile(ctx)
	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		tx := db.DB.WithContext(ctx).Begin()
		if ok, _ := db.InitFileRepo(tx).Delete(file.ID); !ok {
			tx.Rollback()
			ctx.JSON(http.StatusNotFound, gin.H{"msg": "FileNotFound", "data": file.ID})
			return
		}
		tx.Commit()
	}
	ctx.Writer.Header().Add("Content-Disposition", "attachment; filename="+file.Filename)
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	ctx.File(file.Path)
}

func DownloadChallenge(ctx *gin.Context) {
	var form f.DownloadChallengeForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	challenge := middleware.GetChallenge(ctx)
	var path string
	switch form.File {
	case model.AttachmentFile, model.GeneratorFile:
		path = fmt.Sprintf("%s/%s", challenge.BasicDir(), form.File)
	default:
		ctx.JSON(http.StatusOK, gin.H{"msg": "InvalidFileName", "data": nil})
		return
	}
	if _, err := os.Stat(path); err != nil {
		log.Logger.Warningf("Failed to get file: %s", err)
		if errors.Is(err, os.ErrNotExist) {
			ctx.JSON(http.StatusNotFound, gin.H{"msg": "FileNotFound", "data": nil})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"msg": "UnknownError", "data": nil})
		return
	}
	ctx.File(path)
}

func UploadAvatar(v string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile(model.Avatar)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
			return
		}
		tx := db.DB.WithContext(ctx).Begin()
		record, ok, msg := service.SaveAvatar(tx, middleware.GetSelfID(ctx), file)
		if !ok {
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		if err = ctx.SaveUploadedFile(file, record.Path); err != nil {
			tx.Rollback()
			log.Logger.Warningf("Failed to save file: %s", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "UnknownError", "data": nil})
			return
		}
		var id uint
		switch v {
		case "admin", "self-user":
			id = middleware.GetSelfID(ctx)
		case "user":
			id = middleware.GetUser(ctx).ID
		case "contest":
			id = middleware.GetContest(ctx).ID
		case "team":
			id = middleware.GetTeam(ctx).ID
		}
		path, ok, msg := service.UpdateAvatar(tx, v, id, record)
		if !ok {
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		tx.Commit()
		path = fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(path, "/"))
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": path})
	}
}

func UploadWriteUp(ctx *gin.Context) {
	file, err := ctx.FormFile(model.WriteUP)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	record, ok, msg := service.SaveWriteUp(tx, middleware.GetContest(ctx).ID, middleware.GetTeam(ctx).ID, file)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if err = ctx.SaveUploadedFile(file, record.Path); err != nil {
		tx.Rollback()
		log.Logger.Warningf("Failed to save file: %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "UnknownError", "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func GetAttachment(ctx *gin.Context) {
	usage := middleware.GetUsage(ctx)
	team := middleware.GetTeam(ctx)
	path := usage.Challenge.AttachmentPath(team.ID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": "FileNotFound", "data": nil})
		return
	}
	ctx.File(path)
}

func UploadChallenge(ctx *gin.Context) {
	challenge := middleware.GetChallenge(ctx)
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	var path string
	switch challenge.Type {
	case model.StaticChallenge, model.DockerChallenge, model.DockersChallenge:
		if file.Filename != model.AttachmentFile {
			ctx.JSON(http.StatusOK, gin.H{"msg": "InvalidFileName", "data": nil})
			return
		}
		path = fmt.Sprintf("%s/%s", challenge.BasicDir(), model.AttachmentFile)
	case model.DynamicChallenge:
		if file.Filename != model.GeneratorFile {
			ctx.JSON(http.StatusOK, gin.H{"msg": "InvalidFileName", "data": nil})
			return
		}
		path = fmt.Sprintf("%s/%s", challenge.BasicDir(), model.GeneratorFile)
	default:
		ctx.JSON(http.StatusOK, gin.H{"msg": "InvalidChallengeType", "data": nil})
		return
	}
	if err := ctx.SaveUploadedFile(file, path); err != nil {
		log.Logger.Warningf("Failed to save file: %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "UnknownError", "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
}

func GetAvatars(ctx *gin.Context) {
	var form f.GetModelsForm
	if _, exists := ctx.GetQuery("limit"); !exists {
		form.Limit = 10
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		form.Offset = 0
	}
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	avatars, count, ok, msg := db.InitFileRepo(db.DB.WithContext(ctx)).GetAll("avatar", form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, avatar := range avatars {
		data = append(data, resp.GetFileResp(avatar))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "avatars": data}})
}

func DeleteAvatars(ctx *gin.Context) {
	var form f.DeleteFileForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	repo := db.InitFileRepo(db.DB.WithContext(ctx))
	for _, id := range form.FileIDL {
		if file, ok, _ := repo.GetByID(id); ok {
			_ = os.Remove(file.Path)
		}
	}
	_, _ = repo.Delete(form.FileIDL...)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
}

func GetWriteUPs(ctx *gin.Context) {
	contest := middleware.GetContest(ctx)
	team := middleware.GetTeam(ctx)
	path := fmt.Sprintf("%s/writeups/%d/%d", config.Env.Path, contest.ID, team.ID)
	dir, err := os.ReadDir(path)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": "FileNotFound", "data": nil})
		return
	}
	var files []string
	for _, file := range dir {
		files = append(files, file.Name())
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": files})
}
