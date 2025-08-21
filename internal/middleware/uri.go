package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetUser 保存 model.User 至上下文
func SetUser(ctx *gin.Context) {
	type userIDUri struct {
		UserID uint `uri:"userID" binding:"required"`
	}
	var userID userIDUri
	if err := ctx.ShouldBindUri(&userID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	user, ok, msg := db.InitUserRepo(db.DB.WithContext(ctx)).GetByID(userID.UserID)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
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
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	contest, ok, msg := db.InitContestRepo(db.DB.WithContext(ctx)).GetByID(contestID.ContestID)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if !IsAdmin(ctx) && contest.Hidden {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.ContestNotFound, "data": nil})
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

// SetTeam 保存 model.Team 至上下文
func SetTeam(ctx *gin.Context) {
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
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	team, ok, msg = db.InitTeamRepo(db.DB.WithContext(ctx)).GetByID(teamID.TeamID)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set("Team", team)
	ctx.Next()
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
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.Forbidden, "data": nil})
		return
	}
	team, ok, msg = db.InitTeamRepo(db.DB.WithContext(ctx)).GetBy2ID(self.ID, GetContest(ctx).ID)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
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
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
			return
		}
		file, ok, msg := db.InitFileRepo(db.DB.WithContext(ctx)).GetByRandID(fileID.FileID)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		if file.Type != t {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
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
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	notice, ok, msg := db.InitNoticeRepo(db.DB.WithContext(ctx)).GetByID(noticeID.NoticeID)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
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
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	challenge, ok, msg := db.InitChallengeRepo(db.DB.WithContext(ctx)).GetByRandID(challengeID.ChallengeID)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
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
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	challenge, ok, msg := db.InitChallengeRepo(db.DB.WithContext(ctx)).GetByRandID(challengeID.ChallengeID)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	contestChallenge, ok, msg := db.InitContestChallengeRepo(db.DB.WithContext(ctx)).Get(db.GetOptions{
		Conditions: map[string]any{"challenge_id": challenge.ID, "contest_id": GetContest(ctx).ID},
	})
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set("ContestChallenge", contestChallenge)
	ctx.Set("Challenge", challenge)
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
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	contestFlag, ok, msg := db.InitContestFlagRepo(db.DB.WithContext(ctx)).GetByID(flagID.FlagID)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
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
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	victim, ok, msg := db.InitVictimRepo(db.DB.WithContext(ctx)).GetByID(victimID.VictimID, db.GetOptions{Deleted: true})
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
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

func SetCheat(ctx *gin.Context) {
	type cheatIDUri struct {
		CheatID uint `uri:"cheatID" binding:"required"`
	}
	var cheatID cheatIDUri
	if err := ctx.ShouldBindUri(&cheatID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	cheat, ok, msg := db.InitCheatRepo(db.DB.WithContext(ctx)).GetByID(cheatID.CheatID)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set("Cheat", cheat)
	ctx.Next()
}

func GetCheat(ctx *gin.Context) model.Cheat {
	if cheat, ok := ctx.Get("Cheat"); !ok || cheat == nil {
		return model.Cheat{}
	} else {
		return cheat.(model.Cheat)
	}
}

func SetOauth(ctx *gin.Context) {
	type oauthIDUri struct {
		OauthID uint `uri:"oauthID" binding:"required"`
	}
	var oauthID oauthIDUri
	if err := ctx.ShouldBindUri(&oauthID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	oauth, ok, msg := db.InitOauthRepo(db.DB.WithContext(ctx)).GetByID(oauthID.OauthID)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set("Oauth", oauth)
	ctx.Next()
}

func GetOauth(ctx *gin.Context) model.Oauth {
	if oauth, ok := ctx.Get("Oauth"); !ok || oauth == nil {
		return model.Oauth{}
	} else {
		return oauth.(model.Oauth)
	}
}

// SetOauthUri 不使用数据库查询, 只传递名称, 后续使用内存中的 map 进行获取
func SetOauthUri(ctx *gin.Context) {
	type oauthUri struct {
		OauthUri string `uri:"oauth" binding:"required"`
	}
	var oauth oauthUri
	if err := ctx.ShouldBindUri(&oauth); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	ctx.Set("OauthUri", oauth.OauthUri)
	ctx.Next()
}

func GetOauthUri(ctx *gin.Context) string {
	if oauth, ok := ctx.Get("OauthUri"); !ok || oauth == nil {
		return ""
	} else {
		return oauth.(string)
	}
}

func SetSmtp(ctx *gin.Context) {
	type smtpIDUri struct {
		SmtpID uint `uri:"smtpID" binding:"required"`
	}
	var smtpID smtpIDUri
	if err := ctx.ShouldBindUri(&smtpID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	smtp, ok, msg := db.InitSmtpRepo(db.DB.WithContext(ctx)).GetByID(smtpID.SmtpID)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set("Smtp", smtp)
}

func GetSmtp(ctx *gin.Context) model.Smtp {
	if smtp, ok := ctx.Get("Smtp"); !ok || smtp == nil {
		return model.Smtp{}
	} else {
		return smtp.(model.Smtp)
	}
}

func SetWebhook(ctx *gin.Context) {
	type webhookIDUti struct {
		WebhookID uint `uri:"webhookID" binding:"required"`
	}
	var webhookID webhookIDUti
	if err := ctx.ShouldBindUri(&webhookID); err != nil {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	webhook, ok, msg := db.InitWebhookRepo(db.DB.WithContext(ctx)).GetByID(webhookID.WebhookID)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set("Webhook", webhook)
}

func GetWebhook(ctx *gin.Context) model.Webhook {
	if webhook, ok := ctx.Get("Webhook"); !ok || webhook == nil {
		return model.Webhook{}
	} else {
		return webhook.(model.Webhook)
	}
}
