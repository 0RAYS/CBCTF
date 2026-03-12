package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"slices"

	"github.com/gin-gonic/gin"
)

func GetContestChallengeImage(ctx *gin.Context) {
	contest := middleware.GetContest(ctx)
	contestChallengeImageList, ret := service.GetContestChallengeImageList(db.DB, contest)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}

	nodeImageMap, ret := service.GetNodeImageList()
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, contestChallengeImage := range contestChallengeImageList {
		status := make([]map[string]any, 0)
		for node, nodeImage := range nodeImageMap {
			if slices.Contains(nodeImage, contestChallengeImage) {
				status = append(status, map[string]any{
					"node":   node,
					"status": true,
				})
			} else {
				status = append(status, map[string]any{
					"node":   node,
					"status": false,
				})
			}
		}
		data = append(data, gin.H{contestChallengeImage: status})
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
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, ret)
}
