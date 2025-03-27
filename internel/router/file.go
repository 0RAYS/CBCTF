package router

import (
	"CBCTF/internel/middleware"
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

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
