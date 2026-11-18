package resp

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
	"CBCTF/internal/view"
)

func TestGetChallengeUsageResp(t *testing.T) {
	usage := view.ChallengeUsageView{
		ContestID:   7,
		ContestName: "Example Contest",
		ContestChallenge: model.ContestChallenge{
			BaseModel: model.BaseModel{ID: 3},
			RandID:    "cc-rand-id",
			Name:      "Copied Challenge",
			Type:      model.PodsChallengeType,
			Category:  "web",
			Hidden:    true,
		},
		Solvers: 12,
	}
	want := gin.H{
		"contest_id":   uint(7),
		"contest_name": "Example Contest",
		"id":           "cc-rand-id",
		"name":         "Copied Challenge",
		"type":         model.PodsChallengeType,
		"category":     "web",
		"hidden":       true,
		"solvers":      int64(12),
	}
	if got := GetChallengeUsageResp(usage); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}
