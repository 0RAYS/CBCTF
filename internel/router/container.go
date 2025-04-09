package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/log"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	"CBCTF/internel/redis"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"CBCTF/internel/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
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
	containers, ok, msg := service.StartContainer(tx, user, team, usage)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	usage.Containers = containers
	status := service.GetRemoteStatus(tx, usage)
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": status})
}

func IncreaseDuration(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	usage := middleware.GetUsage(ctx)
	DB := db.DB.WithContext(ctx)
	repo := db.InitContainerRepo(DB)
	containers, ok, msg := repo.GetBy2ID(team.ID, usage.ID, false)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	for _, container := range containers {
		if !container.Start.Add(container.Duration).Before(time.Now().Add(20 * time.Minute)) {
			ctx.JSON(http.StatusOK, gin.H{"msg": "HasMuchTime", "data": nil})
			return
		}
	}
	tx := DB.Begin()
	repo = db.InitContainerRepo(tx)
	data := make([]gin.H, 0)
	for _, container := range containers {
		duration := container.Duration + 1*time.Hour
		ok, msg := repo.Update(container.ID, db.UpdateContainerOptions{
			Duration: &duration,
		})
		if !ok {
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		data = append(data, gin.H{
			"target":    container.RemoteAddr(),
			"remaining": container.Remaining().Seconds(),
			"status":    "Running",
		})
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": data})
}

func StopContainer(ctx *gin.Context) {
	DB := db.DB.WithContext(ctx)
	team := middleware.GetTeam(ctx)
	usage := middleware.GetUsage(ctx)
	repo := db.InitContainerRepo(DB)
	containers, ok, msg := repo.GetBy2ID(team.ID, usage.ID, false)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	for _, container := range containers {
		tx := DB.Begin()
		ok, msg = service.StopContainer(tx, container)
		if !ok {
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
}

func GetContainer(ctx *gin.Context) {
	container := middleware.GetContainer(ctx)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": resp.GetContainerResp(container)})
}

func GetContainers(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	if _, exists := ctx.GetQuery("limit"); !exists {
		form.Limit = 5
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		form.Offset = 0
	}
	team := middleware.GetTeam(ctx)
	repo := db.InitContainerRepo(db.DB.WithContext(ctx))
	containers, count, ok, msg := repo.GetByTeam(team.ID, form.Limit, form.Offset, true)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, container := range containers {
		data = append(data, resp.GetContainerResp(container))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"containers": data, "count": count}})
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
	if _, exists := ctx.GetQuery("limit"); !exists {
		form.Limit = 5
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		form.Offset = 0
	}
	container := middleware.GetContainer(ctx)
	repo := db.InitTrafficRepo(db.DB.WithContext(ctx))
	traffics, count, ok, msg := repo.GetAll(container.ID, form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"traffics": &traffics, "count": count}})
}
