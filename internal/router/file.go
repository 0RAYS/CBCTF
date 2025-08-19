package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func DownloadFile(eventType string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.Set(middleware.CTXEventTypeKey, eventType)
		file := middleware.GetFile(ctx)
		if _, err := os.Stat(file.Path); err != nil {
			if os.IsNotExist(err) {
				tx := db.DB.WithContext(ctx).Begin()
				if ok, _ := db.InitFileRepo(tx).Delete(file.ID); !ok {
					tx.Rollback()
				} else {
					tx.Commit()
				}
				ctx.JSON(http.StatusOK, gin.H{"msg": i18n.FileNotFound, "data": file.ID})
				return
			}
			log.Logger.Warningf("Failed to get file: %s", err)
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
			return
		}
		ctx.Set(middleware.CTXEventSuccessKey, true)
		ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Filename))
		ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
		ctx.File(file.Path)
	}
}

func DownloadChallengeFile(ctx *gin.Context) {
	var form f.DownloadChallengeForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.DownloadAttachmentEventType)
	challenge := middleware.GetChallenge(ctx)
	path := fmt.Sprintf("%s/%s", challenge.BasicDir(), form.File)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.FileNotFound, "data": nil})
			return
		}
		log.Logger.Warningf("Failed to get file: %s", err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.File(path)
}

// DownloadAttachment 需要预加载 Challenge
func DownloadAttachment(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DownloadAttachmentEventType)
	contestChallenge := middleware.GetContestChallenge(ctx)
	team := middleware.GetTeam(ctx)
	path := contestChallenge.Challenge.AttachmentPath(team.ID)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.FileNotFound, "data": nil})
			return
		}
		log.Logger.Warningf("Failed to get attachment: %s", err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.File(path)
}

func DownloadTraffic(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DownloadTrafficEventType)
	victim := middleware.GetVictim(ctx)
	if _, err := os.Stat(victim.TrafficZipPath()); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.FileNotFound, "data": nil})
			return
		}
		log.Logger.Warningf("Failed to get file: %s", err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.File(victim.TrafficZipPath())
}

func UploadAvatar(v string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile(model.AvatarFile)
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
			options.AdminID = &id
		case "self-user":
			id = middleware.GetSelfID(ctx)
			options.UserID = &id
		case "user":
			id = middleware.GetUser(ctx).ID
			selfID := middleware.GetSelfID(ctx)
			options.AdminID = &selfID
			options.UserID = &id
		case "contest":
			id = middleware.GetContest(ctx).ID
			selfID := middleware.GetSelfID(ctx)
			options.AdminID = &selfID
			options.ContestID = &id
		case "team":
			id = middleware.GetTeam(ctx).ID
			options.TeamID = &id
			selfID := middleware.GetSelfID(ctx)
			if middleware.IsAdmin(ctx) {
				options.AdminID = &selfID
			} else {
				options.UserID = &selfID
			}
		case "oauth":
			id = middleware.GetOauth(ctx).ID
			options.OauthID = &id
			selfID := middleware.GetSelfID(ctx)
			options.AdminID = &selfID
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
	file, err := ctx.FormFile(model.ChallengeFile)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UploadChallengeFileEventType)
	var path string
	challenge := middleware.GetChallenge(ctx)
	switch challenge.Type {
	case model.StaticChallengeType, model.QuestionChallengeType, model.PodsChallengeType:
		if file.Filename != model.AttachmentFile {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.InvalidFileName, "data": nil})
			return
		}
		path = fmt.Sprintf("%s/%s", challenge.BasicDir(), model.AttachmentFile)
	case model.DynamicChallengeType:
		if file.Filename != model.GeneratorFile {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.InvalidFileName, "data": nil})
			return
		}
		path = fmt.Sprintf("%s/%s", challenge.BasicDir(), model.GeneratorFile)
	default:
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.InvalidChallengeType, "data": nil})
		return
	}
	if err = ctx.SaveUploadedFile(file, path); err != nil {
		log.Logger.Warningf("Failed to save file: %s", err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
		return
	}
	prometheus.UpdateFileUploadMetrics(".zip", file.Size)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
}

func UploadWriteUp(ctx *gin.Context) {
	file, err := ctx.FormFile(model.WriteUPFile)
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

func GetAvatars(ctx *gin.Context) {
	var form f.GetModelsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	avatars, count, ok, msg := db.InitFileRepo(db.DB.WithContext(ctx)).List(form.Limit, form.Offset, db.GetOptions{
		Conditions: map[string]any{"type": model.AvatarFile},
	})
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

func GetWriteUPs(ctx *gin.Context) {
	var form f.GetModelsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	team := middleware.GetTeam(ctx)
	writeups, count, ok, msg := db.InitFileRepo(db.DB.WithContext(ctx)).
		List(form.Limit, form.Offset, db.GetOptions{
			Conditions: map[string]any{"type": model.WriteUPFile, "team_id": team.ID},
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

func DeleteAvatars(ctx *gin.Context) {
	var form f.DeleteFileForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteAvatarEventType)
	repo := db.InitFileRepo(db.DB.WithContext(ctx))
	// 保留文件
	//for _, id := range form.FileIDL {
	//	if file, ok, _ := db.GetByID(id); ok {
	//		_ = os.Remove(file.Path)
	//	}
	//}
	repo.DeleteByRandID(form.FileIDL...)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
}
