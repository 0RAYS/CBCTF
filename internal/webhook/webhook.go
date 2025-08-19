package webhook

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"bytes"
	"encoding/json"
	"net/http"
	"slices"
	"sync"
	"time"
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

func SelectWebhook(event model.Event) []model.Webhook {
	targets := make([]model.Webhook, 0)
	lock.RLock()
	for _, webhook := range Webhooks {
		if len(webhook.Events) == 0 || slices.Contains(webhook.Events, event.Type) {
			targets = append(targets, webhook)
		}
	}
	lock.RUnlock()
	return targets
}

func SendPayload(event model.Event, target model.Webhook) error {
	payload := Payload{
		IsAdmin: event.IsAdmin,
		Type:    event.Type,
		IP:      event.IP,
		Models:  event.Models,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		log.Logger.Warningf("Failed to marshal payload: %v", err)
		return err
	}
	start := time.Now()
	req, err := http.NewRequest(target.Method, target.URL, bytes.NewBuffer(data))
	if err != nil {
		log.Logger.Warningf("Failed to create request: %v", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "CBCTF-WEBHOOK")
	for k, v := range target.Headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: time.Duration(target.Timeout) * time.Second}
	resp, err := client.Do(req)
	options := db.CreateWebhookHistoryOptions{
		WebhookID: target.ID,
		EventID:   event.ID,
		RespCode:  0,
		Duration:  time.Since(start),
		Success:   false,
		Error:     "",
	}
	if err != nil {
		log.Logger.Warningf("Failed to send request: %v", err)
		options.Success = false
		options.Error = err.Error()
	} else {
		options.RespCode = resp.StatusCode
		options.Success = resp.StatusCode >= 200 && resp.StatusCode < 300
		if !options.Success {
			options.Error = resp.Status
		}
	}
	db.InitWebhookHistoryRepo(db.DB).Create(options)
	return nil
}
