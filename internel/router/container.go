package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/log"
	"CBCTF/internel/middleware"
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func StartContainer(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	usage := middleware.GetUsage(ctx)
	user := middleware.GetSelf(ctx).(model.User)
	if ok, err := redis.CheckContainersCreate(team.ID, usage.ChallengeID); ok || err != nil {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"msg": "TooQuick", "data": nil})
		return
	}
	if err := redis.RecordContainersCreate(team.ID, usage.ChallengeID); err != nil {
		log.Logger.Warningf("Failed to record container create: %v", err)
	}
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.CreateContainer(tx, user, team, usage)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}

}
func GetContainer(ctx *gin.Context) {
	container := middleware.GetContainer(ctx)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": &container})
}

func GetContainers(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	team := middleware.GetTeam(ctx)
	repo := db.InitContainerRepo(db.DB.WithContext(ctx))
	containers, count, ok, msg := repo.GetAll(team.ID, form.Limit, form.Offset, false, 0)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"containers": &containers, "count": count}})
}

func DownloadTraffic(ctx *gin.Context) {
	container := middleware.GetContainer(ctx)
	if _, err := os.Stat(container.TrafficPath()); err != nil {
		log.Logger.Warningf("Failed to get file: %s", err)
		if errors.Is(err, os.ErrNotExist) {
			ctx.JSON(http.StatusOK, gin.H{"msg": "FileNotFound", "data": nil})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"msg": "UnknownError", "data": nil})
		return
	}
	ctx.File(container.TrafficPath())
}

func LoadTraffic(ctx *gin.Context) {
	container := middleware.GetContainer(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.SaveTraffic(tx, container)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
}

func GetTraffics(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	container := middleware.GetContainer(ctx)
	repo := db.InitTrafficRepo(db.DB.WithContext(ctx))
	traffics, count, ok, msg := repo.GetAll(container.ID, form.Limit, form.Offset, false, 0)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"traffics": &traffics, "count": count}})
}
