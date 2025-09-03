package resp

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

// GetSolvedStateResp model.ContestFlag Preload model.ContestChallenge
func GetSolvedStateResp(solved []model.ContestFlag, all []model.ContestFlag) []gin.H {
	categories := make(map[uint]string)
	for _, v := range all {
		categories[v.ContestChallengeID] = v.ContestChallenge.Category
	}
	allCount := make(map[string]int64)
	for _, v := range all {
		allCount[v.ContestChallenge.Category] += 1
	}
	solvedCount := make(map[string]int64)
	for _, flag := range solved {
		solvedCount[categories[flag.ContestChallengeID]] += 1
	}
	data := make([]gin.H, 0)
	for k, v := range allCount {
		if _, ok := solvedCount[k]; !ok {
			solvedCount[k] = 0
		}
		data = append(data, gin.H{"category": k, "solved": solvedCount[k], "all": v})
	}
	return data
}

func GetTeamResp(team model.Team) gin.H {
	data := gin.H{
		"id":         team.ID,
		"contest_id": team.ContestID,
		"name":       team.Name,
		"score":      team.Score,
		"avatar":     team.Avatar,
		"last":       team.Last,
		"rank":       team.Rank,
		"users":      db.InitTeamRepo(db.DB).CountAssociation(team, "Users"),
		"desc":       team.Desc,
		"captain_id": team.CaptainID,
	}
	return data
}
