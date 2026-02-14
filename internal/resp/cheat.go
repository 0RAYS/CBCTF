package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetCheatResp(cheat model.Cheat) gin.H {
	return gin.H{
		"id":          cheat.ID,
		"contest_id":  cheat.ContestID,
		"model":       cheat.Model,
		"magic":       cheat.Magic,
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
