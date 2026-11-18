package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/view"
)

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
		"generator_image":  challenge.GeneratorImage,
		"flags":            flags,
		"docker_compose":   challengeView.DockerCompose,
		"network_policies": challenge.NetworkPolicies,
		"file":             challengeView.FileName,
	}
	markUnavailable(data, challengeView.Unavailable)
	return data
}

func GetSimpleChallengeResp(challengeView view.SimpleChallengeView) gin.H {
	challenge := challengeView.Challenge
	return gin.H{
		"id":               challenge.RandID,
		"name":             challenge.Name,
		"description":      challenge.Description,
		"category":         challenge.Category,
		"type":             challenge.Type,
		"generator_image":  challenge.GeneratorImage,
		"network_policies": challenge.NetworkPolicies,
	}
}

// GetChallengeUsageResp 返回引用了题库题目的比赛题目实例。
func GetChallengeUsageResp(usageView view.ChallengeUsageView) gin.H {
	contestChallenge := usageView.ContestChallenge
	return gin.H{
		"contest_id":   usageView.ContestID,
		"contest_name": usageView.ContestName,
		"id":           contestChallenge.RandID,
		"name":         contestChallenge.Name,
		"type":         contestChallenge.Type,
		"category":     contestChallenge.Category,
		"hidden":       contestChallenge.Hidden,
		"solvers":      usageView.Solvers,
	}
}
