package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetContainer(ctx *gin.Context) {
	docker, ok, msg := db.GetDockerBy3ID(db.DB.WithContext(ctx), middleware.GetContest(ctx).ID, middleware.GetTeam(ctx).ID, middleware.GetChallenge(ctx).ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"target": docker.RemoteAddr(), "remaining": docker.Remaining().Seconds()}})
}

func StartContainer(ctx *gin.Context) {
	flag, ok, msg := db.GetFlagBy3ID(db.DB.WithContext(ctx), middleware.GetContest(ctx).ID, middleware.GetTeam(ctx).ID, middleware.GetChallenge(ctx).ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	docker, ok, msg := db.CreateDocker(tx, flag, middleware.GetChallenge(ctx), middleware.GetSelfID(ctx))
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"target": docker.RemoteAddr(), "remaining": docker.Remaining().Seconds()}})
}

func IncreaseDuration(ctx *gin.Context) {
	docker, ok, msg := db.GetDockerBy3ID(db.DB.WithContext(ctx), middleware.GetContest(ctx).ID, middleware.GetTeam(ctx).ID, middleware.GetChallenge(ctx).ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if !docker.Start.Add(docker.Duration).Before(time.Now().Add(20 * time.Second)) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "HasMuchTime", "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg = db.UpdateDocker(tx, docker.ID, map[string]interface{}{"duration": docker.Duration + 1*time.Hour})
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"target": docker.RemoteAddr(), "remaining": docker.Remaining().Seconds()}})
}

func StopContainer(ctx *gin.Context) {
	docker, ok, msg := db.GetDockerBy3ID(db.DB.WithContext(ctx), middleware.GetContest(ctx).ID, middleware.GetTeam(ctx).ID, middleware.GetChallenge(ctx).ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg = db.DeleteDocker(tx, docker)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
}
