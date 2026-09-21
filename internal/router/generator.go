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

func GetGenerators(ctx *gin.Context) {
	var form dto.ListGeneratorsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	generators, count, ret := service.ListGenerators(db.DB, middleware.GetContest(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, generator := range generators {
		data = append(data, resp.GetGeneratorResp(generator))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"count": count, "generators": data}))
}

func StartGenerator(ctx *gin.Context) {
	var form dto.StartGeneratorsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.StartGeneratorEventType)
	contest := middleware.GetContest(ctx)
	go service.StartGenerators(db.TaskDB, contest.ID, form)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}

func StopGenerator(ctx *gin.Context) {
	var form dto.StopGeneratorsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.StopGeneratorEventType)
	go service.StopGenerators(db.TaskDB, middleware.GetContest(ctx).ID, form)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}

// GetGeneratorLogs 获取指定 generator 的 Pod 日志（pending/running/terminating 状态）
func GetGeneratorLogs(ctx *gin.Context) {
	generator := middleware.GetGenerator(ctx)
	switch generator.Status {
	case model.PendingGeneratorStatus, model.RunningGeneratorStatus, model.TerminatingGeneratorStatus:
	default:
		resp.JSON(ctx, model.RetVal{Msg: i18n.K8S.GetError, Attr: map[string]any{"Model": "PodLog", "Error": "generator is not active"}})
		return
	}

	var form dto.GetGeneratorLogsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}

	ctxTimeout, cancel := context.WithTimeout(ctx.Request.Context(), 30*time.Second)
	defer cancel()
	pod, ret := k8s.GetPod(ctxTimeout, generator.Name)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	if !k8s.GeneratorOwnsPod(generator, pod) {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	logs, ret := k8s.GetPodLogs(ctxTimeout, generator.Name, "generator", form.Lines)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"logs": logs}))
}

// GetGeneratorStatus shares the admin/contest authorization chain with logs.
func GetGeneratorStatus(ctx *gin.Context) {
	generator := middleware.GetGenerator(ctx)
	timeout, cancel := context.WithTimeout(ctx.Request.Context(), 15*time.Second)
	defer cancel()
	pod, ret := k8s.GetPod(timeout, generator.Name)
	pods := []k8s.PodDiagnostics{}
	if !ret.OK {
		if ret.Msg != i18n.K8S.NotFound {
			resp.JSON(ctx, ret)
			return
		}
	} else {
		if !k8s.GeneratorOwnsPod(generator, pod) {
			resp.JSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
			return
		}
		pods = append(pods, k8s.DescribePod(pod))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"pods": pods}))
}
