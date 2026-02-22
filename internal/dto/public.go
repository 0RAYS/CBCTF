package dto

import (
	"CBCTF/internal/model"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

// ListModelsForm for get models list
type ListModelsForm struct {
	Offset int `form:"offset" json:"offset" binding:"gte=0"`
	Limit  int `form:"limit" json:"limit" binding:"gte=0,lte=100"`
}

func (f *ListModelsForm) Validate(ctx *gin.Context) model.RetVal {
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 10
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
	return model.SuccessRetVal()
}

type SearchModelsForm struct {
	Offset int               `form:"offset" json:"offset" binding:"gte=0"`
	Limit  int               `form:"limit" json:"limit" binding:"gte=0,lte=100"`
	Model  string            `form:"model" json:"model"`
	Sort   map[string]string `form:"sort" json:"sort"`
	Search map[string]string `form:"search" json:"search"`
}

func (f *SearchModelsForm) Validate(ctx *gin.Context) model.RetVal {
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 10
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
	f.Sort = make(map[string]string)
	for k, v := range ctx.QueryMap("sort") {
		if slices.Contains([]string{"desc", "asc"}, strings.ToLower(v)) {
			f.Sort[strings.ToLower(k)] = strings.ToLower(v)
		}
	}
	f.Search = make(map[string]string)
	for k, v := range ctx.QueryMap("search") {
		f.Search[strings.ToLower(k)] = v
	}
	return model.SuccessRetVal()
}

// ChangePasswordForm for user or admin change password
type ChangePasswordForm struct {
	OldPassword string `form:"old" json:"old" binding:"required,nefield=NewPassword"`
	NewPassword string `form:"new" json:"new" binding:"required,nefield=OldPassword"`
}
