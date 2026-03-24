package resp

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetContestResp(contest model.Contest, admin bool) gin.H {
	contestRepo := db.InitContestRepo(db.DB)
	teams, _ := contestRepo.CountTeams(contest.ID)
	users, _ := contestRepo.CountUsers(contest.ID)
	notices, _ := contestRepo.CountNotices(contest.ID)
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
		"teams":       teams,
		"users":       users,
		"notices":     notices,
		"prefix":      contest.Prefix,
		"victims":     contest.Victims,
		"picture":     contest.Picture,
		"hidden":      contest.Hidden,
		"blood":       contest.Blood,
	}
	if admin {
		data["captcha"] = contest.Captcha
	}
	return data
}
