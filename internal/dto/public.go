package dto

import (
	"CBCTF/internal/model"
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
	ListModelsForm
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
		f.Sort[strings.ToLower(k)] = strings.ToLower(v)
	}
	f.Search = make(map[string]string)
	for k, v := range ctx.QueryMap("search") {
		f.Search[strings.ToLower(k)] = v
	}
	return model.SuccessRetVal()
}
