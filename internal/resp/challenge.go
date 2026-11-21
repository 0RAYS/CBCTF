package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
	"CBCTF/internal/view"
)

// GetChallengeResp 返回后台题目详情。
func GetChallengeResp(challengeView view.ChallengeView) gin.H {
	challenge := challengeView.Challenge
	flags := make([]gin.H, 0)
	for _, flag := range challengeView.Flags {
		flags = append(flags, gin.H{"id": flag.ID, "value": flag.Value})
	}
	data := gin.H{
		"id":               challenge.RandID,
		"name":             challenge.Name,
		"description":      challenge.Description,
		"category":         challenge.Category,
		"type":             challenge.Type,
		"contest_id":       challenge.ContestID.V,
		"contest_name":     challengeView.ContestName,
		"hidden":           challenge.Hidden,
		"attempt":          challenge.Attempt,
		"hints":            challenge.Hints,
		"tags":             challenge.Tags,
		"generator_image":  challenge.GeneratorImage,
		"flags":            flags,
		"docker_compose":   challengeView.DockerCompose,
		"network_policies": challenge.NetworkPolicies,
		"file":             challengeView.FileName,
	}
	markUnavailable(data, challengeView.Unavailable)
	return data
}

func getChallengeBaseResp(challenge model.Challenge) gin.H {
	var score float64
	var solvers int64
	for _, flag := range challenge.ChallengeFlags {
		if !flag.Active {
			continue
		}
		score += flag.CurrentScore
		solvers += flag.Solvers
	}
	return gin.H{
		"id":          challenge.RandID,
		"name":        challenge.Name,
		"description": challenge.Description,
		"attempt":     challenge.Attempt,
		"type":        challenge.Type,
		"category":    challenge.Category,
		"hidden":      challenge.Hidden,
		"score":       score,
		"solvers":     solvers,
		"hints":       challenge.Hints,
		"tags":        challenge.Tags,
		"contest_id":  challenge.ContestID.V,
	}
}

func GetChallengeRuntimeResp(challengeView view.ChallengeRuntimeView) gin.H {
	data := getChallengeBaseResp(challengeView.Challenge)
	data["attempts"] = challengeView.Attempts
	data["init"] = challengeView.Init
	data["solved"] = challengeView.Solved
	data["remote"] = GetVictimStatusResp(challengeView.Remote)
	data["file"] = challengeView.FileName
	return data
}

func GetAdminChallengeResp(challenge model.Challenge) gin.H {
	return getChallengeBaseResp(challenge)
}

func GetChallengeStatusResp(status view.ChallengeStatusView) gin.H {
	return gin.H{
		"attempts": status.Attempts,
		"init":     status.Init,
		"solved":   status.Solved,
		"remote":   GetVictimStatusResp(status.Remote),
		"file":     status.FileName,
	}
}
