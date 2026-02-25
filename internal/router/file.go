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
				ctx.JSON(http.StatusNotFound, model.RetVal{Msg: i18n.Model.File.NotFound})
				return
			}
			log.Logger.Warningf("Failed to get file: %s", err)
			resp.JSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
			return
		}
		ctx.Set(middleware.CTXEventSuccessKey, true)
		ctx.FileAttachment(file.Path, file.Filename)
	}
}

func UploadPicture(v string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile(string(model.PictureFileType))
		if err != nil {
			resp.JSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
			return
		}
		ctx.Set(middleware.CTXEventTypeKey, model.UploadPictureEventType)
		options := db.CreateFileOptions{}
		var id uint
		switch v {
		case "self":
			id = middleware.GetSelf(ctx).ID
			options.Model = model.User{}.ModelName()
			options.ModelID = id
		case "user":
			id = middleware.GetUser(ctx).ID
			options.Model = model.User{}.ModelName()
			options.ModelID = id
		case "contest":
			id = middleware.GetContest(ctx).ID
			options.Model = model.Contest{}.ModelName()
			options.ModelID = id
		case "team":
			id = middleware.GetTeam(ctx).ID
			options.Model = model.Team{}.ModelName()
			options.ModelID = id
		case "oauth":
			id = middleware.GetOauth(ctx).ID
			options.Model = model.Oauth{}.ModelName()
			options.ModelID = id
		}
		record, ret := service.SavePicture(db.DB, options, file)
		if !ret.OK {
			resp.JSON(ctx, ret)
			return
		}
		path, ret := service.UpdatePicture(db.DB, v, id, record)
		if !ret.OK {
			resp.JSON(ctx, ret)
			return
		}
		if err = ctx.SaveUploadedFile(file, record.Path); err != nil {
			log.Logger.Warningf("Failed to save file: %s", err)
			resp.JSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
			return
		}
		if v != "contest" {
			_, _ = task.EnqueueResizeImageTask(record.Path, 100, 100)
		}
		path = fmt.Sprintf("%s/%s", config.Env.Host, strings.TrimPrefix(path, "/"))
		ctx.Set(middleware.CTXEventSuccessKey, true)
		resp.JSON(ctx, model.SuccessRetVal(path))
	}
}

func UploadChallengeFile(ctx *gin.Context) {
	file, err := ctx.FormFile(string(model.ChallengeFileType))
	if err != nil {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UploadChallengeFileEventType)
	var path string
	challenge := middleware.GetChallenge(ctx)
	switch challenge.Type {
	case model.StaticChallengeType, model.QuestionChallengeType, model.PodsChallengeType:
		path = fmt.Sprintf("%s/%s", challenge.BasicDir(), model.AttachmentFileName)
	case model.DynamicChallengeType:
		path = fmt.Sprintf("%s/%s", challenge.BasicDir(), model.GeneratorFileName)
	default:
		resp.JSON(ctx, model.RetVal{Msg: i18n.Model.Challenge.InvalidType})
		return
	}
	record, ret := service.SaveChallengeFile(db.DB, challenge, file, path)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	if err = ctx.SaveUploadedFile(file, record.Path); err != nil {
		log.Logger.Warningf("Failed to save file: %s", err)
		resp.JSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}

func UploadWriteUp(ctx *gin.Context) {
	file, err := ctx.FormFile(string(model.WriteupFileType))
	if err != nil {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UploadWriteUpEventType)
	contest := middleware.GetContest(ctx)
	team := middleware.GetTeam(ctx)
	record, ret := service.SaveWriteUp(db.DB, contest, team, file)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	if err = ctx.SaveUploadedFile(file, record.Path); err != nil {
		log.Logger.Warningf("Failed to save file: %s", err)
		resp.JSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, ret)
}

func GetFiles(ctx *gin.Context) {
	var form dto.GetFilesForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	options := db.GetOptions{}
	if form.Type != "" {
		options.Conditions = map[string]any{"type": form.Type}
	}
	files, count, ret := db.InitFileRepo(db.DB).List(form.Limit, form.Offset, options)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, file := range files {
		data = append(data, resp.GetFileResp(file))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"count": count, "files": data}))
}

func GetWriteUPs(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	team := middleware.GetTeam(ctx)
	writeups, count, ret := db.InitFileRepo(db.DB).List(form.Limit, form.Offset, db.GetOptions{
		Conditions: map[string]any{"model": team.ModelName(), "model_id": team.ID, "type": model.WriteupFileType},
	})
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, writeup := range writeups {
		data = append(data, resp.GetFileResp(writeup))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"count": count, "writeups": data}))
}

func DeleteFiles(ctx *gin.Context) {
	var form dto.DeleteFileForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.DeletePictureEventType)
	ret := db.InitFileRepo(db.DB).DeleteByRandID(form.FileIDs...)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}
