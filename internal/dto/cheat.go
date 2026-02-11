package dto

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

type GetCheatsForm struct {
	Offset int    `form:"offset" json:"offset" binding:"gte=0"`
	Limit  int    `form:"limit" json:"limit" binding:"gte=0,lte=100"`
	Type   string `form:"type" json:"type" binding:"omitempty,oneof=suspicious cheater pass"`
}

func (f *GetCheatsForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

type UpdateCheatForm struct {
	Reason  *string `form:"reason" json:"reason"`
	Type    *string `form:"type" json:"type" binding:"omitempty,oneof=suspicious cheater pass"`
	Checked *bool   `form:"checked" json:"checked"`
	Comment *string `form:"comment" json:"comment"`
}

func (f *UpdateCheatForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
