package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
	"CBCTF/internal/view"
)

func GetContestResp(contestView view.ContestView, admin bool) gin.H {
	contest := contestView.Contest
	if contest.Rules == nil {
		contest.Rules = []string{}
	}
	if contest.Prizes == nil {
		contest.Prizes = []model.Prize{}
	}
	if contest.Timelines == nil {
		contest.Timelines = []model.Timeline{}
	}
	data := gin.H{
		"id":                 contest.ID,
		"name":               contest.Name,
		"description":        contest.Description,
		"start":              contest.Start,
		"end":                contest.End(),
		"duration":           int64(contest.Duration.Seconds()),
		"rules":              contest.Rules,
		"prizes":             contest.Prizes,
		"size":               contest.Size,
		"timelines":          contest.Timelines,
		"teams":              contestView.TeamCount,
		"users":              contestView.UserCount,
		"notices":            contestView.NoticeCount,
		"prefix":             contest.Prefix,
		"victims":            contest.Victims,
		"picture":            contest.Picture,
		"hidden":             contest.Hidden,
		"blood":              contest.Blood,
		"paused":             contest.Paused,
		"leaderboard_hidden": contest.LeaderboardHidden,
	}
	if contestView.StatsReady {
		// 排行榜黑灯期间最高分同样属于榜单信息，选手不可见。
		if admin || !contest.LeaderboardHidden {
			data["highest"] = contestView.Highest
		}
		data["solved"] = contestView.SolvedCount
	}
	if admin {
		data["captcha"] = contest.Captcha
	}
	markUnavailable(data, contestView.Unavailable)
	return data
}
