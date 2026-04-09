package resp

import (
	"CBCTF/internal/view"

	"github.com/gin-gonic/gin"
)

func GetChallengeResp(challengeView view.ChallengeView) gin.H {
	challenge := challengeView.Challenge
	flags := make([]gin.H, 0)
	for _, flag := range challengeView.Flags {
		flags = append(flags, gin.H{"id": flag.ID, "value": flag.Value})
	}
	return gin.H{
		"id":               challenge.RandID,
		"name":             challenge.Name,
		"description":      challenge.Description,
		"category":         challenge.Category,
		"type":             challenge.Type,
		"generator_image":  challenge.GeneratorImage,
		"flags":            flags,
		"docker_compose":   challengeView.DockerCompose,
		"options":          challenge.Options,
		"network_policies": challenge.NetworkPolicies,
		"file":             challengeView.FileName,
	}
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
		"options":          challenge.Options,
		"network_policies": challenge.NetworkPolicies,
	}
}
