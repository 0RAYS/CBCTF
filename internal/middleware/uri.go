package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"

	"github.com/gin-gonic/gin"
)

func SetFullAccess(ctx *gin.Context) {
	ctx.Set("FullAccess", true)
	ctx.Next()
}

func IsFullAccess(ctx *gin.Context) bool {
	return ctx.GetBool("FullAccess")
}

// SetRole 保存 model.Role 至上下文
func SetRole(ctx *gin.Context) {
	type roleIDUri struct {
		RoleID uint `uri:"roleID" binding:"required"`
	}
	var roleID roleIDUri
	if err := ctx.ShouldBindUri(&roleID); err != nil {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	role, ret := db.InitRoleRepo(db.DB).GetByID(roleID.RoleID)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("Role", role)
	ctx.Next()
}

// GetRole 从上下文中获取 model.Role
func GetRole(ctx *gin.Context) model.Role {
	role, ok := ctx.Get("Role")
	if !ok || role == nil {
		return model.Role{}
	}
	return role.(model.Role)
}

// SetGroup 保存 model.Group 至上下文
func SetGroup(ctx *gin.Context) {
	type groupIDUri struct {
		GroupID uint `uri:"groupID" binding:"required"`
	}
	var groupID groupIDUri
	if err := ctx.ShouldBindUri(&groupID); err != nil {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	group, ret := db.InitGroupRepo(db.DB).GetByID(groupID.GroupID)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("Group", group)
	ctx.Next()
}

// GetGroup 从上下文中获取 model.Group
func GetGroup(ctx *gin.Context) model.Group {
	group, ok := ctx.Get("Group")
	if !ok || group == nil {
		return model.Group{}
	}
	return group.(model.Group)
}

// SetPermission 保存 model.Permission 至上下文
func SetPermission(ctx *gin.Context) {
	type permissionIDUri struct {
		PermissionID uint `uri:"permissionID" binding:"required"`
	}
	var permissionID permissionIDUri
	if err := ctx.ShouldBindUri(&permissionID); err != nil {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	permission, ret := db.InitPermissionRepo(db.DB).GetByID(permissionID.PermissionID)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("Permission", permission)
	ctx.Next()
}

// GetPermission 从上下文中获取 model.Permission
func GetPermission(ctx *gin.Context) model.Permission {
	permission, ok := ctx.Get("Permission")
	if !ok || permission == nil {
		return model.Permission{}
	}
	return permission.(model.Permission)
}

// SetUser 保存 model.User 至上下文
func SetUser(ctx *gin.Context) {
	type userIDUri struct {
		UserID uint `uri:"userID" binding:"required"`
	}
	var userID userIDUri
	if err := ctx.ShouldBindUri(&userID); err != nil {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	user, ret := db.InitUserRepo(db.DB).GetByID(userID.UserID)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("User", user)
	ctx.Next()
}

// GetUser 从上下文中获取 model.User
func GetUser(ctx *gin.Context) model.User {
	user, ok := ctx.Get("User")
	if !ok || user == nil {
		return model.User{}
	}
	return user.(model.User)
}

// SetContest 保存 model.Contest 至上下文
func SetContest(ctx *gin.Context) {
	type contestIDUri struct {
		ContestID uint `uri:"contestID" binding:"required"`
	}
	var contestID contestIDUri
	if err := ctx.ShouldBindUri(&contestID); err != nil {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	contest, ret := db.InitContestRepo(db.DB).GetByID(contestID.ContestID)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	if !IsFullAccess(ctx) && contest.Hidden {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Model.Contest.NotFound})
		return
	}
	ctx.Set("Contest", contest)
	ctx.Next()
}

// GetContest 从上下文中获取 model.Contest
func GetContest(ctx *gin.Context) model.Contest {
	contest, ok := ctx.Get("Contest")
	if !ok || contest == nil {
		return model.Contest{}
	}
	return contest.(model.Contest)
}

// SetTeam 保存 model.Team 至上下文
func SetTeam(ctx *gin.Context) {
	type teamIDUri struct {
		TeamID uint `uri:"teamID" binding:"required"`
	}
	var teamID teamIDUri
	if err := ctx.ShouldBindUri(&teamID); err != nil {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	team, ret := db.InitTeamRepo(db.DB).GetByID(teamID.TeamID)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("Team", team)
	ctx.Next()
}

// SetTeamByUser 依照 model.User model.Contest 保存 model.Team 至上下文, 调用前前文须设置 model.Contest
func SetTeamByUser(ctx *gin.Context) {
	team, ret := db.InitTeamRepo(db.DB).GetBy2ID(GetSelf(ctx).ID, GetContest(ctx).ID)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("Team", team)
	ctx.Next()
}

// GetTeam 从上下文中获取 model.Team
func GetTeam(ctx *gin.Context) model.Team {
	team, ok := ctx.Get("Team")
	if !ok || team == nil {
		return model.Team{}
	}
	return team.(model.Team)
}

// SetFile 保存 model.File 至上下文
func SetFile(t model.FileType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		type fileIDUri struct {
			FileID string `uri:"fileID" binding:"required"`
		}
		var fileID fileIDUri
		if err := ctx.ShouldBindUri(&fileID); err != nil {
			resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
			return
		}
		options := db.GetOptions{}
		if t != "" {
			options = db.GetOptions{Conditions: map[string]any{"type": t}}
		}
		file, ret := db.InitFileRepo(db.DB).GetByUniqueField("rand_id", fileID.FileID, options)
		if !ret.OK {
			resp.AbortJSON(ctx, ret)
			return
		}
		ctx.Set("File", file)
		ctx.Next()
	}
}

func SetChallengeFile(ctx *gin.Context) {
	challenge := GetChallenge(ctx)
	file, ret := db.InitFileRepo(db.DB).Get(db.GetOptions{
		Conditions: map[string]any{"model": model.ModelName(challenge), "model_id": challenge.ID, "type": model.ChallengeFileType}},
	)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("File", file)
	ctx.Next()
}

func SetTrafficFile(ctx *gin.Context) {
	file, ret := db.InitFileRepo(db.DB).Get(db.GetOptions{
		Conditions: map[string]any{"model": model.ModelName(GetVictim(ctx)), "model_id": GetVictim(ctx).ID, "type": model.TrafficFileType},
	})
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("File", file)
	ctx.Next()
}

func SetAttachmentFile(test bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		challenge := GetChallenge(ctx)
		if test && challenge.Type == model.DynamicChallengeType {
			if ret := service.GenTestAttachment(db.DB, challenge); !ret.OK {
				resp.AbortJSON(ctx, ret)
				return
			}
		}
		path := challenge.AttachmentPath(GetTeam(ctx).ID)
		record, ret := db.InitFileRepo(db.DB).Get(db.GetOptions{
			Conditions: map[string]any{"model": model.ModelName(challenge), "model_id": challenge.ID, "type": model.ChallengeFileType}},
		)
		if ret.OK && string(record.Path) == path {
			ctx.Set("File", record)
			ctx.Next()
			return
		}
		ctx.Set("File", model.File{Filename: "attachment.zip", Path: model.FilePath(path)})
		ctx.Next()
	}
}

// GetFile 从上下文中获取 model.File
func GetFile(ctx *gin.Context) model.File {
	file, ok := ctx.Get("File")
	if !ok || file == nil {
		return model.File{}
	}
	return file.(model.File)
}

// SetNotice 保存 model.Notice 至上下文
func SetNotice(ctx *gin.Context) {
	type noticeIDUri struct {
		NoticeID uint `uri:"noticeID" binding:"required"`
	}
	var noticeID noticeIDUri
	if err := ctx.ShouldBindUri(&noticeID); err != nil {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	notice, ret := db.InitNoticeRepo(db.DB).GetByID(noticeID.NoticeID)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("Notice", notice)
	ctx.Next()
}

// GetNotice 从上下文中获取 model.Notice
func GetNotice(ctx *gin.Context) model.Notice {
	notice, ok := ctx.Get("Notice")
	if !ok || notice == nil {
		return model.Notice{}
	}
	return notice.(model.Notice)
}

// SetChallenge 保存 model.Challenge 至上下文
func SetChallenge(ctx *gin.Context) {
	type challengeIDUri struct {
		ChallengeID string `uri:"challengeID" binding:"required"`
	}
	var challengeID challengeIDUri
	if err := ctx.ShouldBindUri(&challengeID); err != nil {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	challenge, ret := db.InitChallengeRepo(db.DB).GetByRandID(challengeID.ChallengeID)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("Challenge", challenge)
	ctx.Next()
}

// GetChallenge 从上下文中获取 model.Challenge
func GetChallenge(ctx *gin.Context) model.Challenge {
	challenge, ok := ctx.Get("Challenge")
	if !ok || challenge == nil {
		return model.Challenge{}
	}
	return challenge.(model.Challenge)
}

func SetContestChallenge(ctx *gin.Context) {
	type challengeIDUri struct {
		ChallengeID string `uri:"challengeID" binding:"required"`
	}
	var challengeID challengeIDUri
	if err := ctx.ShouldBindUri(&challengeID); err != nil {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	challenge, ret := db.InitChallengeRepo(db.DB).GetByRandID(challengeID.ChallengeID)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	contestChallenge, ret := db.InitContestChallengeRepo(db.DB).Get(db.GetOptions{
		Conditions: map[string]any{"challenge_id": challenge.ID, "contest_id": GetContest(ctx).ID},
	})
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("ContestChallenge", contestChallenge)
	ctx.Set("Challenge", challenge)
	ctx.Next()
}

// GetContestChallenge 从上下文中获取 model.ContestChallenge
func GetContestChallenge(ctx *gin.Context) model.ContestChallenge {
	contestChallenge, ok := ctx.Get("ContestChallenge")
	if !ok || contestChallenge == nil {
		return model.ContestChallenge{}
	}
	return contestChallenge.(model.ContestChallenge)
}

func SetContestFlag(ctx *gin.Context) {
	type flagIDUri struct {
		FlagID uint `uri:"flagID" binding:"required"`
	}
	var flagID flagIDUri
	if err := ctx.ShouldBindUri(&flagID); err != nil {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	contestFlag, ret := db.InitContestFlagRepo(db.DB).GetByID(flagID.FlagID)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("ContestFlag", contestFlag)
	ctx.Next()
}

func GetContestFlag(ctx *gin.Context) model.ContestFlag {
	contestFlag, ok := ctx.Get("ContestFlag")
	if !ok || contestFlag == nil {
		return model.ContestFlag{}
	}
	return contestFlag.(model.ContestFlag)
}

func SetVictim(ctx *gin.Context) {
	type victimIDUri struct {
		VictimID uint `uri:"victimID" binding:"required"`
	}
	var victimID victimIDUri
	if err := ctx.ShouldBindUri(&victimID); err != nil {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	victim, ret := db.InitVictimRepo(db.DB).GetByID(victimID.VictimID, db.GetOptions{Deleted: true})
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("Victim", victim)
	ctx.Next()
}

func GetVictim(ctx *gin.Context) model.Victim {
	victim, ok := ctx.Get("Victim")
	if !ok || victim == nil {
		return model.Victim{}
	}
	return victim.(model.Victim)
}

func SetCheat(ctx *gin.Context) {
	type cheatIDUri struct {
		CheatID uint `uri:"cheatID" binding:"required"`
	}
	var cheatID cheatIDUri
	if err := ctx.ShouldBindUri(&cheatID); err != nil {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	cheat, ret := db.InitCheatRepo(db.DB).GetByID(cheatID.CheatID)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("Cheat", cheat)
	ctx.Next()
}

func GetCheat(ctx *gin.Context) model.Cheat {
	cheat, ok := ctx.Get("Cheat")
	if !ok || cheat == nil {
		return model.Cheat{}
	}
	return cheat.(model.Cheat)
}

func SetOauth(ctx *gin.Context) {
	type oauthIDUri struct {
		OauthID uint `uri:"oauthID" binding:"required"`
	}
	var oauthID oauthIDUri
	if err := ctx.ShouldBindUri(&oauthID); err != nil {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	oauth, ret := db.InitOauthRepo(db.DB).GetByID(oauthID.OauthID)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("Oauth", oauth)
	ctx.Next()
}

func GetOauth(ctx *gin.Context) model.Oauth {
	oauth, ok := ctx.Get("Oauth")
	if !ok || oauth == nil {
		return model.Oauth{}
	}
	return oauth.(model.Oauth)
}

// SetOauthUri 不使用数据库查询, 只传递名称, 后续使用内存中的 map 进行获取
func SetOauthUri(ctx *gin.Context) {
	type oauthUri struct {
		OauthUri string `uri:"oauth" binding:"required"`
	}
	var oauth oauthUri
	if err := ctx.ShouldBindUri(&oauth); err != nil {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	ctx.Set("OauthUri", oauth.OauthUri)
	ctx.Next()
}

func GetOauthUri(ctx *gin.Context) string {
	oauth, ok := ctx.Get("OauthUri")
	if !ok || oauth == nil {
		return ""
	}
	return oauth.(string)
}

func SetSmtp(ctx *gin.Context) {
	type smtpIDUri struct {
		SmtpID uint `uri:"smtpID" binding:"required"`
	}
	var smtpID smtpIDUri
	if err := ctx.ShouldBindUri(&smtpID); err != nil {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	smtp, ret := db.InitSmtpRepo(db.DB).GetByID(smtpID.SmtpID)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("Smtp", smtp)
	ctx.Next()
}

func GetSmtp(ctx *gin.Context) model.Smtp {
	smtp, ok := ctx.Get("Smtp")
	if !ok || smtp == nil {
		return model.Smtp{}
	}
	return smtp.(model.Smtp)
}

func SetCronJob(ctx *gin.Context) {
	type cronJobIDUri struct {
		CronJobID uint `uri:"cronJobID" binding:"required"`
	}
	var cronJobID cronJobIDUri
	if err := ctx.ShouldBindUri(&cronJobID); err != nil {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	cronJob, ret := db.InitCronJobRepo(db.DB).GetByID(cronJobID.CronJobID)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("CronJob", cronJob)
	ctx.Next()
}

func GetCronJob(ctx *gin.Context) model.CronJob {
	cronJob, ok := ctx.Get("CronJob")
	if !ok || cronJob == nil {
		return model.CronJob{}
	}
	return cronJob.(model.CronJob)
}

func SetWebhook(ctx *gin.Context) {
	type webhookIDUti struct {
		WebhookID uint `uri:"webhookID" binding:"required"`
	}
	var webhookID webhookIDUti
	if err := ctx.ShouldBindUri(&webhookID); err != nil {
		resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	webhook, ret := db.InitWebhookRepo(db.DB).GetByID(webhookID.WebhookID)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("Webhook", webhook)
}

func GetWebhook(ctx *gin.Context) model.Webhook {
	webhook, ok := ctx.Get("Webhook")
	if !ok || webhook == nil {
		return model.Webhook{}
	}
	return webhook.(model.Webhook)
}
