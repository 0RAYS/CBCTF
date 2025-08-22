package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetCheatResp(cheat model.Cheat) gin.H {
	return gin.H{
		"id":                   cheat.ID,
		"user_id":              cheat.UserID.V,
		"team_id":              cheat.TeamID.V,
		"contest_id":           cheat.ContestID.V,
		"contest_challenge_id": cheat.ContestChallengeID.V,
		"contest_flag_id":      cheat.ContestFlagID.V,
		"magic":                cheat.Magic,
		"ip":                   cheat.IP,
		"reason":               cheat.Reason,
		"type":                 cheat.Type,
		"checked":              cheat.Checked,
		"hash":                 cheat.Hash,
		"comment":              cheat.Comment,
		"references":           cheat.References,
		"created_at":           cheat.CreatedAt,
	}
}
