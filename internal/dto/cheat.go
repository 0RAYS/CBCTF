package dto

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

type GetCheatsForm struct {
	ListModelsForm
	Type       model.CheatType `form:"type" json:"type" binding:"omitempty,oneof=suspicious cheater pass"`
	ReasonType string          `form:"reason_type" json:"reason_type" binding:"omitempty,oneof=same_device same_web_ip same_victim_ip wrong_flag token_magic"`
}

func (f *GetCheatsForm) Validate(ctx *gin.Context) model.RetVal {
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 10
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
	return model.SuccessRetVal()
}

type UpdateCheatForm struct {
	Reason  *string          `form:"reason" json:"reason"`
	Type    *model.CheatType `form:"type" json:"type" binding:"omitempty,oneof=suspicious cheater pass"`
	Checked *bool            `form:"checked" json:"checked"`
	Comment *string          `form:"comment" json:"comment"`
}
