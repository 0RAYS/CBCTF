package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetContainer(ctx *gin.Context) {
	docker, ok, msg := db.GetDockerBy3ID(ctx, middleware.GetContest(ctx).ID, middleware.GetTeam(ctx).ID, middleware.GetChallenge(ctx).ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"target": docker.RemoteAddr(), "remaining": docker.Remaining()}})
}

func StartContainer(ctx *gin.Context) {
	flag, ok, msg := db.GetFlagBy3ID(ctx, middleware.GetContest(ctx).ID, middleware.GetTeam(ctx).ID, middleware.GetChallenge(ctx).ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	docker, ok, msg := db.CreateDocker(ctx, flag, middleware.GetChallenge(ctx), middleware.GetSelfID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"target": docker.RemoteAddr(), "remaining": docker.Remaining()}})
}

func IncreaseDuration(ctx *gin.Context) {
	docker, ok, msg := db.GetDockerBy3ID(ctx, middleware.GetContest(ctx).ID, middleware.GetTeam(ctx).ID, middleware.GetChallenge(ctx).ID)
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
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"target": docker.RemoteAddr(), "remaining": docker.Remaining()}})
}

func StopContainer(ctx *gin.Context) {
	docker, ok, msg := db.GetDockerBy3ID(ctx, middleware.GetContest(ctx).ID, middleware.GetTeam(ctx).ID, middleware.GetChallenge(ctx).ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ok, msg = db.DeleteDocker(ctx, docker)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
}
