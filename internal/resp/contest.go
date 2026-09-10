package resp

import (
	"CBCTF/internal/model"
	"CBCTF/internal/view"

	"github.com/gin-gonic/gin"
)

func GetContestResp(contestView view.ContestView, admin bool) gin.H {
	contest := contestView.Contest
	if contest.Rules == nil {
		contest.Rules = model.StringList{}
	}
	if contest.Prizes == nil {
		contest.Prizes = model.Prizes{}
	}
	if contest.Timelines == nil {
		contest.Timelines = model.Timelines{}
	}
	data := gin.H{
		"id":          contest.ID,
		"name":        contest.Name,
		"description": contest.Description,
		"start":       contest.Start,
		"duration":    int64(contest.Duration.Seconds()),
		"rules":       contest.Rules,
		"prizes":      contest.Prizes,
		"size":        contest.Size,
		"timelines":   contest.Timelines,
		"teams":       contestView.TeamCount,
		"users":       contestView.UserCount,
		"notices":     contestView.NoticeCount,
		"prefix":      contest.Prefix,
		"victims":     contest.Victims,
		"picture":     contest.Picture,
		"hidden":      contest.Hidden,
		"blood":       contest.Blood,
	}
	if contestView.StatsReady {
		data["highest"] = contestView.Highest
		data["solved"] = contestView.SolvedCount
	}
	if admin {
		data["captcha"] = contest.Captcha
	}
	return data
}
