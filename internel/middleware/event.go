package middleware

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Events(ctx *gin.Context) {
	method, path, ip := ctx.Request.Method, ctx.FullPath(), ctx.ClientIP()

	ctx.Next()

	references := model.References{
		UserID:    GetSelfID(ctx),
		TeamID:    GetTeam(ctx).ID,
		ContestID: GetContest(ctx).ID,
		UsageID:   GetUsage(ctx).ID,
	}
	options := db.CreateEventOptions{
		References: references,
		IP:         ip,
		Magic:      GetMagic(ctx),
	}

	switch method {
	case "GET":
		switch path {
		case "/contests/:contestID/challenges/:challengeID/attachment":
			options.Type = model.DownloadAttachmentEventType
			options.Desc = fmt.Sprintf("User %d download attachment for usage %d in contest %d as team %d", references.UserID, references.UsageID, references.ContestID, references.TeamID)
		default:
			return
		}
	case "POST":
		switch path {
		case "/register":
			if references.UserID == 0 {
				return
			}
			options.Type = model.UserRegisterEventType
			options.Desc = fmt.Sprintf("User %d register", references.UserID)
		case "/login":
			if references.UserID == 0 {
				return
			}
			options.Type = model.UserLoginEventType
			options.Desc = fmt.Sprintf("User %d login", references.UserID)
		case "/me/activate":
			if references.UserID == 0 {
				return
			}
			options.Type = model.UserVerifyEmailEventType
			options.Desc = fmt.Sprintf("User %d verify email", references.UserID)
		case "/contests/:contestID/teams/join":
			if references.TeamID == 0 {
				return
			}
			options.Type = model.JoinTeamEventType
			options.Desc = fmt.Sprintf("User %d join team %d", references.UserID, references.TeamID)
		case "/contests/:contestID/teams/create":
			if references.TeamID == 0 {
				return
			}
			options.Type = model.CreateTeamEventType
			options.Desc = fmt.Sprintf("User %d create team %d", references.UserID, references.TeamID)
		case "/contests/:contestID/teams/me/leave":
			options.Type = model.LeaveTeamEventType
			options.Desc = fmt.Sprintf("User %d leave team %d", references.UserID, references.TeamID)
		case "/contests/:contestID/teams/me/kick":
			options.Type = model.KickMemberEventType
			options.Desc = fmt.Sprintf("User %d kick member from team %d", references.UserID, references.TeamID)
		case "/contests/:contestID/challenges/:challengeID/init":
			options.Type = model.InitChallengeEventType
			options.Desc = fmt.Sprintf("User %d init usage %d in contest %d as team %d", references.UserID, references.UsageID, references.ContestID, references.TeamID)
		case "/contests/:contestID/challenges/:challengeID/reset":
			options.Type = model.ResetChallengeEventType
			options.Desc = fmt.Sprintf("User %d reset usage %d in contest %d as team %d", references.UserID, references.UsageID, references.ContestID, references.TeamID)
		case "/contests/:contestID/challenges/:challengeID/submit":
			options.Type = model.SubmitFlagEventType
			options.Desc = fmt.Sprintf("User %d submit flag for usage %d in contest %d as team %d", references.UserID, references.UsageID, references.ContestID, references.TeamID)
		case "/contests/:contestID/challenges/:challengeID/start":
			options.Type = model.StartVictimEventType
			options.Desc = fmt.Sprintf("User %d start victim for usage %d in contest %d as team %d", references.UserID, references.UsageID, references.ContestID, references.TeamID)
		case "/contests/:contestID/challenges/:challengeID/increase":
			options.Type = model.IncreaseVictimEventType
			options.Desc = fmt.Sprintf("User %d increase victim for usage %d in contest %d as team %d", references.UserID, references.UsageID, references.ContestID, references.TeamID)
		case "/contests/:contestID/challenges/:challengeID/stop":
			options.Type = model.StopVictimEventType
			options.Desc = fmt.Sprintf("User %d stop victim for usage %d in contest %d as team %d", references.UserID, references.UsageID, references.ContestID, references.TeamID)
		case "/contests/:contestID/writeups":
			options.Type = model.UploadWriteUpEventType
			options.Desc = fmt.Sprintf("User %d upload writeup for contest %d as team %d", references.UserID, references.ContestID, references.TeamID)
		default:
			return
		}
	case "PUT":
		switch path {
		case "/me":
			if references.UserID == 0 {
				return
			}
			options.Type = model.UserUpdateEventType
			options.Desc = fmt.Sprintf("User %d update self", references.UserID)
		case "/me/password":
			if references.UserID == 0 {
				return
			}
			options.Type = model.UserUpdatePasswordEventType
			options.Desc = fmt.Sprintf("User %d update password", references.UserID)
		case "/contests/:contestID/teams/me":
			options.Type = model.UpdateTeamEventType
			options.Desc = fmt.Sprintf("User %d update team %d", references.UserID, references.TeamID)
		default:
			return
		}
	case "DELETE":
		switch path {
		case "/me":
			if references.UserID == 0 {
				return
			}
			options.Type = model.UserDeleteEventType
			options.Desc = fmt.Sprintf("User %d delete self", references.UserID)
		default:
			return
		}
	default:
		return
	}
	if _, ok, msg := db.InitEventRepo(db.DB.WithContext(ctx)).Create(options); !ok {
		log.Logger.Warningf("Failed to record event: %v beacause of %s", options, msg)
	}
}
