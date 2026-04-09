package router

import (
	"CBCTF/internal/dto"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"

	"github.com/gin-gonic/gin"
)

func GetAllowQueryModels(ctx *gin.Context) {
	data := gin.H{}
	for name, info := range service.GetAllowQueryModels() {
		data[name] = gin.H{
			"query":  info.Query,
			"search": info.Search,
		}
	}
	resp.JSON(ctx, model.SuccessRetVal(data))
}

func Search(ctx *gin.Context) {
	var form dto.SearchModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ms, count, ret := service.SearchModels(nil, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"count": count, "models": ms}))
}
