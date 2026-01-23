package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
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
	contestChallengeImageList, ret := service.GetContestChallengeImageList(db.DB, contest)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}

	nodeImageMap, ret := service.GetNodeImageList()
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
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
	ctx.JSON(http.StatusOK, model.SuccessRetVal(data))
}

func WarmUpContestChallengeImage(ctx *gin.Context) {
	var form f.WarmUpImageForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.PullImageEventType)
	ret := service.WarmUpContestChallengeImage(form)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, ret)
}

func GetContestVictims(ctx *gin.Context) {
	var form f.GetContestVictimsForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	contest := middleware.GetContest(ctx)
	victims, count, _ := service.GetContestVictims(db.DB, contest, form)
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
	total, ret := db.InitVictimRepo(db.DB).Count(db.CountOptions{
		Conditions: map[string]any{"contest_id": contest.ID}, Deleted: true,
	})
	if !ret.OK {
		total = count
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"victims": data, "count": total, "running": count}))
}

func StartContestVictims(ctx *gin.Context) {
	var form f.StartContestVictimsForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.StartVictimEventType)
	contest := middleware.GetContest(ctx)
	go func(selfID uint) {
		if ret := service.StartContestVictims(db.DB, contest, form); !ret.OK {
			websocket.Send(true, selfID, wsm.ErrorLevel, wsm.StartVictimWSType, "Victims Warmup", "Failed")
			return
		}
		websocket.Send(true, selfID, wsm.SuccessLevel, wsm.StartVictimWSType, "Victims Warmup", "Done")
	}(middleware.GetSelfID(ctx))
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal())
}

func StopContestVictims(ctx *gin.Context) {
	var form f.StopContestVictimsForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.StopVictimEventType)
	go func(selfID uint) {
		if ret := service.StopContestVictims(db.DB, form); !ret.OK {
			websocket.Send(true, selfID, wsm.ErrorLevel, wsm.StopVictimWSType, "Victims Stop", "Failed")
			return
		}
		websocket.Send(true, selfID, wsm.SuccessLevel, wsm.StopVictimWSType, "Victims Stop", "Done")
	}(middleware.GetSelfID(ctx))
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal())
}
