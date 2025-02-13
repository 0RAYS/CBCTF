package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetContainer(deleted bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if deleted {
			docker := middleware.GetContainer(ctx)
			ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": docker})
			return
		}
		docker, ok, msg := db.GetDockerBy3ID(db.DB.WithContext(ctx), middleware.GetContest(ctx).ID, middleware.GetTeam(ctx).ID, middleware.GetChallenge(ctx).ID)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"target": docker.RemoteAddr(), "remaining": docker.Remaining().Seconds()}})
	}
}

func GetContainers(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	dockers, count, ok, msg := db.GetDockerByTeamID(db.DB.WithContext(ctx), middleware.GetTeam(ctx).ID, form.Limit, form.Offset, true)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"dockers": dockers, "count": count}})
}

func StartContainer(ctx *gin.Context) {
	var DB = db.DB.WithContext(ctx)
	if err := redis.RecordDockerCreate(middleware.GetTeam(ctx).ID, middleware.GetChallenge(ctx).ID); err != nil {
		log.Logger.Warningf("Failed to record docker create: %v", err)
	}
	flag, ok, msg := db.GetFlagBy3ID(DB, middleware.GetContest(ctx).ID, middleware.GetTeam(ctx).ID, middleware.GetChallenge(ctx).ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx := DB.Begin()
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
	var DB = db.DB.WithContext(ctx)
	docker, ok, msg := db.GetDockerBy3ID(DB, middleware.GetContest(ctx).ID, middleware.GetTeam(ctx).ID, middleware.GetChallenge(ctx).ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if !docker.Start.Add(docker.Duration).Before(time.Now().Add(20 * time.Second)) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "HasMuchTime", "data": nil})
		return
	}
	tx := DB.Begin()
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
	if middleware.GetRole(ctx) == "admin" {
		docker := middleware.GetContainer(ctx)
		if docker.DeletedAt.Valid {
			ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
			return
		}
		tx := db.DB.WithContext(ctx).Begin()
		ok, msg := db.DeleteDocker(tx, docker)
		if !ok {
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		tx.Commit()
		ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
		return
	} else {
		var DB = db.DB.WithContext(ctx)
		if ok, err := redis.CheckDockerCreate(middleware.GetTeam(ctx).ID, middleware.GetChallenge(ctx).ID); ok || err != nil {
			ctx.JSON(http.StatusTooManyRequests, gin.H{"msg": "TooQuick", "data": nil})
			return
		}
		docker, ok, msg := db.GetDockerBy3ID(DB, middleware.GetContest(ctx).ID, middleware.GetTeam(ctx).ID, middleware.GetChallenge(ctx).ID)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		tx := DB.Begin()
		ok, msg = db.DeleteDocker(tx, docker)
		if !ok {
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		tx.Commit()
		ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
	}
}
