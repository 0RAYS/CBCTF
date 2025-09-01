package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetFileResp(file model.File) gin.H {
	return gin.H{
		"id":           file.RandID,
		"filename":     file.Filename,
		"type":         file.Type,
		"hash":         file.Hash,
		"size":         file.Size,
		"date":         file.CreatedAt,
		"suffix":       file.Suffix,
		"admin_id":     file.AdminID.V,
		"user_id":      file.UserID.V,
		"team_id":      file.TeamID.V,
		"contest_id":   file.ContestID.V,
		"oauth_id":     file.OauthID.V,
		"challenge_id": file.ChallengeID.V,
	}
}
