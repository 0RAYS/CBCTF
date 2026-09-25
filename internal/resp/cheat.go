package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
)

func GetCheatResp(cheat model.Cheat) gin.H {
	return gin.H{
		"id":          cheat.ID,
		"contest_id":  cheat.ContestID,
		"model":       cheat.Model,
		"ip":          cheat.IP,
		"reason":      cheat.Reason,
		"reason_type": cheat.ReasonType,
		"type":        cheat.Type,
		"checked":     cheat.Checked,
		"hash":        cheat.Hash,
		"comment":     cheat.Comment,
		"time":        cheat.Time,
	}
}
