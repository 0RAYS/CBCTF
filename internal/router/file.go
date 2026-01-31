package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"CBCTF/internal/task"
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
				//db.InitFileRepo(db.DB).Delete(file.ID)
				ctx.JSON(http.StatusNotFound, model.RetVal{Msg: i18n.Model.NotFound, Attr: map[string]any{"Model": file.ModelName()}})
				return
			}
			log.Logger.Warningf("Failed to get file: %s", err)
			ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err}})
			return
		}
		ctx.Set(middleware.CTXEventSuccessKey, true)
		ctx.FileAttachment(file.Path, file.Filename)
	}
}

func UploadPicture(v string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile(model.PictureFileType)
		if err != nil {
			ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Request.BadRequest})
			return
		}
		ctx.Set(middleware.CTXEventTypeKey, model.UploadPictureEventType)
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
		record, ret := service.SavePicture(db.DB, options, file)
		if !ret.OK {
			ctx.JSON(http.StatusOK, ret)
			return
		}
		path, ret := service.UpdatePicture(db.DB, v, id, record)
		if !ret.OK {
			ctx.JSON(http.StatusOK, ret)
			return
		}
		if err = ctx.SaveUploadedFile(file, record.Path); err != nil {
			log.Logger.Warningf("Failed to save file: %s", err)
			ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err}})
			return
		}
		if v != "contest" {
			_, _ = task.EnqueueResizeImageTask(record.Path, 100, 100)
		}
		path = fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(path, "/"))
		ctx.Set(middleware.CTXEventSuccessKey, true)
		ctx.JSON(http.StatusOK, model.SuccessRetVal(path))
	}
}

func UploadChallengeFile(ctx *gin.Context) {
	file, err := ctx.FormFile(model.ChallengeFileType)
	if err != nil {
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Request.BadRequest})
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
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Model.Challenge.InvalidType})
		return
	}
	record, ret := service.SaveChallengeFile(db.DB, challenge, file, path)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	if err = ctx.SaveUploadedFile(file, record.Path); err != nil {
		log.Logger.Warningf("Failed to save file: %s", err)
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err}})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal())
}

func UploadWriteUp(ctx *gin.Context) {
	file, err := ctx.FormFile(model.WriteupFileType)
	if err != nil {
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Request.BadRequest})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UploadWriteUpEventType)
	user := middleware.GetSelf(ctx).(model.User)
	contest := middleware.GetContest(ctx)
	team := middleware.GetTeam(ctx)
	record, ret := service.SaveWriteUp(db.DB, user, contest, team, file)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	if err = ctx.SaveUploadedFile(file, record.Path); err != nil {
		log.Logger.Warningf("Failed to save file: %s", err)
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err}})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, ret)
}

func GetFiles(ctx *gin.Context) {
	var form dto.GetFilesForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	options := db.GetOptions{}
	if form.Type != "" {
		options.Conditions = map[string]any{"type": form.Type}
	}
	files, count, ret := db.InitFileRepo(db.DB).List(form.Limit, form.Offset, options)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, file := range files {
		data = append(data, resp.GetFileResp(file))
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"count": count, "files": data}))
}

func GetWriteUPs(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	team := middleware.GetTeam(ctx)
	writeups, count, ret := db.InitFileRepo(db.DB).List(form.Limit, form.Offset, db.GetOptions{
		Conditions: map[string]any{"type": model.WriteupFileType, "team_id": team.ID},
	})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, writeup := range writeups {
		data = append(data, resp.GetFileResp(writeup))
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"count": count, "writeups": data}))
}

func DeleteFiles(ctx *gin.Context) {
	var form dto.DeleteFileForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.DeletePictureEventType)
	ret := db.InitFileRepo(db.DB).DeleteByRandID(form.FileIDL...)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal())
}
