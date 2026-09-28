package webhook

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

func TestHTTPFailuresAreRecordedAndReturnedForRetry(t *testing.T) {
	log.Init()
	previous := db.TaskDB
	t.Cleanup(func() { db.TaskDB = previous })
	var err error
	db.TaskDB, err = gorm.Open(postgres.Open("host=127.0.0.1 user=test dbname=test"), &gorm.Config{
		DryRun: true, DisableAutomaticPing: true, SkipDefaultTransaction: true,
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	pool, _ := db.TaskDB.DB()
	t.Cleanup(func() { _ = pool.Close() })
	var histories []model.WebhookHistory
	if err := db.TaskDB.Callback().Create().After("gorm:create").Register("test:history", func(tx *gorm.DB) {
		if history, ok := tx.Statement.Dest.(*model.WebhookHistory); ok {
			histories = append(histories, *history)
		}
	}); err != nil {
		t.Fatal(err)
	}
	for _, status := range []int{200, 204, 400, 503} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var payload Payload
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Type != "submit_flag" {
				t.Errorf("unexpected webhook payload: %+v, %v", payload, err)
			}
			w.WriteHeader(status)
		}))
		err := SendPayload(model.Event{Type: "submit_flag"}, model.Webhook{URL: server.URL, Method: "POST", Timeout: 1})
		server.Close()
		if (err != nil) != (status >= 300) {
			t.Fatalf("HTTP %d: unexpected retry result %v", status, err)
		}
		if len(histories) == 0 {
			t.Fatal("delivery attempt was not recorded")
		}
		history := histories[len(histories)-1]
		if history.RespCode != status || history.Success != (status < 300) {
			t.Fatalf("HTTP %d: unexpected history %+v", status, history)
		}
	}
}
