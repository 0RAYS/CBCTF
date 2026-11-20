package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
)

// GetCheatResp 返回线索详情，包含完整证据链。
func GetCheatResp(cheat model.Cheat) gin.H {
	data := GetCheatListResp(cheat)
	data["evidence"] = cheat.Evidence
	return data
}

// GetCheatListResp 返回线索列表项，只带证据摘要，避免列表被证据撑大。
func GetCheatListResp(cheat model.Cheat) gin.H {
	return gin.H{
		"id":               cheat.ID,
		"contest_id":       cheat.ContestID,
		"model":            cheat.Model,
		"ip":               cheat.IP,
		"reason":           cheat.Reason,
		"reason_type":      cheat.ReasonType,
		"type":             cheat.Type,
		"checked":          cheat.Checked,
		"hash":             cheat.Hash,
		"comment":          cheat.Comment,
		"time":             cheat.Time,
		"evidence_summary": cheat.Evidence.Summary(),
		"reviewed_by":      cheat.ReviewedBy,
		"reviewed_at":      cheat.ReviewedAt,
	}
}
