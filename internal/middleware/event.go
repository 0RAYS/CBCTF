package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/task"
	"CBCTF/internal/webhook"

	"github.com/gin-gonic/gin"
)

const (
	CTXEventTypeKey    = "EventType"
	CTXEventSuccessKey = "EventSuccess"
	CTXEventModelsKey  = "EventModels"
)

func Events(ctx *gin.Context) {
	ctx.Next()

	t := ctx.GetString(CTXEventTypeKey)
	if t == "" || t == model.SkipEventType {
		return
	}
	models := []model.Model{
		GetUser(ctx), GetContest(ctx), GetTeam(ctx), GetFile(ctx), GetNotice(ctx), GetChallenge(ctx), GetWebhook(ctx),
		GetContestChallenge(ctx), GetContestFlag(ctx), GetVictim(ctx), GetCheat(ctx), GetOauth(ctx), GetSmtp(ctx),
	}
	options := db.CreateEventOptions{
		Type:    t,
		Success: ctx.GetBool(CTXEventSuccessKey),
		IP:      ctx.ClientIP(),
		Magic:   GetMagic(ctx),
		Models:  make(model.UintMap),
	}
	for _, m := range models {
		if id := m.GetBaseModel().ID; id > 0 {
			options.Models[model.ModelName(m)] = id
		}
	}
	if value, ok := ctx.Get(CTXEventModelsKey); ok {
		if eventModels, ok := value.(model.UintMap); ok {
			for k, v := range eventModels {
				options.Models[k] = v
			}
		}
	}
	options.Models["Self"] = GetSelf(ctx).ID
	if event, ret := db.InitEventRepo(db.DB).Create(options); ret.OK {
		for _, target := range webhook.SelectWebhook(event) {
			if _, err := task.EnqueueWebhookTask(event, target); err != nil {
				log.Logger.Warningf("Failed to enqueue webhook task: %s", err)
			}
		}
	}
}
