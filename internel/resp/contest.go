package resp

import (
	"CBCTF/internel/config"
	"CBCTF/internel/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
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
		"prefix":    contest.Prefix,
		"notices":   len(contest.Notices),
		"avatar":    fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(contest.Avatar, "/")),
		"solved": func() int {
			count := 0
			for _, submission := range contest.Submissions {
				if submission.Solved {
					count++
				}
			}
			return count
		}(),
	}
	if admin {
		data["hidden"] = contest.Hidden
		data["captcha"] = contest.Captcha
	}
	return data
}
