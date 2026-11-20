package dto

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
)

type ListModelsForm struct {
	Offset int `form:"offset" json:"offset" binding:"gte=0"`
	Limit  int `form:"limit" json:"limit" binding:"gte=0,lte=100"`
}

// ApplyPaginationDefaults keeps the default page size when the query does not
// set one.
func (f *ListModelsForm) ApplyPaginationDefaults(ctx *gin.Context, defaultLimit int) {
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = defaultLimit
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
}

func (f *ListModelsForm) Validate(ctx *gin.Context) model.RetVal {
	f.ApplyPaginationDefaults(ctx, 10)
	return model.SuccessRetVal()
}

// LeaderboardForm 排行榜查询参数，ShowHidden 决定隐藏队伍是否参与排名与名次计算。
type LeaderboardForm struct {
	ListModelsForm
	ShowHidden bool `form:"show_hidden" json:"show_hidden"`
}
