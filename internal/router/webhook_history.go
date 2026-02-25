package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"

	"github.com/gin-gonic/gin"
)

func GetWebhookHistory(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	options := db.GetOptions{
		Preloads: map[string]db.GetOptions{"Webhook": {}, "Event": {}},
	}
	webhook := middleware.GetWebhook(ctx)
	if webhook.ID > 0 {
		options.Conditions = map[string]any{"webhook_id": webhook.ID}
	}
	histories, count, ret := db.InitWebhookHistoryRepo(db.DB).List(form.Limit, form.Offset, options)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, history := range histories {
		data = append(data, resp.GetWebhookHistoryResp(history))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"histories": data, "count": count}))
}
