package middleware

import (
    "CBCTF/internal/model"
    db "CBCTF/internal/repo"
    "fmt"

    "github.com/gin-gonic/gin"
)

func Events(ctx *gin.Context) {
    method, path, ip := ctx.Request.Method, ctx.FullPath(), ctx.ClientIP()

	ctx.Next()

	userID := GetSelfID(ctx)
	teamID := GetTeam(ctx).ID
	contestID := GetContest(ctx).ID
	contestChallengeID := GetContestChallenge(ctx).ID
	options := db.CreateEventOptions{
        IP:       ip,
        Magic:    GetMagic(ctx),
        Status:   ctx.Writer.Status(),
        Success:  ctx.Writer.Status() >= 200 && ctx.Writer.Status() < 400,
        UserAgent: func() string {
            ua := ctx.Request.UserAgent()
            if len(ua) > 255 {
                return ua[:255]
            }
            return ua
        }(),
        TraceID: GetTraceID(ctx),
	}
	if userID == 0 {
		options.UserID = nil
	} else {
		options.UserID = &userID
	}
	if teamID == 0 {
		options.TeamID = nil
	} else {
		options.TeamID = &teamID
	}
	if contestID == 0 {
		options.ContestID = nil
	} else {
		options.ContestID = &contestID
	}
	if contestChallengeID == 0 {
		options.ContestChallengeID = nil
	} else {
		options.ContestChallengeID = &contestChallengeID
	}

	switch method {
	case "GET":
		switch path {
		case "/contests/:contestID/challenges/:challengeID/attachment":
			options.Type = model.DownloadAttachmentEventType
			options.Desc = fmt.Sprintf("User %d download attachment for usage %d in contest %d as team %d", userID, contestChallengeID, contestID, teamID)
		default:
			return
		}
	case "POST":
		switch path {
		case "/register":
			if userID == 0 {
				return
			}
			options.Type = model.UserRegisterEventType
			options.Desc = fmt.Sprintf("User %d register", userID)
		case "/login":
			if userID == 0 {
				return
			}
			options.Type = model.UserLoginEventType
			options.Desc = fmt.Sprintf("User %d login", userID)
		case "/me/activate":
			if userID == 0 {
				return
			}
			options.Type = model.UserVerifyEmailEventType
			options.Desc = fmt.Sprintf("User %d verify email", userID)
		case "/contests/:contestID/teams/join":
			if teamID == 0 {
				return
			}
			options.Type = model.JoinTeamEventType
			options.Desc = fmt.Sprintf("User %d join team %d", userID, teamID)
		case "/contests/:contestID/teams/create":
			if teamID == 0 {
				return
			}
			options.Type = model.CreateTeamEventType
			options.Desc = fmt.Sprintf("User %d create team %d", userID, teamID)
		case "/contests/:contestID/teams/me/leave":
			options.Type = model.LeaveTeamEventType
			options.Desc = fmt.Sprintf("User %d leave team %d", userID, teamID)
		case "/contests/:contestID/teams/me/kick":
			options.Type = model.KickMemberEventType
			options.Desc = fmt.Sprintf("User %d kick member from team %d", userID, teamID)
		case "/contests/:contestID/challenges/:challengeID/init":
			options.Type = model.InitChallengeEventType
			options.Desc = fmt.Sprintf("User %d init usage %d in contest %d as team %d", userID, contestChallengeID, contestID, teamID)
		case "/contests/:contestID/challenges/:challengeID/reset":
			options.Type = model.ResetChallengeEventType
			options.Desc = fmt.Sprintf("User %d reset usage %d in contest %d as team %d", userID, contestChallengeID, contestID, teamID)
		case "/contests/:contestID/challenges/:challengeID/submit":
			options.Type = model.SubmitFlagEventType
			options.Desc = fmt.Sprintf("User %d submit flag for usage %d in contest %d as team %d", userID, contestChallengeID, contestID, teamID)
		case "/contests/:contestID/challenges/:challengeID/start":
			options.Type = model.StartVictimEventType
			options.Desc = fmt.Sprintf("User %d start victim for usage %d in contest %d as team %d", userID, contestChallengeID, contestID, teamID)
		case "/contests/:contestID/challenges/:challengeID/increase":
			options.Type = model.IncreaseVictimEventType
			options.Desc = fmt.Sprintf("User %d increase victim for usage %d in contest %d as team %d", userID, contestChallengeID, contestID, teamID)
		case "/contests/:contestID/challenges/:challengeID/stop":
			options.Type = model.StopVictimEventType
			options.Desc = fmt.Sprintf("User %d stop victim for usage %d in contest %d as team %d", userID, contestChallengeID, contestID, teamID)
		case "/contests/:contestID/writeups":
			options.Type = model.UploadWriteUpEventType
			options.Desc = fmt.Sprintf("User %d upload writeup for contest %d as team %d", userID, contestID, teamID)
		default:
			return
		}
	case "PUT":
		switch path {
		case "/me":
			if userID == 0 {
				return
			}
			options.Type = model.UserUpdateEventType
			options.Desc = fmt.Sprintf("User %d update self", userID)
		case "/me/password":
			if userID == 0 {
				return
			}
			options.Type = model.UserUpdatePasswordEventType
			options.Desc = fmt.Sprintf("User %d update password", userID)
		case "/contests/:contestID/teams/me":
			options.Type = model.UpdateTeamEventType
			options.Desc = fmt.Sprintf("User %d update team %d", userID, teamID)
		default:
			return
		}
	case "DELETE":
		switch path {
		case "/me":
			if userID == 0 {
				return
			}
			options.Type = model.UserDeleteEventType
			options.Desc = fmt.Sprintf("User %d delete self", userID)
		default:
			return
		}
	default:
		return
	}
	db.InitEventRepo(db.DB.WithContext(ctx)).Create(options)
}
