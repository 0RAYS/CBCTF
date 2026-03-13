package router

import (
	"CBCTF/internal/cron"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"

	"github.com/gin-gonic/gin"
)

func GetCronJobs(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	cronJobs, count, ret := db.InitCronJobRepo(db.DB).List(form.Limit, form.Offset)
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
	cronJob := middleware.GetCronJob(ctx)
	if ret := db.InitCronJobRepo(db.DB).Update(cronJob.ID, db.UpdateCronJobOptions{
		Schedule: form.Schedule,
	}); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	newCronJob, ret := db.InitCronJobRepo(db.DB).GetByID(cronJob.ID)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	if ret = cron.ReloadCronJob(newCronJob.Name); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetCronJobResp(newCronJob)))
}
