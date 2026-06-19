package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/k8s"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func StartVictim(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.StartVictimEventType)
	user := middleware.GetSelf(ctx)
	team := middleware.GetTeam(ctx)
	contest := middleware.GetContest(ctx)
	challenge := middleware.GetChallenge(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	ret := service.StartVictim(db.DB, user.ID, team.ID, contest.ID, contestChallenge.ID, challenge.ID)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func ExtendVictimDuration(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.ExtendVictimEventType)
	victim, ret := service.ExtendVictimDuration(db.DB, middleware.GetTeam(ctx), middleware.GetChallenge(ctx))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(gin.H{
		"target":    victim.RemoteAddr(),
		"duration":  victim.Duration.Seconds(),
		"remaining": victim.Remaining().Seconds(),
		"status":    "Running",
	}))
}

func StopVictim(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.StopVictimEventType)
	if ret := service.StopAliveVictim(db.DB, middleware.GetTeam(ctx), middleware.GetChallenge(ctx)); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}

func GetVictimHistories(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	victims, count, ret := service.ListVictimHistories(db.DB, middleware.GetTeam(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, victim := range victims {
		data = append(data, resp.GetVictimResp(victim))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"victims": data, "count": count}))
}

func GetVictims(ctx *gin.Context) {
	var form dto.GetVictimsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	contest := middleware.GetContest(ctx)
	victims, running, total, ret := service.GetVictims(db.DB, contest, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
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
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"victims": data, "count": total, "running": running}))
}

func StartVictims(ctx *gin.Context) {
	var form dto.StartVictimsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.StartVictimEventType)
	contest := middleware.GetContest(ctx)
	ret := service.StartVictims(db.DB, contest, form)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

func StopVictims(ctx *gin.Context) {
	var form dto.StopVictimsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.StopVictimEventType)
	go service.StopVictims(db.DB, form)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}

// GetVictimPods 列出指定 victim 关联的 Pods（pending/running/terminating 状态）
func GetVictimPods(ctx *gin.Context) {
	victim := middleware.GetVictim(ctx)
	switch victim.Status {
	case model.PendingVictimStatus, model.RunningVictimStatus, model.TerminatingVictimStatus:
	default:
		resp.JSON(ctx, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "Pod", "Error": "victim is not active"}})
		return
	}

	labels := k8s.VictimLabels(victim)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	podList, ret := k8s.ListPods(ctxTimeout, labels)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}

	pods := make([]gin.H, 0, len(podList.Items))
	for _, pod := range podList.Items {
		containers := make([]string, 0, len(pod.Spec.InitContainers)+len(pod.Spec.Containers))
		for _, c := range pod.Spec.InitContainers {
			if c.Name == k8s.CaptureContainerName {
				continue
			}
			containers = append(containers, c.Name)
		}
		for _, c := range pod.Spec.Containers {
			if c.Name == k8s.CaptureContainerName {
				continue
			}
			containers = append(containers, c.Name)
		}
		if len(containers) == 0 {
			continue
		}
		pods = append(pods, gin.H{
			"name":       pod.Name,
			"status":     string(pod.Status.Phase),
			"containers": containers,
		})
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"pods": pods}))
}

// GetVictimPodLogs 获取指定 victim 的某个 Pod 日志（pending/running/terminating 状态）
func GetVictimPodLogs(ctx *gin.Context) {
	victim := middleware.GetVictim(ctx)
	switch victim.Status {
	case model.PendingVictimStatus, model.RunningVictimStatus, model.TerminatingVictimStatus:
	default:
		resp.JSON(ctx, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "PodLog", "Error": "victim is not active"}})
		return
	}

	var form dto.GetVictimPodLogsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	podList, ret := k8s.ListPods(ctxTimeout, k8s.VictimLabels(victim))
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	podFound := false
	containerFound := false
	for _, pod := range podList.Items {
		if pod.Name != form.PodName {
			continue
		}
		podFound = true
		for _, c := range pod.Spec.InitContainers {
			if c.Name == form.Container && c.Name != k8s.CaptureContainerName {
				containerFound = true
				break
			}
		}
		if containerFound {
			break
		}
		for _, c := range pod.Spec.Containers {
			if c.Name == form.Container && c.Name != k8s.CaptureContainerName {
				containerFound = true
				break
			}
		}
		break
	}
	if !podFound {
		resp.JSON(ctx, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "PodLog", "Error": "pod does not belong to victim"}})
		return
	}
	if !containerFound {
		resp.JSON(ctx, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "PodLog", "Error": "container does not belong to victim"}})
		return
	}
	logs, ret := k8s.GetPodLogs(ctxTimeout, form.PodName, form.Container, form.Lines)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"logs": logs}))
}
