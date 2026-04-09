package resp

import (
	"CBCTF/internal/model"
	"CBCTF/internal/view"
	"slices"
	"strings"

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
	slices.SortFunc(data, func(a, b gin.H) int {
		return strings.Compare(a["category"].(string), b["category"].(string))
	})
	return data
}

func GetTeamResp(teamView view.TeamView, isAdmin bool) gin.H {
	team := teamView.Team
	data := gin.H{
		"id":          team.ID,
		"contest_id":  team.ContestID,
		"name":        team.Name,
		"score":       team.Score,
		"picture":     team.Picture,
		"last":        team.Last,
		"rank":        team.Rank,
		"users":       teamView.UserCount,
		"description": team.Description,
		"captain_id":  team.CaptainID,
		"banned":      team.Banned,
		"hidden":      team.Hidden,
	}
	if isAdmin {
		data["captcha"] = team.Captcha
	}
	return data
}
