package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetContainer(ctx *gin.Context) {
	team, ok, msg := db.GetTeamByUserID(ctx, middleware.GetSelfID(ctx), middleware.GetContestID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	docker, ok, msg := db.GetDockerBy3ID(ctx, middleware.GetContestID(ctx), team.ID, middleware.GetChallengeID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": docker.RemoteAddr()})
}

func StartContainer(ctx *gin.Context) {
	team, ok, msg := db.GetTeamByUserID(ctx, middleware.GetSelfID(ctx), middleware.GetContestID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	flag, ok, msg := db.GetFlagBy3ID(ctx, middleware.GetContestID(ctx), team.ID, middleware.GetChallengeID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	docker, ok, msg := db.CreateDocker(ctx, flag, middleware.GetSelfID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": docker.RemoteAddr()})
}

func IncreaseDuration(ctx *gin.Context) {
	team, ok, msg := db.GetTeamByUserID(ctx, middleware.GetSelfID(ctx), middleware.GetContestID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	docker, ok, msg := db.GetDockerBy3ID(ctx, middleware.GetContestID(ctx), team.ID, middleware.GetChallengeID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if !docker.Start.Add(docker.Duration).Before(time.Now().Add(20 * time.Second)) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "HasMuchTime", "data": nil})
		return
	}
	ok, msg = db.UpdateDocker(ctx, docker.ID, map[string]interface{}{"duration": docker.Duration + 1*time.Hour})
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": docker.RemoteAddr()})
}

func StopContainer(ctx *gin.Context) {
	team, ok, msg := db.GetTeamByUserID(ctx, middleware.GetSelfID(ctx), middleware.GetContestID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	docker, ok, msg := db.GetDockerBy3ID(ctx, middleware.GetContestID(ctx), team.ID, middleware.GetChallengeID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ok, msg = db.DeleteDocker(ctx, docker.ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
}
