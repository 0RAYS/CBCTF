package dto

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

type SearchIP struct {
	IP string `form:"ip" json:"ip" binding:"required,ip|cidr"`
}

// ListModelsForm for get models list
type ListModelsForm struct {
	Offset int               `form:"offset" json:"offset" binding:"gte=0"`
	Limit  int               `form:"limit" json:"limit" binding:"gte=0,lte=100"`
	Sort   map[string]string `form:"sort" json:"sort"`
	Search map[string]string `form:"search" json:"search"`
}

func (f *ListModelsForm) Validate(ctx *gin.Context) model.RetVal {
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
	OldPassword string `form:"old" json:"old" binding:"required,nefield=NewPassword"`
	NewPassword string `form:"new" json:"new" binding:"required,nefield=OldPassword"`
}
