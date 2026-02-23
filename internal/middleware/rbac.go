package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckPermission(ctx *gin.Context) {
	var permission string
	switch fmt.Sprintf("%s %s", ctx.Request.Method, ctx.FullPath()) {
	// /me
	case "GET /me":
		permission = model.PermSelfRead
	case "PUT /me":
		permission = model.PermSelfUpdate
	case "PUT /me/password":
		permission = model.PermSelfUpdate
	case "DELETE /me":
		permission = model.PermSelfDelete
	case "POST /me/picture":
		permission = model.PermSelfUpdate
	case "POST /me/activate":
		permission = model.PermSelfActivate

	// /contests/:contestID
	case "GET /contests/:contestID":
		permission = model.PermUserContestRead
	case "GET /contests/:contestID/rank":
		permission = model.PermUserContestRank
	case "GET /contests/:contestID/scoreboard":
		permission = model.PermUserContestRank
	case "GET /contests/:contestID/timeline":
		permission = model.PermUserContestRank
	case "POST /contests/:contestID/teams/join":
		permission = model.PermUserTeamJoin
	case "POST /contests/:contestID/teams/create":
		permission = model.PermUserTeamCreate
	case "GET /contests/:contestID/teams/me":
		permission = model.PermUserTeamRead
	case "GET /contests/:contestID/teams/me/captcha":
		permission = model.PermUserTeamRead
	case "GET /contests/:contestID/teams/me/users":
		permission = model.PermUserTeamRead
	case "PUT /contests/:contestID/teams/me/captcha":
		permission = model.PermUserTeamUpdate
	case "PUT /contests/:contestID/teams/me":
		permission = model.PermUserTeamUpdate
	case "POST /contests/:contestID/teams/me/picture":
		permission = model.PermUserTeamUpdate
	case "DELETE /contests/:contestID/teams/me":
		permission = model.PermUserTeamDelete
	case "POST /contests/:contestID/teams/me/kick":
		permission = model.PermUserTeamUpdate
	case "POST /contests/:contestID/teams/me/leave":
		permission = model.PermUserTeamRead
	case "GET /contests/:contestID/notices":
		permission = model.PermUserNoticeList
	case "GET /contests/:contestID/challenges":
		permission = model.PermUserChallengeList
	case "GET /contests/:contestID/challenges/categories":
		permission = model.PermUserChallengeList
	case "GET /contests/:contestID/challenges/:challengeID":
		permission = model.PermUserChallengeRead
	case "GET /contests/:contestID/challenges/:challengeID/attachment":
		permission = model.PermUserChallengeRead
	case "POST /contests/:contestID/challenges/:challengeID/init":
		permission = model.PermUserChallengeInit
	case "POST /contests/:contestID/challenges/:challengeID/reset":
		permission = model.PermUserChallengeReset
	case "POST /contests/:contestID/challenges/:challengeID/start":
		permission = model.PermUserVictimControl
	case "POST /contests/:contestID/challenges/:challengeID/increase":
		permission = model.PermUserVictimControl
	case "POST /contests/:contestID/challenges/:challengeID/stop":
		permission = model.PermUserVictimControl
	case "POST /contests/:contestID/challenges/:challengeID/submit":
		permission = model.PermUserChallengeSubmit
	case "POST /contests/:contestID/writeups":
		permission = model.PermUserWriteupUpload
	case "GET /contests/:contestID/writeups":
		permission = model.PermUserWriteupList

	// /admin 基础
	case "GET /admin/ip":
		permission = model.PermAdminIPSearch
	case "GET /admin/models":
		permission = model.PermAdminModelsSearch
	case "GET /admin/search":
		permission = model.PermAdminModelsSearch

	// /admin/system
	case "GET /admin/system/status":
		permission = model.PermAdminSystemRead
	case "GET /admin/system/config":
		permission = model.PermAdminSystemRead
	case "PUT /admin/system/config":
		permission = model.PermAdminSystemUpdate
	case "POST /admin/system/restart":
		permission = model.PermAdminSystemRestart

	// /admin/permissions
	case "GET /admin/permissions":
		permission = model.PermAdminPermissionList
	case "PUT /admin/permissions/:permissionID":
		permission = model.PermAdminPermissionUpdate

	// /admin/roles
	case "GET /admin/roles":
		permission = model.PermAdminRoleList
	case "POST /admin/roles":
		permission = model.PermAdminRoleCreate
	case "GET /admin/roles/:roleID":
		permission = model.PermAdminRoleRead
	case "GET /admin/roles/:roleID/permissions":
		permission = model.PermAdminPermissionList
	case "PUT /admin/roles/:roleID":
		permission = model.PermAdminRoleUpdate
	case "DELETE /admin/roles/:roleID":
		permission = model.PermAdminRoleDelete
	case "POST /admin/roles/:roleID/permissions":
		permission = model.PermAdminRoleAssign
	case "DELETE /admin/roles/:roleID/permissions":
		permission = model.PermAdminRoleRevoke

	// /admin/groups
	case "GET /admin/groups":
		permission = model.PermAdminGroupList
	case "POST /admin/groups":
		permission = model.PermAdminGroupCreate
	case "GET /admin/groups/:groupID":
		permission = model.PermAdminGroupRead
	case "GET /admin/groups/:groupID/users":
		permission = model.PermAdminUserList
	case "PUT /admin/groups/:groupID":
		permission = model.PermAdminGroupUpdate
	case "DELETE /admin/groups/:groupID":
		permission = model.PermAdminGroupDelete
	case "POST /admin/groups/:groupID/users":
		permission = model.PermAdminUserAssign
	case "DELETE /admin/groups/:groupID/users":
		permission = model.PermAdminUserRevoke

	// /admin/users
	case "GET /admin/users":
		permission = model.PermAdminUserList
	case "POST /admin/users":
		permission = model.PermAdminUserCreate
	case "GET /admin/users/:userID":
		permission = model.PermAdminUserRead
	case "PUT /admin/users/:userID":
		permission = model.PermAdminUserUpdate
	case "DELETE /admin/users/:userID":
		permission = model.PermAdminUserDelete
	case "POST /admin/users/:userID/picture":
		permission = model.PermAdminUserUpdate

	// /admin/oauth
	case "GET /admin/oauth":
		permission = model.PermAdminOauthList
	case "POST /admin/oauth":
		permission = model.PermAdminOauthCreate
	case "PUT /admin/oauth/:oauthID":
		permission = model.PermAdminOauthUpdate
	case "POST /admin/oauth/:oauthID/picture":
		permission = model.PermAdminOauthUpdate
	case "DELETE /admin/oauth/:oauthID":
		permission = model.PermAdminOauthDelete

	// /admin/email + /admin/smtp
	case "GET /admin/email":
		permission = model.PermAdminSMTPList
	case "GET /admin/smtp":
		permission = model.PermAdminSMTPList
	case "POST /admin/smtp":
		permission = model.PermAdminSMTPCreate
	case "PUT /admin/smtp/:smtpID":
		permission = model.PermAdminSMTPUpdate
	case "DELETE /admin/smtp/:smtpID":
		permission = model.PermAdminSMTPDelete
	case "GET /admin/smtp/:smtpID/email":
		permission = model.PermAdminSMTPList

	// /admin/webhook
	case "GET /admin/webhook":
		permission = model.PermAdminWebhookList
	case "GET /admin/webhook/events":
		permission = model.PermAdminWebhookList
	case "GET /admin/webhook/history":
		permission = model.PermAdminWebhookList
	case "POST /admin/webhook":
		permission = model.PermAdminWebhookCreate
	case "PUT /admin/webhook/:webhookID":
		permission = model.PermAdminWebhookUpdate
	case "DELETE /admin/webhook/:webhookID":
		permission = model.PermAdminWebhookDelete
	case "GET /admin/webhook/:webhookID/history":
		permission = model.PermAdminWebhookList

	// /admin/challenges
	case "GET /admin/challenges":
		permission = model.PermAdminChallengeList
	case "GET /admin/challenges/categories":
		permission = model.PermAdminChallengeList
	case "POST /admin/challenges":
		permission = model.PermAdminChallengeCreate
	case "GET /admin/challenges/:challengeID/download":
		permission = model.PermAdminChallengeRead
	case "PUT /admin/challenges/:challengeID":
		permission = model.PermAdminChallengeUpdate
	case "DELETE /admin/challenges/:challengeID":
		permission = model.PermAdminChallengeDelete
	case "POST /admin/challenges/:challengeID/upload":
		permission = model.PermAdminChallengeUpdate
	case "GET /admin/challenges/:challengeID/test":
		permission = model.PermAdminChallengeTest
	case "GET /admin/challenges/:challengeID/test/attachment":
		permission = model.PermAdminChallengeTest
	case "POST /admin/challenges/:challengeID/test/start":
		permission = model.PermAdminChallengeTest
	case "POST /admin/challenges/:challengeID/test/stop":
		permission = model.PermAdminChallengeTest

	// /admin/contests
	case "GET /admin/contests":
		permission = model.PermAdminContestList
	case "POST /admin/contests":
		permission = model.PermAdminContestCreate
	case "GET /admin/contests/:contestID":
		permission = model.PermAdminContestRead
	case "PUT /admin/contests/:contestID":
		permission = model.PermAdminContestUpdate
	case "DELETE /admin/contests/:contestID":
		permission = model.PermAdminContestDelete
	case "POST /admin/contests/:contestID/picture":
		permission = model.PermAdminContestUpdate
	case "GET /admin/contests/:contestID/rank":
		permission = model.PermAdminContestRank
	case "GET /admin/contests/:contestID/scoreboard":
		permission = model.PermAdminContestRank
	case "GET /admin/contests/:contestID/timeline":
		permission = model.PermAdminContestRank

	// /admin/contests/:contestID/teams
	case "GET /admin/contests/:contestID/teams":
		permission = model.PermAdminTeamList
	case "GET /admin/contests/:contestID/teams/:teamID":
		permission = model.PermAdminTeamRead
	case "GET /admin/contests/:contestID/teams/:teamID/users":
		permission = model.PermAdminTeamRead
	case "PUT /admin/contests/:contestID/teams/:teamID":
		permission = model.PermAdminTeamUpdate
	case "DELETE /admin/contests/:contestID/teams/:teamID":
		permission = model.PermAdminTeamDelete
	case "POST /admin/contests/:contestID/teams/:teamID/kick":
		permission = model.PermAdminTeamUpdate
	case "POST /admin/contests/:contestID/teams/:teamID/picture":
		permission = model.PermAdminTeamUpdate
	case "GET /admin/contests/:contestID/teams/:teamID/flags":
		permission = model.PermAdminTeamRead
	case "GET /admin/contests/:contestID/teams/:teamID/submissions":
		permission = model.PermAdminTeamRead
	case "GET /admin/contests/:contestID/teams/:teamID/victims":
		permission = model.PermAdminTeamRead
	case "GET /admin/contests/:contestID/teams/:teamID/victims/:victimID/traffic":
		permission = model.PermAdminTeamRead
	case "GET /admin/contests/:contestID/teams/:teamID/victims/:victimID/traffic/download":
		permission = model.PermAdminTeamRead
	case "GET /admin/contests/:contestID/teams/:teamID/writeups":
		permission = model.PermAdminTeamWriteupList
	case "GET /admin/contests/:contestID/teams/:teamID/writeups/:fileID":
		permission = model.PermAdminTeamWriteupRead

	// /admin/contests/:contestID/notices
	case "GET /admin/contests/:contestID/notices":
		permission = model.PermAdminNoticeList
	case "POST /admin/contests/:contestID/notices":
		permission = model.PermAdminNoticeCreate
	case "PUT /admin/contests/:contestID/notices/:noticeID":
		permission = model.PermAdminNoticeUpdate
	case "DELETE /admin/contests/:contestID/notices/:noticeID":
		permission = model.PermAdminNoticeDelete

	// /admin/contests/:contestID/cheats
	case "GET /admin/contests/:contestID/cheats":
		permission = model.PermAdminCheatList
	case "DELETE /admin/contests/:contestID/cheats":
		permission = model.PermAdminCheatDelete
	case "POST /admin/contests/:contestID/cheats":
		permission = model.PermAdminCheatCreate
	case "PUT /admin/contests/:contestID/cheats/:cheatID":
		permission = model.PermAdminCheatUpdate
	case "DELETE /admin/contests/:contestID/cheats/:cheatID":
		permission = model.PermAdminCheatDelete

	// /admin/contests/:contestID/challenges
	case "GET /admin/contests/:contestID/challenges":
		permission = model.PermAdminContestChallengeList
	case "GET /admin/contests/:contestID/challenges/others":
		permission = model.PermAdminChallengeList
	case "GET /admin/contests/:contestID/challenges/categories":
		permission = model.PermAdminContestChallengeList
	case "POST /admin/contests/:contestID/challenges":
		permission = model.PermAdminContestChallengeCreate
	case "PUT /admin/contests/:contestID/challenges/:challengeID":
		permission = model.PermAdminContestChallengeUpdate
	case "DELETE /admin/contests/:contestID/challenges/:challengeID":
		permission = model.PermAdminContestChallengeDelete
	case "GET /admin/contests/:contestID/challenges/:challengeID/flags":
		permission = model.PermAdminContestChallengeFlagList
	case "PUT /admin/contests/:contestID/challenges/:challengeID/flags/:flagID":
		permission = model.PermAdminContestChallengeFlagUpdate

	// /admin/contests/:contestID/images
	case "GET /admin/contests/:contestID/images":
		permission = model.PermAdminImagePull
	case "POST /admin/contests/:contestID/images":
		permission = model.PermAdminImagePull

	// /admin/contests/:contestID/victims
	case "GET /admin/contests/:contestID/victims":
		permission = model.PermAdminVictimControl
	case "POST /admin/contests/:contestID/victims":
		permission = model.PermAdminVictimControl
	case "DELETE /admin/contests/:contestID/victims":
		permission = model.PermAdminVictimControl

	// /admin/files
	case "GET /admin/files":
		permission = model.PermAdminFileList
	case "DELETE /admin/files":
		permission = model.PermAdminFileDelete
	case "GET /admin/files/:fileID":
		permission = model.PermAdminFileRead

	// /admin/logs
	case "GET /admin/logs":
		permission = model.PermAdminLogRead

	default:
		ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Request.Forbidden})
		return
	}

	pass, ret := db.InitPermissionRepo(db.DB).CheckUserPermission(GetSelf(ctx).ID, permission)
	if !ret.OK {
		ctx.AbortWithStatusJSON(http.StatusOK, ret)
		return
	}
	if !pass {
		ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Request.Forbidden})
		return
	}
	ctx.Next()
}
