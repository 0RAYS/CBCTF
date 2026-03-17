package resp

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetContestResp(contest model.Contest, admin bool) gin.H {
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
		"teams":       db.InitContestRepo(db.DB).CountAssociation(contest, "Teams"),
		"users":       db.InitContestRepo(db.DB).CountAssociation(contest, "Users"),
		"notices":     db.InitContestRepo(db.DB).CountAssociation(contest, "Notices"),
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
