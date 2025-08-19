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
		IsAdmin: IsAdmin(ctx),
		Type:    t,
		Success: ctx.GetBool(CTXEventSuccessKey),
		IP:      ctx.ClientIP(),
		Magic:   GetMagic(ctx),
		Models:  make(model.UintMap),
	}
	for _, m := range models {
		if id := m.GetBasicModel().ID; id > 0 {
			options.Models[m.GetModelName()] = id
		}
	}
	if value, ok := ctx.Get(CTXEventModelsKey); ok {
		if eventModels, ok := value.(model.UintMap); ok {
			for k, v := range eventModels {
				options.Models[k] = v
			}
		}
	}
	options.Models["Self"] = GetSelfID(ctx)
	if event, ok, _ := db.InitEventRepo(db.DB.WithContext(ctx)).Create(options); ok {
		for _, target := range webhook.SelectWebhook(event) {
			if _, err := task.EnqueueWebhookTask(event, target); err != nil {
				log.Logger.Warningf("Failed to enqueue webhook task: %s", err)
			}
		}
	}
}
