package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SetUser 设置 user 对象，判断 user 是否存在，并保存，供后续使用，减少数据库操作
func SetUser(ctx *gin.Context) {
	type userIDUri struct {
		UserID uint `uri:"userID" binding:"required"`
	}
	var userID userIDUri
	trace := GetTraceID(ctx)
	if err := ctx.ShouldBindUri(&userID); err != nil {
		log.Logger.Infof("| %s | %s", trace, err)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		ctx.Abort()
		return
	}
	user, ok, msg := db.GetUserByID(userID.UserID)
	if !ok {
		log.Logger.Infof("| %s | %s", trace, msg)
		ctx.JSONP(http.StatusNotFound, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("user", user)
	ctx.Next()
}

// SetTeam 设置 team 对象，判断 team 是否存在，并保存，供后续使用，减少数据库操作
func SetTeam(ctx *gin.Context) {
	type teamIDUri struct {
		TeamID uint `uri:"teamID" binding:"required"`
	}
	var teamID teamIDUri
	trace := GetTraceID(ctx)
	if err := ctx.ShouldBindUri(&teamID); err != nil {
		log.Logger.Infof("| %s | %s", trace, err)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		ctx.Abort()
		return
	}
	team, ok, msg := db.GetTeamByID(teamID.TeamID)
	if !ok {
		log.Logger.Infof("| %s | %s", trace, msg)
		ctx.JSONP(http.StatusNotFound, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("team", team)
	ctx.Next()

}

// SetContest 设置 contest 对象，判断 contest 是否存在，并保存，供后续使用，减少数据库操作
func SetContest(ctx *gin.Context) {
	type contestIDUri struct {
		ContestID uint `uri:"contestID" binding:"required"`
	}
	var contestID contestIDUri
	trace := GetTraceID(ctx)
	if err := ctx.ShouldBindUri(&contestID); err != nil {
		log.Logger.Infof("| %s | %s", trace, err)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		ctx.Abort()
		return
	}
	contest, ok, msg := db.GetContestByID(contestID.ContestID)
	if !ok {
		log.Logger.Infof("| %s | %s", trace, msg)
		ctx.JSONP(http.StatusNotFound, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("contest", contest)
	ctx.Next()
}
