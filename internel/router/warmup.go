package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/middleware"
	db "CBCTF/internel/repo"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
)

func GetContestChallengeImage(ctx *gin.Context) {
	contest := middleware.GetContest(ctx)
	contestChallengeImageList, ok, msg := service.GetContestChallengeImageList(db.DB.WithContext(ctx), contest)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}

	nodeImageMap, ok, msg := service.GetNodeImageList()
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, contestChallengeImage := range contestChallengeImageList {
		status := make([]map[string]any, 0)
		for node, nodeImage := range nodeImageMap {
			if slices.Contains(nodeImage, contestChallengeImage) {
				status = append(status, map[string]any{
					"node":   node,
					"status": "exists",
				})
			} else {
				status = append(status, map[string]any{
					"node":   node,
					"status": "not exists",
				})
			}
		}
		data = append(data, gin.H{
			contestChallengeImage: status,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": data})
}

func WarmUpContestChallengeImage(ctx *gin.Context) {
	var form f.WarmUpImageForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	_, msg := service.WarmUpContestChallengeImage(form)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
