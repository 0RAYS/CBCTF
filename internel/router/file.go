package router

import (
	"CBCTF/internel/config"
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"CBCTF/internel/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func DownloadFile(ctx *gin.Context) {
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
	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Filename))
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	ctx.File(file.Path)
}

func DownloadChallengeFile(ctx *gin.Context) {
	var form f.DownloadChallengeForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
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
		if os.IsNotExist(err) {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.FileNotFound, "data": nil})
			return
		}
		log.Logger.Warningf("Failed to get file: %s", err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
		return
	}
	ctx.File(path)
}

// DownloadAttachment 需要预加载 Challenge
func DownloadAttachment(ctx *gin.Context) {
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
	ctx.File(path)
}

func UploadAvatar(v string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile(model.AvatarFile)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
			return
		}
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
			if middleware.GetRole(ctx) == "admin" {
				options.AdminID = &selfID
			} else {
				options.UserID = &selfID
			}
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
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": path})
	}
}

func UploadChallengeFile(ctx *gin.Context) {
	challenge := middleware.GetChallenge(ctx)
	file, err := ctx.FormFile(model.ChallengeFile)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	var path string
	switch challenge.Type {
	case model.StaticChallengeType, model.PodsChallengeType:
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
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
}

func UploadWriteUp(ctx *gin.Context) {
	file, err := ctx.FormFile(model.WriteUPFile)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
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
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
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
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	avatars, count, ok, msg := db.InitFileRepo(db.DB.WithContext(ctx)).ListWithConditions(form.Limit, form.Offset, db.GetOptions{
		{Key: "type", Value: model.AvatarFile, Op: "and"},
	}, false)
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
	if _, exists := ctx.GetQuery("limit"); !exists {
		form.Limit = 10
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		form.Offset = 0
	}
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	team := middleware.GetTeam(ctx)
	writeups, count, ok, msg := db.InitFileRepo(db.DB.WithContext(ctx)).
		ListWithConditions(form.Limit, form.Offset, db.GetOptions{
			{Key: "type", Value: model.WriteUPFile, Op: "and"},
			{Key: "team_id", Value: team.ID, Op: "or"},
		}, false)
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
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	repo := db.InitFileRepo(db.DB.WithContext(ctx))
	// 保留文件
	//for _, id := range form.FileIDL {
	//	if file, ok, _ := repo.GetByID(id); ok {
	//		_ = os.Remove(file.Path)
	//	}
	//}
	repo.DeleteByRandID(form.FileIDL...)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
}
