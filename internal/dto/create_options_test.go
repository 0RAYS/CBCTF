package dto

import (
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin/binding"
)

func TestCreateFormsAcceptEnabledCheckbox(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		data, _ := json.Marshal(map[string]bool{"on": enabled})
		var smtp CreateSmtpForm
		var webhook CreateWebhookForm
		if err := json.Unmarshal(data, &smtp); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(data, &webhook); err != nil {
			t.Fatal(err)
		}
		if smtp.On != enabled || webhook.On != enabled {
			t.Fatalf("create checkbox was lost: smtp=%v webhook=%v", smtp.On, webhook.On)
		}
	}
}

func TestVictimRatioSupportsAllTeams(t *testing.T) {
	for _, ratio := range []float64{-0.1, 0, 0.01, 0.5, 1, 1.01} {
		form := StartVictimsForm{Challenges: []string{"550e8400-e29b-41d4-a716-446655440000"}, TeamRatio: ratio}
		valid := binding.Validator.ValidateStruct(form) == nil
		if valid != (ratio > 0 && ratio <= 1) {
			t.Fatalf("unexpected validation for ratio %v: %v", ratio, valid)
		}
	}
}
