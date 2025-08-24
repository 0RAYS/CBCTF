package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/middleware"
	"CBCTF/internal/resp"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetWebhookHistory(ctx *gin.Context) {
	var form f.GetModelsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	options := db.GetOptions{
		Preloads: map[string]db.GetOptions{"Webhook": {}, "Event": {}},
	}
	webhook := middleware.GetWebhook(ctx)
	if webhook.ID > 0 {
		options.Conditions = map[string]any{"webhook_id": webhook.ID}
	}
	histories, count, ok, msg := db.InitWebhookHistoryRepo(db.DB.WithContext(ctx)).List(form.Limit, form.Offset, options)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, history := range histories {
		data = append(data, resp.GetWebhookHistoryResp(history))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"histories": data, "count": count}})
}
