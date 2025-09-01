package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func DownloadFile(eventType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(middleware.CTXEventTypeKey, eventType)
		file := middleware.GetFile(ctx)
		if _, err := os.Stat(file.Path); err != nil {
			if os.IsNotExist(err) {
				// 保留数据库记录
				//tx := db.DB.WithContext(ctx).Begin()
				//if ok, _ := db.InitFileRepo(tx).Delete(file.ID); !ok {
				//	tx.Rollback()
				//} else {
				//	tx.Commit()
				//}
				ctx.JSON(http.StatusOK, gin.H{"msg": i18n.FileNotFound, "data": nil})
				return
			}
			log.Logger.Warningf("Failed to get file: %s", err)
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
			return
		}
		ctx.Set(middleware.CTXEventSuccessKey, true)
		ctx.FileAttachment(file.Path, file.Filename)
	}
}

func UploadAvatar(v string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile(model.AvatarFileType)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
			return
		}
		ctx.Set(middleware.CTXEventTypeKey, model.UploadAvatarEventType)
		options := db.CreateFileOptions{}
		var id uint
		switch v {
		case "admin":
			id = middleware.GetSelfID(ctx)
			options.AdminID = sql.Null[uint]{V: id, Valid: true}
		case "self-user":
			id = middleware.GetSelfID(ctx)
			options.UserID = sql.Null[uint]{V: id, Valid: true}
		case "user":
			id = middleware.GetUser(ctx).ID
			selfID := middleware.GetSelfID(ctx)
			options.AdminID = sql.Null[uint]{V: selfID, Valid: true}
			options.UserID = sql.Null[uint]{V: id, Valid: true}
		case "contest":
			id = middleware.GetContest(ctx).ID
			selfID := middleware.GetSelfID(ctx)
			options.AdminID = sql.Null[uint]{V: selfID, Valid: true}
			options.ContestID = sql.Null[uint]{V: id, Valid: true}
		case "team":
			id = middleware.GetTeam(ctx).ID
			options.TeamID = sql.Null[uint]{V: id, Valid: true}
			selfID := middleware.GetSelfID(ctx)
			if middleware.IsAdmin(ctx) {
				options.AdminID = sql.Null[uint]{V: selfID, Valid: true}
			} else {
				options.UserID = sql.Null[uint]{V: selfID, Valid: true}
			}
		case "oauth":
			id = middleware.GetOauth(ctx).ID
			options.OauthID = sql.Null[uint]{V: id, Valid: true}
			selfID := middleware.GetSelfID(ctx)
			options.AdminID = sql.Null[uint]{V: selfID, Valid: true}
		}
		tx := db.DB.WithContext(ctx).Begin()
		record, ok, msg := service.SaveAvatar(tx, options, file)
		if !ok {
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		if err = ctx.SaveUploadedFile(file, record.Path); err != nil {
			log.Logger.Warningf("Failed to save file: %s", err)
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
			return
		}
		path, ok, msg := service.UpdateAvatar(tx, v, id, record)
		if !ok {
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		tx.Commit()
		path = fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(path, "/"))
		ctx.Set(middleware.CTXEventSuccessKey, true)
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": path})
	}
}

func UploadChallengeFile(ctx *gin.Context) {
	file, err := ctx.FormFile(model.ChallengeFileType)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UploadChallengeFileEventType)
	var path string
	challenge := middleware.GetChallenge(ctx)
	switch challenge.Type {
	case model.StaticChallengeType, model.QuestionChallengeType, model.PodsChallengeType:
		path = fmt.Sprintf("%s/%s", challenge.BasicDir(), model.AttachmentFile)
	case model.DynamicChallengeType:
		path = fmt.Sprintf("%s/%s", challenge.BasicDir(), model.GeneratorFile)
	default:
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.InvalidChallengeType, "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	record, ok, msg := service.SaveChallengeFile(tx, challenge, file, path)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if err = ctx.SaveUploadedFile(file, record.Path); err != nil {
		tx.Rollback()
		log.Logger.Warningf("Failed to save file: %s", err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
		return
	}
	tx.Commit()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
}

func UploadWriteUp(ctx *gin.Context) {
	file, err := ctx.FormFile(model.WriteUPFileType)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UploadWriteUpEventType)
	user := middleware.GetSelf(ctx).(model.User)
	contest := middleware.GetContest(ctx)
	team := middleware.GetTeam(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	record, ok, msg := service.SaveWriteUp(tx, user, contest, team, file)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if err = ctx.SaveUploadedFile(file, record.Path); err != nil {
		tx.Rollback()
		log.Logger.Warningf("Failed to save file: %s", err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
		return
	}
	tx.Commit()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func GetFiles(ctx *gin.Context) {
	var form f.GetFilesForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	options := db.GetOptions{}
	if form.Type != "" {
		options.Conditions = map[string]any{"type": form.Type}
	}
	files, count, ok, msg := db.InitFileRepo(db.DB.WithContext(ctx)).List(form.Limit, form.Offset, options)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, file := range files {
		data = append(data, resp.GetFileResp(file))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "files": data}})
}

func GetWriteUPs(ctx *gin.Context) {
	var form f.GetModelsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	team := middleware.GetTeam(ctx)
	writeups, count, ok, msg := db.InitFileRepo(db.DB.WithContext(ctx)).
		List(form.Limit, form.Offset, db.GetOptions{
			Conditions: map[string]any{"type": model.WriteUPFileType, "team_id": team.ID},
		})
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, writeup := range writeups {
		data = append(data, resp.GetFileResp(writeup))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "writeups": data}})
}

func DeleteFiles(ctx *gin.Context) {
	var form f.DeleteFileForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteAvatarEventType)
	ok, msg := db.InitFileRepo(db.DB.WithContext(ctx)).DeleteByRandID(form.FileIDL...)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
}
