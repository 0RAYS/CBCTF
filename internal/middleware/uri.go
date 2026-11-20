package middleware

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
)

// byIDLookup is the shape every by-id repository lookup shares.
type byIDLookup[T any] interface {
	GetByID(id uint, options ...db.GetOptions) (T, model.RetVal)
}

// byRandIDLookup is the shape every by-random-id repository lookup shares.
type byRandIDLookup[T any] interface {
	GetByRandID(randID string, options ...db.GetOptions) (T, model.RetVal)
}

// getResource returns the resource a middleware stored under key, or the zero
// value when the middleware did not run.
func getResource[T any](ctx *gin.Context, key string) T {
	value, ok := ctx.Get(key)
	if !ok || value == nil {
		var zero T
		return zero
	}
	return value.(T)
}

// uriUint reads a required numeric path parameter.
func uriUint(ctx *gin.Context, param, kind, targetKey string) (uint, bool) {
	raw := ctx.Param(param)
	id, err := strconv.ParseUint(raw, 10, strconv.IntSize)
	if err != nil {
		resp.AbortJSON(ctx, model.FailRetVal(i18n.Response.BadRequest, kind, targetKey+"="+raw, err))
		return 0, false
	}
	if id == 0 {
		resp.AbortJSON(ctx, model.FailRetVal(i18n.Response.BadRequest, kind, targetKey+"="+raw, fmt.Errorf("path parameter %s must not be zero", param)))
		return 0, false
	}
	return uint(id), true
}

// uriString reads a required path parameter.
func uriString(ctx *gin.Context, param, kind, targetKey string) (string, bool) {
	raw := ctx.Param(param)
	if raw == "" {
		resp.AbortJSON(ctx, model.FailRetVal(i18n.Response.BadRequest, kind, targetKey+"="+raw, fmt.Errorf("missing %s path parameter", param)))
		return "", false
	}
	return raw, true
}

// loadByID loads the resource addressed by a numeric path parameter and stores
// it in the context under key.
func loadByID[T any](ctx *gin.Context, key, param, kind, targetKey string, repo byIDLookup[T], options ...db.GetOptions) bool {
	id, ok := uriUint(ctx, param, kind, targetKey)
	if !ok {
		return false
	}
	value, ret := repo.GetByID(id, options...)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return false
	}
	ctx.Set(key, value)
	return true
}

// loadByRandID loads the resource addressed by a random id path parameter and
// stores it in the context under key.
func loadByRandID[T any](ctx *gin.Context, key, param, kind, targetKey string, repo byRandIDLookup[T], options ...db.GetOptions) bool {
	randID, ok := uriString(ctx, param, kind, targetKey)
	if !ok {
		return false
	}
	value, ret := repo.GetByRandID(randID, options...)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return false
	}
	ctx.Set(key, value)
	return true
}

func SetFullAccess(ctx *gin.Context) {
	ctx.Set("FullAccess", true)
	ctx.Next()
}

func IsFullAccess(ctx *gin.Context) bool {
	return ctx.GetBool("FullAccess")
}

func SetRole(ctx *gin.Context) {
	if !loadByID(ctx, "Role", "roleID", "Role", "role_id", db.InitRoleRepo(db.DB.WithContext(ctx.Request.Context()))) {
		return
	}
	ctx.Next()
}

func GetRole(ctx *gin.Context) model.Role {
	return getResource[model.Role](ctx, "Role")
}

func SetGroup(ctx *gin.Context) {
	if !loadByID(ctx, "Group", "groupID", "Group", "group_id", db.InitGroupRepo(db.DB.WithContext(ctx.Request.Context()))) {
		return
	}
	ctx.Next()
}

func GetGroup(ctx *gin.Context) model.Group {
	return getResource[model.Group](ctx, "Group")
}

func SetPermission(ctx *gin.Context) {
	if !loadByID(ctx, "Permission", "permissionID", "Permission", "permission_id", db.InitPermissionRepo(db.DB.WithContext(ctx.Request.Context()))) {
		return
	}
	ctx.Next()
}

func GetPermission(ctx *gin.Context) model.Permission {
	return getResource[model.Permission](ctx, "Permission")
}

func SetUser(ctx *gin.Context) {
	if !loadByID(ctx, "User", "userID", "User", "user_id", db.InitUserRepo(db.DB.WithContext(ctx.Request.Context()))) {
		return
	}
	ctx.Next()
}

func GetUser(ctx *gin.Context) model.User {
	return getResource[model.User](ctx, "User")
}

func SetContest(ctx *gin.Context) {
	if !loadByID(ctx, "Contest", "contestID", "Contest", "contest_id", db.InitContestRepo(db.DB.WithContext(ctx.Request.Context()))) {
		return
	}
	contest := GetContest(ctx)
	if !IsFullAccess(ctx) && contest.Hidden {
		resp.AbortJSON(ctx, model.FailRetVal(i18n.Model.Contest.NotFound, "Contest", fmt.Sprintf("contest_id=%d", contest.ID), nil))
		return
	}
	ctx.Next()
}

func GetContest(ctx *gin.Context) model.Contest {
	return getResource[model.Contest](ctx, "Contest")
}

func SetTeam(ctx *gin.Context) {
	if !loadByID(ctx, "Team", "teamID", "Team", "team_id", db.InitTeamRepo(db.DB.WithContext(ctx.Request.Context()))) {
		return
	}
	team := GetTeam(ctx)
	if contest := GetContest(ctx); contest.ID != 0 && team.ContestID != contest.ID {
		resp.AbortJSON(ctx, model.FailRetVal(i18n.Model.Team.NotFound, "Team", fmt.Sprintf("team_id=%d", team.ID), nil))
		return
	}
	ctx.Next()
}

func SetTeamByUser(ctx *gin.Context) {
	team, ret := db.InitTeamRepo(db.DB.WithContext(ctx.Request.Context())).GetBy2ID(GetSelf(ctx).ID, GetContest(ctx).ID)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("Team", team)
	ctx.Next()
}

func GetTeam(ctx *gin.Context) model.Team {
	return getResource[model.Team](ctx, "Team")
}

func SetFile(t model.FileType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		randID, ok := uriString(ctx, "fileID", "File", "file_id")
		if !ok {
			return
		}
		options := db.GetOptions{}
		if t != "" {
			options = db.GetOptions{Conditions: map[string]any{"type": t}}
		}
		file, ret := db.InitFileRepo(db.DB.WithContext(ctx.Request.Context())).GetByUniqueField("rand_id", randID, options)
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
	file, ret := db.InitFileRepo(db.DB.WithContext(ctx.Request.Context())).Get(db.GetOptions{
		Conditions: map[string]any{"model": model.Name(challenge), "model_id": challenge.ID, "type": model.ChallengeFileType}},
	)
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("File", file)
	ctx.Next()
}

func SetTrafficFile(ctx *gin.Context) {
	file, ret := db.InitFileRepo(db.DB.WithContext(ctx.Request.Context())).Get(db.GetOptions{
		Conditions: map[string]any{"model": model.Name(GetVictim(ctx)), "model_id": GetVictim(ctx).ID, "type": model.TrafficFileType},
	})
	if !ret.OK {
		resp.AbortJSON(ctx, ret)
		return
	}
	ctx.Set("File", file)
	ctx.Next()
}

func SetTeamWriteupFile(ctx *gin.Context) {
	if !loadByRandID(ctx, "File", "fileID", "File", "file_id", db.InitFileRepo(db.DB.WithContext(ctx.Request.Context()))) {
		return
	}
	file := GetFile(ctx)
	if team := GetTeam(ctx); !(file.Model == model.Name(team) && file.ModelID == team.ID) {
		resp.AbortJSON(ctx, model.FailRetVal(i18n.Model.File.NotFound, "File", fmt.Sprintf("file_id=%d,team_id=%d", file.ID, team.ID), nil))
		return
	}
	ctx.Next()
}

// SetAttachmentFile 准备题目附件，比赛内用比赛题目副本的附件，测试模式用题库题目的附件
func SetAttachmentFile(test bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var source model.GeneratorSource
		if contestChallenge := GetContestChallenge(ctx); contestChallenge.ID > 0 {
			source = contestChallenge
		} else {
			challenge := GetChallenge(ctx)
			if test && challenge.Type == model.DynamicChallengeType {
				if ret := service.GenTestAttachment(db.DB.WithContext(ctx.Request.Context()), challenge); !ret.OK {
					resp.AbortJSON(ctx, ret)
					return
				}
			}
			source = challenge
		}
		path, pathRet := service.AttachmentPath(db.DB.WithContext(ctx.Request.Context()), source, GetTeam(ctx).ID)
		if !pathRet.OK {
			resp.AbortJSON(ctx, pathRet)
			return
		}
		record, ret := db.InitFileRepo(db.DB.WithContext(ctx.Request.Context())).Get(db.GetOptions{
			Conditions: map[string]any{"model": model.Name(source), "model_id": source.SourceID(), "type": model.ChallengeFileType}},
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

func GetFile(ctx *gin.Context) model.File {
	return getResource[model.File](ctx, "File")
}

func SetNotice(ctx *gin.Context) {
	if !loadByID(ctx, "Notice", "noticeID", "Notice", "notice_id", db.InitNoticeRepo(db.DB.WithContext(ctx.Request.Context()))) {
		return
	}
	notice := GetNotice(ctx)
	if contest := GetContest(ctx); contest.ID != 0 && notice.ContestID != contest.ID {
		resp.AbortJSON(ctx, model.FailRetVal(i18n.Model.Notice.NotFound, "Notice", fmt.Sprintf("notice_id=%d", notice.ID), nil))
		return
	}
	ctx.Next()
}

func GetNotice(ctx *gin.Context) model.Notice {
	return getResource[model.Notice](ctx, "Notice")
}

func SetChallenge(ctx *gin.Context) {
	if !loadByRandID(ctx, "Challenge", "challengeID", "Challenge", "challenge_id", db.InitChallengeRepo(db.DB.WithContext(ctx.Request.Context()))) {
		return
	}
	ctx.Next()
}

func GetChallenge(ctx *gin.Context) model.Challenge {
	return getResource[model.Challenge](ctx, "Challenge")
}

func SetContestChallenge(ctx *gin.Context) {
	if !loadByRandID(ctx, "ContestChallenge", "challengeID", "Challenge", "challenge_id", db.InitContestChallengeRepo(db.DB.WithContext(ctx.Request.Context())), db.GetOptions{
		Conditions: map[string]any{"contest_id": GetContest(ctx).ID},
	}) {
		return
	}
	contestChallenge := GetContestChallenge(ctx)
	if contestChallenge.Hidden && !IsFullAccess(ctx) {
		resp.AbortJSON(ctx, model.FailRetVal(i18n.Model.Contest.NotFound, "ContestChallenge", fmt.Sprintf("contest_challenge_id=%d", contestChallenge.ID), nil))
		return
	}
	ctx.Next()
}

func GetContestChallenge(ctx *gin.Context) model.ContestChallenge {
	return getResource[model.ContestChallenge](ctx, "ContestChallenge")
}

func SetContestFlag(ctx *gin.Context) {
	if !loadByID(ctx, "ContestFlag", "flagID", "ContestFlag", "contest_flag_id", db.InitContestFlagRepo(db.DB.WithContext(ctx.Request.Context()))) {
		return
	}
	contestFlag := GetContestFlag(ctx)
	if contest := GetContest(ctx); contest.ID != 0 && contestFlag.ContestID != contest.ID {
		resp.AbortJSON(ctx, model.FailRetVal(i18n.Model.ContestFlag.NotFound, "ContestFlag", fmt.Sprintf("contest_flag_id=%d", contestFlag.ID), nil))
		return
	}
	if contestChallenge := GetContestChallenge(ctx); contestChallenge.ID != 0 && contestFlag.ContestChallengeID != contestChallenge.ID {
		resp.AbortJSON(ctx, model.FailRetVal(i18n.Model.ContestFlag.NotFound, "ContestFlag", fmt.Sprintf("contest_flag_id=%d", contestFlag.ID), nil))
		return
	}
	ctx.Next()
}

func GetContestFlag(ctx *gin.Context) model.ContestFlag {
	return getResource[model.ContestFlag](ctx, "ContestFlag")
}

func SetGenerator(ctx *gin.Context) {
	if !loadByID(ctx, "Generator", "generatorID", "Generator", "generator_id", db.InitGeneratorRepo(db.DB.WithContext(ctx.Request.Context()))) {
		return
	}
	generator := GetGenerator(ctx)
	if contest := GetContest(ctx); contest.ID != 0 && (!generator.ContestID.Valid || generator.ContestID.V != contest.ID) {
		resp.AbortJSON(ctx, model.FailRetVal(i18n.Response.BadRequest, "Generator", fmt.Sprintf("generator_id=%d", generator.ID), nil))
		return
	}
	ctx.Next()
}

func GetGenerator(ctx *gin.Context) model.Generator {
	return getResource[model.Generator](ctx, "Generator")
}

func SetVictim(ctx *gin.Context) {
	if !loadByID(ctx, "Victim", "victimID", "Victim", "victim_id", db.InitVictimRepo(db.DB.WithContext(ctx.Request.Context())), db.GetOptions{Deleted: true}) {
		return
	}
	victim := GetVictim(ctx)
	if contest := GetContest(ctx); contest.ID != 0 && (!victim.ContestID.Valid || victim.ContestID.V != contest.ID) {
		resp.AbortJSON(ctx, model.FailRetVal(i18n.Model.Victim.NotFound, "Victim", fmt.Sprintf("victim_id=%d", victim.ID), nil))
		return
	}
	if team := GetTeam(ctx); team.ID != 0 && (!victim.TeamID.Valid || victim.TeamID.V != team.ID) {
		resp.AbortJSON(ctx, model.FailRetVal(i18n.Model.Victim.NotFound, "Victim", fmt.Sprintf("victim_id=%d", victim.ID), nil))
		return
	}
	ctx.Next()
}

func GetVictim(ctx *gin.Context) model.Victim {
	return getResource[model.Victim](ctx, "Victim")
}

func SetCheat(ctx *gin.Context) {
	if !loadByID(ctx, "Cheat", "cheatID", "Cheat", "cheat_id", db.InitCheatRepo(db.DB.WithContext(ctx.Request.Context()))) {
		return
	}
	cheat := GetCheat(ctx)
	if contest := GetContest(ctx); contest.ID != 0 && cheat.ContestID != contest.ID {
		resp.AbortJSON(ctx, model.FailRetVal(i18n.Model.Cheat.NotFound, "Cheat", fmt.Sprintf("cheat_id=%d", cheat.ID), nil))
		return
	}
	ctx.Next()
}

func GetCheat(ctx *gin.Context) model.Cheat {
	return getResource[model.Cheat](ctx, "Cheat")
}

func SetOauth(ctx *gin.Context) {
	if !loadByID(ctx, "Oauth", "oauthID", "OAuthProvider", "oauth_id", db.InitOauthRepo(db.DB.WithContext(ctx.Request.Context()))) {
		return
	}
	ctx.Next()
}

func GetOauth(ctx *gin.Context) model.Oauth {
	return getResource[model.Oauth](ctx, "Oauth")
}

// SetOauthUri 只传递名称，不查库，后续由内存中的 map 解析
func SetOauthUri(ctx *gin.Context) {
	name, ok := uriString(ctx, "oauth", "OAuthProvider", "oauth")
	if !ok {
		return
	}
	ctx.Set("OauthUri", name)
	ctx.Next()
}

func GetOauthUri(ctx *gin.Context) string {
	return getResource[string](ctx, "OauthUri")
}

func SetSmtp(ctx *gin.Context) {
	if !loadByID(ctx, "Smtp", "smtpID", "Smtp", "smtp_id", db.InitSmtpRepo(db.DB.WithContext(ctx.Request.Context()))) {
		return
	}
	ctx.Next()
}

func GetSmtp(ctx *gin.Context) model.Smtp {
	return getResource[model.Smtp](ctx, "Smtp")
}

func SetCronJob(ctx *gin.Context) {
	if !loadByID(ctx, "CronJob", "cronJobID", "CronJob", "cron_job_id", db.InitCronJobRepo(db.DB.WithContext(ctx.Request.Context()))) {
		return
	}
	ctx.Next()
}

func GetCronJob(ctx *gin.Context) model.CronJob {
	return getResource[model.CronJob](ctx, "CronJob")
}

func SetWebhook(ctx *gin.Context) {
	if !loadByID(ctx, "Webhook", "webhookID", "Webhook", "webhook_id", db.InitWebhookRepo(db.DB.WithContext(ctx.Request.Context()))) {
		return
	}
	ctx.Next()
}

func GetWebhook(ctx *gin.Context) model.Webhook {
	return getResource[model.Webhook](ctx, "Webhook")
}
