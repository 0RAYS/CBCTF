package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"sort"

	"github.com/gin-gonic/gin"
)

func formatNodeImages(nodeImageMap map[string][]string) []gin.H {
	data := make([]gin.H, 0, len(nodeImageMap))
	for node, images := range nodeImageMap {
		data = append(data, gin.H{
			"node":   node,
			"images": images,
		})
	}
	return data
}

func GetImages(ctx *gin.Context) {
	nodeImageMap, ret := service.GetNodeImageList()
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}

	resp.JSON(ctx, model.SuccessRetVal(gin.H{
		"nodes": formatNodeImages(nodeImageMap),
		"target_images": func(nodeImageMap map[string][]string) []string {
			imageSet := make(map[string]struct{})
			images := make([]string, 0)
			for _, nodeImages := range nodeImageMap {
				for _, image := range nodeImages {
					if image == "" {
						continue
					}
					if _, ok := imageSet[image]; ok {
						continue
					}
					imageSet[image] = struct{}{}
					images = append(images, image)
				}
			}
			sort.Strings(images)
			return images
		}(nodeImageMap),
	}))
}

func GetContestChallengeImage(ctx *gin.Context) {
	nodeImageMap, ret := service.GetNodeImageList()
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}

	targetImages, ret := service.GetContestChallengeImageList(db.DB, middleware.GetContest(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	nodeImageMap = func(nodeImageMap map[string][]string, targetImages []string) map[string][]string {
		targetSet := make(map[string]struct{}, len(targetImages))
		for _, image := range targetImages {
			targetSet[image] = struct{}{}
		}

		filtered := make(map[string][]string, len(nodeImageMap))
		for node, images := range nodeImageMap {
			current := make([]string, 0, len(images))
			for _, image := range images {
				if _, ok := targetSet[image]; ok {
					current = append(current, image)
				}
			}
			sort.Strings(current)
			filtered[node] = current
		}
		return filtered
	}(nodeImageMap, targetImages)

	resp.JSON(ctx, model.SuccessRetVal(gin.H{
		"nodes":         formatNodeImages(nodeImageMap),
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
	ret := service.PullContestChallengeImage(form)
	ctx.Set(middleware.CTXEventSuccessKey, ret.OK)
	resp.JSON(ctx, ret)
}
