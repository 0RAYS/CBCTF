package router

import (
	"CBCTF/internal/cron"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"

	"github.com/gin-gonic/gin"
)

func GetCronJobs(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	cronJobs, count, ret := service.ListCronJobs(db.DB, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0, len(cronJobs))
	for _, cronJob := range cronJobs {
		data = append(data, resp.GetCronJobResp(cronJob))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"cronjobs": data, "count": count}))
}

func GetCronJob(ctx *gin.Context) {
	resp.JSON(ctx, model.SuccessRetVal(resp.GetCronJobResp(middleware.GetCronJob(ctx))))
}

func UpdateCronJob(ctx *gin.Context) {
	var form dto.UpdateCronJobForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateCronJobEventType)
	cronJob, ret := service.UpdateCronJob(db.DB, middleware.GetCronJob(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	if ret = cron.ReloadCronJob(cronJob.Name); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetCronJobResp(cronJob)))
}
