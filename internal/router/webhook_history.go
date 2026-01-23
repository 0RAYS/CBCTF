package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetWebhookHistory(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
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
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, history := range histories {
		data = append(data, resp.GetWebhookHistoryResp(history))
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"histories": data, "count": count}))
}
