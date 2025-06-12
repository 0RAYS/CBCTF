package resp

import (
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
)

func GetCheatResp(cheat model.Cheat) gin.H {
	return gin.H{
		"id":                   cheat.ID,
		"user_id":              cheat.UserID,
		"team_id":              cheat.TeamID,
		"contest_id":           cheat.ContestID,
		"contest_challenge_id": cheat.ContestChallengeID,
		"contest_flag_id":      cheat.ContestFlagID,
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
