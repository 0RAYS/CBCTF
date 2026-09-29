package router

import (
	"sort"

	"github.com/gin-gonic/gin"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
)

func formatNodeImages(nodeImageMap map[string][]string, targetImages []string) []gin.H {
	targetSet := make(map[string]struct{}, len(targetImages))
	for _, image := range targetImages {
		targetSet[image] = struct{}{}
	}
	data := make([]gin.H, 0, len(nodeImageMap))
	for node, images := range nodeImageMap {
		current := make([]string, 0)
		for _, image := range images {
			if _, ok := targetSet[image]; ok {
				current = append(current, image)
			}
		}
		sort.Strings(current)
		data = append(data, gin.H{
			"node":   node,
			"images": current,
		})
	}
	sort.Slice(data, func(i, j int) bool { return data[i]["node"].(string) < data[j]["node"].(string) })
	return data
}

func GetImages(ctx *gin.Context) {
	nodeImageMap, ret := service.ListNodeImages(ctx.Request.Context())
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}

	targetImages, ret := service.ListChallengeImages(db.DB.WithContext(ctx.Request.Context()))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{
		"nodes":         formatNodeImages(nodeImageMap, targetImages),
		"target_images": targetImages,
	}))
}

func GetContestChallengeImage(ctx *gin.Context) {
	nodeImageMap, ret := service.ListNodeImages(ctx.Request.Context())
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}

	targetImages, ret := service.ListContestChallengeImages(db.DB.WithContext(ctx.Request.Context()), middleware.GetContest(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{
		"nodes":         formatNodeImages(nodeImageMap, targetImages),
		"target_images": targetImages,
	}))
}

func PullImages(ctx *gin.Context) {
	var form dto.PullImageForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.PullImageEventType)
	ret := service.PullContestChallengeImage(ctx.Request.Context(), form)
	ctx.Set(middleware.CTXEventSuccessKey, ret.OK)
	resp.JSON(ctx, ret)
}
