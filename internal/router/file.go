package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/oauth"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"CBCTF/internal/task"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var DefaultPicture = map[string][]byte{
	"github":  oauth.GithubMark,
	"hduhelp": oauth.HDUHelpPicture,
}

func DefaultAssets(ctx *gin.Context) {
	var form dto.GetAssetForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	file, ok := DefaultPicture[form.Filename]
	if !ok {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Model.File.NotFound})
		return
	}
	ctx.Writer.Header().Set("File", "true")
	ctx.Data(http.StatusOK, "application/octet-stream", file)
}

func DownloadFile(eventType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(middleware.CTXEventTypeKey, eventType)
		file := middleware.GetFile(ctx)
		if _, err := os.Stat(string(file.Path)); err != nil {
			if os.IsNotExist(err) {
				// 保留数据库记录
				//db.InitFileRepo(db.DB).Delete(file.ID)
				resp.JSON(ctx, model.RetVal{Msg: i18n.Model.File.NotFound})
				return
			}
			log.Logger.Warningf("Failed to get file: %s", err)
			resp.JSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
			return
		}
		ctx.Writer.Header().Set("File", "true")
		ctx.Set(middleware.CTXEventSuccessKey, true)
		ctx.FileAttachment(string(file.Path), file.Filename)
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
		var modelName string
		var id uint
		switch v {
		case "self":
			id = middleware.GetSelf(ctx).ID
			modelName = model.ModelName(model.User{})
		case "user":
			id = middleware.GetUser(ctx).ID
			modelName = model.ModelName(model.User{})
		case "contest":
			id = middleware.GetContest(ctx).ID
			modelName = model.ModelName(model.Contest{})
		case "team":
			id = middleware.GetTeam(ctx).ID
			modelName = model.ModelName(model.Team{})
		case "oauth":
			id = middleware.GetOauth(ctx).ID
			modelName = model.ModelName(model.Oauth{})
		case "branding":
			var ret model.RetVal
			id, ret = service.GetDefaultBrandingID(db.DB)
			if !ret.OK {
				resp.JSON(ctx, ret)
				return
			}
			modelName = model.ModelName(model.Branding{})
		}
		record, ret := service.SavePicture(db.DB, modelName, id, file)
		if !ret.OK {
			resp.JSON(ctx, ret)
			return
		}
		path, ret := service.UpdatePicture(db.DB, v, id, record)
		if !ret.OK {
			resp.JSON(ctx, ret)
			return
		}
		if err = ctx.SaveUploadedFile(file, string(record.Path)); err != nil {
			log.Logger.Warningf("Failed to save file: %s", err)
			resp.JSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
			return
		}
		if v != "contest" {
			_, _ = task.EnqueueResizeImageTask(string(record.Path), 100, 100)
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
	case model.StaticChallengeType, model.PodsChallengeType:
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
	if err = ctx.SaveUploadedFile(file, string(record.Path)); err != nil {
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
	if err = ctx.SaveUploadedFile(file, string(record.Path)); err != nil {
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
	files, count, ret := service.ListFiles(db.DB, form)
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
	writeups, count, ret := service.ListWriteUps(db.DB, middleware.GetTeam(ctx), form)
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
	ret := service.DeleteFiles(db.DB, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}
