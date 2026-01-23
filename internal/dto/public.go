package dto

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

// ListModelsForm for get models list
type ListModelsForm struct {
	Offset int               `form:"offset" json:"offset"`
	Limit  int               `form:"limit" json:"limit"`
	Sort   map[string]string `form:"sort" json:"sort"`
	Search map[string]string `form:"search" json:"search"`
}

func (f *ListModelsForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	if f.Limit > 100 || f.Limit < 0 {
		f.Limit = 15
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 10
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
	f.Sort = ctx.QueryMap("sort")
	f.Search = ctx.QueryMap("search")
	return model.SuccessRetVal()
}

// ChangePasswordForm for user or admin change password
type ChangePasswordForm struct {
	OldPassword string `form:"old" json:"old" binding:"required"`
	NewPassword string `form:"new" json:"new" binding:"required"`
}

func (f *ChangePasswordForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	if f.OldPassword == f.NewPassword {
		return model.RetVal{Msg: i18n.Model.User.SamePassword}
	}
	return model.SuccessRetVal()
}
