package resp

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

func GetContestResp(contest model.Contest, admin bool) gin.H {
	data := gin.H{
		"id":        contest.ID,
		"name":      contest.Name,
		"desc":      contest.Desc,
		"start":     contest.Start,
		"duration":  contest.Duration.Seconds(),
		"rules":     contest.Rules,
		"prizes":    contest.Prizes,
		"size":      contest.Size,
		"timelines": contest.Timelines,
		"teams":     len(contest.Teams),
		"users":     len(contest.Users),
		"notices":   len(contest.Notices),
		"prefix":    contest.Prefix,
		"avatar":    contest.Avatar,
		"hidden":    contest.Hidden,
		"blood":     contest.Blood,
	}
	if admin {
		data["captcha"] = contest.Captcha
	}
	return data
}
