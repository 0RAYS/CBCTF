package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/middleware"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
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
		status := make([]map[string]string, 0)
		for node, nodeImage := range nodeImageMap {
			if slices.Contains(nodeImage, contestChallengeImage) {
				status = append(status, map[string]string{
					"node":   node,
					"status": "exists",
				})
			} else {
				status = append(status, map[string]string{
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
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	_, msg := service.WarmUpContestChallengeImage(form)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func GetContestVictims(ctx *gin.Context) {
	var form f.GetContestVictimsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	victims, count, _, _ := service.GetContestVictims(db.DB.WithContext(ctx), contest, form)
	data := make([]gin.H, 0)
	for _, victim := range victims {
		info := resp.GetVictimResp(victim)
		remoteAddrL := make([]string, 0)
		for _, pod := range victim.Pods {
			remoteAddrL = append(remoteAddrL, pod.RemoteAddr()...)
		}
		info["remote"] = remoteAddrL
		info["remaining"] = victim.Remaining().Seconds()
		info["team"] = victim.Team.Name
		info["user"] = victim.User.Name
		info["challenge"] = victim.ContestChallenge.Name
		data = append(data, info)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": gin.H{"victims": data, "count": count}})
}
