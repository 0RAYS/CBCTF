package middleware

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Events(ctx *gin.Context) {
	ctx.Next()

	options := db.CreateEventOptions{
		UserID:    GetSelfID(ctx),
		TeamID:    GetTeam(ctx).ID,
		ContestID: GetContest(ctx).ID,
		UsageID:   GetUsage(ctx).ID,
		IP:        ctx.ClientIP(),
		Magic:     GetMagic(ctx),
	}

	switch ctx.Request.Method {
	case "GET":
		switch ctx.FullPath() {
		case "/contests/:contestID/challenges/:challengeID/attachment":
			options.Type = model.DownloadAttachmentEventType
			options.Desc = fmt.Sprintf("User %d download attachment for challenge %d", options.UserID, options.UsageID)
		default:
			return
		}
	case "POST":
		switch ctx.FullPath() {
		case "/register":
			if options.UserID == 0 {
				return
			}
			options.Type = model.UserRegisterEventType
			options.Desc = fmt.Sprintf("User %d register", options.UserID)
		case "/login":
			if options.UserID == 0 {
				return
			}
			options.Type = model.UserLoginEventType
			options.Desc = fmt.Sprintf("User %d login", options.UserID)
		case "/me/activate":
			if options.UserID == 0 {
				return
			}
			options.Type = model.UserVerifyEmailEventType
			options.Desc = fmt.Sprintf("User %d verify email", options.UserID)
		case "/contests/:contestID/teams/join":
			if options.TeamID == 0 {
				return
			}
			options.Type = model.JoinTeamEventType
			options.Desc = fmt.Sprintf("User %d join team %d", options.UserID, options.TeamID)
		case "/contests/:contestID/teams/create":
			if options.TeamID == 0 {
				return
			}
			options.Type = model.CreateTeamEventType
			options.Desc = fmt.Sprintf("User %d create team %d", options.UserID, options.TeamID)
		case "/contests/:contestID/teams/me/leave":
			options.Type = model.LeaveTeamEventType
			options.Desc = fmt.Sprintf("User %d leave team %d", options.UserID, options.TeamID)
		case "/contests/:contestID/teams/me/kick":
			options.Type = model.KickMemberEventType
			options.Desc = fmt.Sprintf("User %d kick member from team %d", options.UserID, options.TeamID)
		case "/contests/:contestID/challenges/:challengeID/init":
			options.Type = model.InitChallengeEventType
			options.Desc = fmt.Sprintf("User %d init challenge %d", options.UserID, options.UsageID)
		case "/contests/:contestID/challenges/:challengeID/reset":
			options.Type = model.ResetChallengeEventType
			options.Desc = fmt.Sprintf("User %d reset challenge %d", options.UserID, options.UsageID)
		case "/contests/:contestID/challenges/:challengeID/submit":
			options.Type = model.SubmitFlagEventType
			options.Desc = fmt.Sprintf("User %d submit flag for challenge %d", options.UserID, options.UsageID)
		case "/contests/:contestID/challenges/:challengeID/start":
			options.Type = model.StartVictimEventType
			options.Desc = fmt.Sprintf("User %d start victim for challenge %d", options.UserID, options.UsageID)
		case "/contests/:contestID/challenges/:challengeID/increase":
			options.Type = model.IncreaseVictimEventType
			options.Desc = fmt.Sprintf("User %d increase victim for challenge %d", options.UserID, options.UsageID)
		case "/contests/:contestID/challenges/:challengeID/stop":
			options.Type = model.StopVictimEventType
			options.Desc = fmt.Sprintf("User %d stop victim for challenge %d", options.UserID, options.UsageID)
		case "/contests/:contestID/writeups":
			options.Type = model.UploadWriteUpEventType
			options.Desc = fmt.Sprintf("User %d upload writeup for contest %d", options.UserID, options.ContestID)
		default:
			return
		}
	case "PUT":
		switch ctx.FullPath() {
		case "/me":
			if options.UserID == 0 {
				return
			}
			options.Type = model.UserUpdateEventType
			options.Desc = fmt.Sprintf("User %d update self", options.UserID)
		case "/me/password":
			if options.UserID == 0 {
				return
			}
			options.Type = model.UserUpdatePasswordEventType
			options.Desc = fmt.Sprintf("User %d update password", options.UserID)
		case "/contests/:contestID/teams/me":
			options.Type = model.UpdateTeamEventType
			options.Desc = fmt.Sprintf("User %d update team %d", options.UserID, options.TeamID)
		default:
			return
		}
	case "DELETE":
		switch ctx.FullPath() {
		case "/me":
			if options.UserID == 0 {
				return
			}
			options.Type = model.UserDeleteEventType
			options.Desc = fmt.Sprintf("User %d delete self", options.UserID)
		default:
			return
		}
	default:
		return
	}
	event, ok, msg := db.InitEventRepo(db.DB.WithContext(ctx)).Create(options)
	if !ok {
		log.Logger.Warningf("Failed to record event: %v beacause of %s", event, msg)
	}
}
