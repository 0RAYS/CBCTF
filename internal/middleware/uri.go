package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetUser(ctx *gin.Context) {
	type userIDUri struct {
		UserID uint `uri:"userID" binding:"required"`
	}
	var userID userIDUri
	if err := ctx.ShouldBindUri(&userID); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
	}
	user, ok, msg := db.GetUserByID(ctx, userID.UserID)
	if !ok {
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
	}
	ctx.Set("User", user)
	ctx.Next()
}

func GetUser(ctx *gin.Context) model.User {
	if user, ok := ctx.Get("User"); !ok {
		return model.User{}
	} else {
		return user.(model.User)
	}
}

func SetContest(ctx *gin.Context) {
	type contestIDUri struct {
		ContestID uint `uri:"contestID" binding:"required"`
	}
	var contestID contestIDUri
	if err := ctx.ShouldBindUri(&contestID); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
	}
	contest, ok, msg := db.GetContestByID(ctx, contestID.ContestID)
	if !ok {
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
	}
	ctx.Set("Contest", contest)
	ctx.Next()
}

func GetContest(ctx *gin.Context) model.Contest {
	if contest, ok := ctx.Get("Contest"); !ok {
		return model.Contest{}
	} else {
		return contest.(model.Contest)
	}
}

func SetTeamByUser(ctx *gin.Context) {
	var (
		self model.User
		team model.Team
		ok   bool
		msg  string
	)
	self, ok = GetSelf(ctx).(model.User)
	if !ok {
		ctx.JSONP(http.StatusForbidden, gin.H{"msg": "Forbidden", "data": nil})
		ctx.Abort()
	}
	team, ok, msg = db.GetTeamByUserID(ctx, self.ID, GetContest(ctx).ID)
	if !ok {
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
	}
	ctx.Set("Team", team)
	ctx.Next()
}

func SetTeamByURI(ctx *gin.Context) {
	var (
		team model.Team
		ok   bool
		msg  string
	)
	type teamIDUri struct {
		TeamID uint `uri:"teamID" binding:"required"`
	}
	var teamID teamIDUri
	if err := ctx.ShouldBindUri(&teamID); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
	}
	team, ok, msg = db.GetTeamByID(ctx, teamID.TeamID)
	if !ok {
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
	}
	ctx.Set("Team", team)
	ctx.Next()
}

func GetTeam(ctx *gin.Context) model.Team {
	if team, ok := ctx.Get("Team"); !ok {
		return model.Team{}
	} else {
		return team.(model.Team)
	}
}

func SetAvatar(ctx *gin.Context) {
	type avatarIDUri struct {
		AvatarID string `uri:"avatarID" binding:"required"`
	}
	var avatarID avatarIDUri
	if err := ctx.ShouldBindUri(&avatarID); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
	}
	avatar, ok, msg := db.GetAvatarByID(ctx, avatarID.AvatarID)
	if !ok {
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
	}
	ctx.Set("Avatar", avatar)
	ctx.Next()
}

func GetAvatar(ctx *gin.Context) model.Avatar {
	if avatar, ok := ctx.Get("Avatar"); !ok {
		return model.Avatar{}
	} else {
		return avatar.(model.Avatar)
	}
}

func SetChallenge(ctx *gin.Context) {
	type challengeIDUri struct {
		ChallengeID string `uri:"challengeID" binding:"required"`
	}
	var challengeID challengeIDUri
	if err := ctx.ShouldBindUri(&challengeID); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		ctx.Abort()
	}
	challenge, ok, msg := db.GetChallengeByID(ctx, challengeID.ChallengeID)
	if !ok {
		ctx.JSONP(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
	}
	ctx.Set("Challenge", challenge)
	ctx.Next()
}

func GetChallenge(ctx *gin.Context) model.Challenge {
	if challenge, ok := ctx.Get("Challenge"); !ok {
		return model.Challenge{}
	} else {
		return challenge.(model.Challenge)
	}
}
