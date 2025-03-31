package router

import (
	"CBCTF/internel/config"
	"CBCTF/internel/middleware"
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func DownloadFile(ctx *gin.Context) {
	file := middleware.GetFile(ctx)
	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		tx := db.DB.WithContext(ctx).Begin()
		if ok, _ := db.InitFileRepo(tx).Delete(file.ID); !ok {
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

func UploadFile(v, t string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile(t)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
			return
		}
		tx := db.DB.WithContext(ctx).Begin()
		record, ok, msg := service.SaveFile(tx, middleware.GetSelfID(ctx), file, t)
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
		switch t {
		case "avatar":
			var id uint
			switch v {
			case "self-admin", "self-user":
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
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": path})
		default:
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		}
	}
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
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": &files})
}
