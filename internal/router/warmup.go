package router

import (
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	db "CBCTF/internal/repo"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"CBCTF/internal/websocket"
	wsm "CBCTF/internal/websocket/model"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
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
					"status": true,
				})
			} else {
				status = append(status, map[string]any{
					"node":   node,
					"status": false,
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
	ctx.Set(middleware.CTXEventTypeKey, model.PullImageEventType)
	_, msg := service.WarmUpContestChallengeImage(form)
	ctx.Set(middleware.CTXEventSuccessKey, true)
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
		info["remote"] = victim.RemoteAddr()
		info["remaining"] = victim.Remaining().Seconds()
		info["team"] = victim.Team.Name
		info["user"] = victim.User.Name
		info["challenge"] = victim.ContestChallenge.Name
		data = append(data, info)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": gin.H{"victims": data, "count": count}})
}

func StartContestVictims(ctx *gin.Context) {
	var form f.StartContestVictimsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.StartVictimEventType)
	contest := middleware.GetContest(ctx)
	go func(ctx *gin.Context) {
		if ok, _ := service.StartContestVictims(db.DB.WithContext(ctx), contest, form); !ok {
			websocket.Send(true, middleware.GetSelfID(ctx), wsm.ErrorLevel, wsm.StartVictimType, "Victims Warmup", "Failed")
			return
		}
		websocket.Send(true, middleware.GetSelfID(ctx), wsm.SuccessLevel, wsm.StartVictimType, "Victims Warmup", "Done")
	}(ctx.Copy())
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
}

func StopContestVictims(ctx *gin.Context) {
	var form f.StopContestVictimsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.StopVictimEventType)
	go func(ctx *gin.Context) {
		if ok, _ := service.StopContestVictims(db.DB.WithContext(ctx), form); !ok {
			websocket.Send(true, middleware.GetSelfID(ctx), wsm.ErrorLevel, wsm.StopVictimType, "Victims Stop", "Failed")
			return
		}
		websocket.Send(true, middleware.GetSelfID(ctx), wsm.SuccessLevel, wsm.StopVictimType, "Victims Stop", "Done")
	}(ctx.Copy())
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
}
