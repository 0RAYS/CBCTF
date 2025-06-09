package middleware

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SetUser 保存 model.User 至上下文
func SetUser(ctx *gin.Context) {
	type userIDUri struct {
		UserID uint `uri:"userID" binding:"required"`
	}
	var userID userIDUri
	if err := ctx.ShouldBindUri(&userID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		ctx.Abort()
		return
	}
	user, ok, msg := db.InitUserRepo(db.DB.WithContext(ctx)).GetByID(userID.UserID, "all")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("User", user)
	ctx.Next()
}

// GetUser 从上下文中获取 model.User
func GetUser(ctx *gin.Context) model.User {
	if user, ok := ctx.Get("User"); !ok || user == nil {
		return model.User{}
	} else {
		return user.(model.User)
	}
}

// SetContest 保存 model.Contest 至上下文
func SetContest(ctx *gin.Context) {
	type contestIDUri struct {
		ContestID uint `uri:"contestID" binding:"required"`
	}
	var contestID contestIDUri
	if err := ctx.ShouldBindUri(&contestID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		ctx.Abort()
		return
	}
	contest, ok, msg := db.InitContestRepo(db.DB.WithContext(ctx)).GetByID(contestID.ContestID, "all")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	if GetRole(ctx) != "admin" && contest.Hidden {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": i18n.ContestNotFound, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("Contest", contest)
	ctx.Next()
}

// GetContest 从上下文中获取 model.Contest
func GetContest(ctx *gin.Context) model.Contest {
	if contest, ok := ctx.Get("Contest"); !ok || contest == nil {
		return model.Contest{}
	} else {
		return contest.(model.Contest)
	}
}

// SetTeamByUser 依照 model.User model.Contest 保存 model.Team 至上下文, 调用前前文须设置 model.Contest
func SetTeamByUser(ctx *gin.Context) {
	var (
		self model.User
		team model.Team
		ok   bool
		msg  string
	)
	self, ok = GetSelf(ctx).(model.User)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{"msg": i18n.Forbidden, "data": nil})
		ctx.Abort()
		return
	}
	team, ok, msg = db.InitTeamRepo(db.DB.WithContext(ctx)).GetBy2ID(self.ID, GetContest(ctx).ID, "all")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("Team", team)
	ctx.Next()
}

// SetTeamByURI 保存 model.Team 至上下文
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
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		ctx.Abort()
		return
	}
	team, ok, msg = db.InitTeamRepo(db.DB.WithContext(ctx)).GetByID(teamID.TeamID, "all")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("Team", team)
	ctx.Next()
}

// GetTeam 从上下文中获取 model.Team
func GetTeam(ctx *gin.Context) model.Team {
	if team, ok := ctx.Get("Team"); !ok || team == nil {
		return model.Team{}
	} else {
		return team.(model.Team)
	}
}

// SetFile 保存 model.File 至上下文
func SetFile(t string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type fileIDUri struct {
			FileID string `uri:"fileID" binding:"required"`
		}
		var fileID fileIDUri
		if err := ctx.ShouldBindUri(&fileID); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
			ctx.Abort()
			return
		}
		file, ok, msg := db.InitFileRepo(db.DB.WithContext(ctx)).GetByRandID(fileID.FileID)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			ctx.Abort()
			return
		}
		if file.Type != t {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
			ctx.Abort()
			return
		}
		ctx.Set("File", file)
		ctx.Next()
	}
}

// GetFile 从上下文中获取 model.File
func GetFile(ctx *gin.Context) model.File {
	if file, ok := ctx.Get("File"); !ok || file == nil {
		return model.File{}
	} else {
		return file.(model.File)
	}
}

// SetNotice 保存 model.Notice 至上下文
func SetNotice(ctx *gin.Context) {
	type noticeIDUri struct {
		NoticeID uint `uri:"noticeID" binding:"required"`
	}
	var noticeID noticeIDUri
	if err := ctx.ShouldBindUri(&noticeID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		ctx.Abort()
		return
	}
	notice, ok, msg := db.InitNoticeRepo(db.DB.WithContext(ctx)).GetByID(noticeID.NoticeID, "all")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("Notice", notice)
	ctx.Next()
}

// GetNotice 从上下文中获取 model.Notice
func GetNotice(ctx *gin.Context) model.Notice {
	if notice, ok := ctx.Get("Notice"); !ok || notice == nil {
		return model.Notice{}
	} else {
		return notice.(model.Notice)
	}
}

// SetChallenge 保存 model.Challenge 至上下文
func SetChallenge(ctx *gin.Context) {
	type challengeIDUri struct {
		ChallengeID string `uri:"challengeID" binding:"required"`
	}
	var challengeID challengeIDUri
	if err := ctx.ShouldBindUri(&challengeID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		ctx.Abort()
		return
	}
	challenge, ok, msg := db.InitChallengeRepo(db.DB.WithContext(ctx)).GetByRandID(challengeID.ChallengeID, "all")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("Challenge", challenge)
	ctx.Next()
}

// GetChallenge 从上下文中获取 model.Challenge
func GetChallenge(ctx *gin.Context) model.Challenge {
	if challenge, ok := ctx.Get("Challenge"); !ok || challenge == nil {
		return model.Challenge{}
	} else {
		return challenge.(model.Challenge)
	}
}

func SetContestChallenge(ctx *gin.Context) {
	type challengeIDUri struct {
		ChallengeID string `uri:"challengeID" binding:"required"`
	}
	var challengeID challengeIDUri
	if err := ctx.ShouldBindUri(&challengeID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		ctx.Abort()
		return
	}
	challenge, ok, msg := db.InitChallengeRepo(db.DB.WithContext(ctx)).GetByRandID(challengeID.ChallengeID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	contestChallenge, ok, msg := db.InitContestChallengeRepo(db.DB.WithContext(ctx)).GetWithConditions(db.GetOptions{
		{Key: "contest_id", Value: GetContest(ctx).ID, Op: "and"},
		{Key: "challenge_id", Value: challenge.ID, Op: "and"},
	}, false, "all")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("ContestChallenge", contestChallenge)
	ctx.Next()
}

// GetContestChallenge 从上下文中获取 model.ContestChallenge
func GetContestChallenge(ctx *gin.Context) model.ContestChallenge {
	if contestChallenge, ok := ctx.Get("ContestChallenge"); !ok || contestChallenge == nil {
		return model.ContestChallenge{}
	} else {
		return contestChallenge.(model.ContestChallenge)
	}
}

func SetContestFlag(ctx *gin.Context) {
	type flagIDUri struct {
		FlagID uint `uri:"flagID" binding:"required"`
	}
	var flagID flagIDUri
	if err := ctx.ShouldBindUri(&flagID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		ctx.Abort()
		return
	}
	contestFlag, ok, msg := db.InitContestFlagRepo(db.DB.WithContext(ctx)).GetByID(flagID.FlagID, "all")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("ContestFlag", contestFlag)
	ctx.Next()
}

func GetContestFlag(ctx *gin.Context) model.ContestFlag {
	if contestFlag, ok := ctx.Get("ContestFlag"); !ok || contestFlag == nil {
		return model.ContestFlag{}
	} else {
		return contestFlag.(model.ContestFlag)
	}
}

func SetVictim(ctx *gin.Context) {
	type victimIDUri struct {
		VictimID uint `uri:"victimID" binding:"required"`
	}
	var victimID victimIDUri
	if err := ctx.ShouldBindUri(&victimID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		ctx.Abort()
		return
	}
	victim, ok, msg := db.InitVictimRepo(db.DB.WithContext(ctx)).GetWithConditions(db.GetOptions{
		{Key: "victim_id", Value: victimID.VictimID, Op: "and"},
	}, true, "all")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		ctx.Abort()
		return
	}
	ctx.Set("Victim", victim)
	ctx.Next()
}

func GetVictim(ctx *gin.Context) model.Victim {
	if victim, ok := ctx.Get("Victim"); !ok || victim == nil {
		return model.Victim{}
	} else {
		return victim.(model.Victim)
	}
}
