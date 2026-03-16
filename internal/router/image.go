package router

import (
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"

	"github.com/gin-gonic/gin"
)

func GetContestChallengeImage(ctx *gin.Context) {
	nodeImageMap, ret := service.GetNodeImageList()
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}

	data := make([]gin.H, 0, len(nodeImageMap))
	for node, images := range nodeImageMap {
		data = append(data, gin.H{
			"node":   node,
			"images": images,
		})
	}
	resp.JSON(ctx, model.SuccessRetVal(data))
}

func PullContestChallengeImage(ctx *gin.Context) {
	var form dto.PullImageForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.PullImageEventType)
	ret := service.PullContestChallengeImage(form)
	ctx.Set(middleware.CTXEventSuccessKey, ret.OK)
	resp.JSON(ctx, ret)
}
