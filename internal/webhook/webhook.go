package webhook

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"slices"
	"sync"
	"time"
)

type Payload struct {
	Type   string        `json:"type"`
	IP     string        `json:"ip"`
	Models model.UintMap `json:"models"`
}

var (
	Webhooks []model.Webhook
	lock     sync.RWMutex
)

func Init() {
	Webhooks, _, _ = db.InitWebhookRepo(db.DB).List(-1, -1, db.GetOptions{
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
	Webhooks = slices.DeleteFunc(Webhooks, func(w model.Webhook) bool {
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

func logURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "invalid-url"
	}
	return parsed.Scheme + "://" + parsed.Host + parsed.Path
}

func SendPayload(event model.Event, target model.Webhook) error {
	payload := Payload{
		Type:   event.Type,
		IP:     event.IP,
		Models: event.Models,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	start := time.Now()
	req, err := http.NewRequest(target.Method, target.URL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "CBCTF-WEBHOOK")
	for k, v := range target.Headers {
		req.Header.Set(k, v)
	}
	timeout := time.Duration(target.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	duration := time.Since(start)
	options := db.CreateWebhookHistoryOptions{
		WebhookID: target.ID,
		EventID:   event.ID,
		RespCode:  0,
		Duration:  duration,
		Success:   false,
		Error:     "",
	}
	if err != nil {
		options.Success = false
		options.Error = err.Error()
		db.InitWebhookHistoryRepo(db.DB).Create(options)
		return err
	}
	defer resp.Body.Close()
	options.RespCode = resp.StatusCode
	options.Success = resp.StatusCode >= 200 && resp.StatusCode < 300
	if !options.Success {
		options.Error = resp.Status
		log.Logger.Warningf("Webhook returned non-success status: event_id=%d webhook_id=%d method=%s url=%s status=%s duration=%s", event.ID, target.ID, target.Method, logURL(target.URL), resp.Status, duration)
	} else {
		log.Logger.Debugf("Webhook sent: event_id=%d webhook_id=%d method=%s url=%s status=%d duration=%s", event.ID, target.ID, target.Method, logURL(target.URL), resp.StatusCode, duration)
	}
	db.InitWebhookHistoryRepo(db.DB).Create(options)
	return nil
}
