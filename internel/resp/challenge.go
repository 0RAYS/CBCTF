package resp

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

func GetChallengeResp(challenge model.Challenge) gin.H {
	return gin.H{
		"id":        challenge.ID,
		"name":      challenge.Name,
		"desc":      challenge.Desc,
		"category":  challenge.Category,
		"type":      challenge.Type,
		"generator": challenge.Generator,
		"flags":     challenge.Flags,
		"dockers":   challenge.Dockers,
	}
}
