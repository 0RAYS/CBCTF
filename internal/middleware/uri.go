package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetUserID(ctx *gin.Context) {
	type userIDUri struct {
		UserID uint `uri:"userID" binding:"required"`
	}
	var userID userIDUri
	if err := ctx.ShouldBindUri(&userID); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
	}
	ctx.Set("UserID", userID.UserID)
	ctx.Next()
}

func GetUserID(ctx *gin.Context) uint {
	if userID, ok := ctx.Get("UserID"); !ok {
		return 0
	} else {
		return userID.(uint)
	}
}

func SetContestID(ctx *gin.Context) {
	type contestIDUri struct {
		ContestID uint `uri:"contestID" binding:"required"`
	}
	var contestID contestIDUri
	if err := ctx.ShouldBindUri(&contestID); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
	}
	ctx.Set("ContestID", contestID.ContestID)
	ctx.Next()
}

func GetContestID(ctx *gin.Context) uint {
	if contestID, ok := ctx.Get("ContestID"); !ok {
		return 0
	} else {
		return contestID.(uint)
	}
}

func SetTeamID(ctx *gin.Context) {
	type teamIDUri struct {
		TeamID uint `uri:"teamID" binding:"required"`
	}
	var teamID teamIDUri
	if err := ctx.ShouldBindUri(&teamID); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
	}
	ctx.Set("TeamID", teamID.TeamID)
	ctx.Next()
}

func GetTeamID(ctx *gin.Context) uint {
	if teamID, ok := ctx.Get("TeamID"); !ok {
		return 0
	} else {
		return teamID.(uint)
	}
}

func SetAvatarID(ctx *gin.Context) {
	type avatarIDUri struct {
		AvatarID string `uri:"avatarID" binding:"required"`
	}
	var avatarID avatarIDUri
	if err := ctx.ShouldBindUri(&avatarID); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
	}
	ctx.Set("AvatarID", avatarID.AvatarID)
	ctx.Next()
}

func GetAvatarID(ctx *gin.Context) string {
	if avatarID, ok := ctx.Get("AvatarID"); !ok {
		return ""
	} else {
		return avatarID.(string)
	}
}

func SetChallengeID(ctx *gin.Context) {
	type challengeIDUri struct {
		ChallengeID uint `uri:"challengeID" binding:"required"`
	}
	var challengeID challengeIDUri
	if err := ctx.ShouldBindUri(&challengeID); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
	}
	ctx.Set("ChallengeID", challengeID.ChallengeID)
	ctx.Next()
}

func GetChallengeID(ctx *gin.Context) uint {
	if challengeID, ok := ctx.Get("ChallengeID"); !ok {
		return 0
	} else {
		return challengeID.(uint)
	}
}
