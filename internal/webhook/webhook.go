package webhook

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"slices"
	"sync"
)

type Payload struct {
	IsAdmin bool          `json:"is_admin"`
	Type    string        `json:"type"`
	IP      string        `json:"ip"`
	Models  model.UintMap `json:"models"`
}

var (
	Webhooks []model.Webhook
	lock     sync.RWMutex
)

func Init() {
	Webhooks, _, _, _ = db.InitWebhookRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"on": true},
	})
}

func AddWebhook(webhook model.Webhook) {
	lock.Lock()
	Webhooks = append(Webhooks, webhook)
	lock.Unlock()
}

func DelWebhook(webhook model.Webhook) {
	lock.Lock()
	slices.DeleteFunc(Webhooks, func(w model.Webhook) bool {
		return w.ID == webhook.ID
	})
	lock.Unlock()
}
