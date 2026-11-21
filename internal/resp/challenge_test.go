package resp

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
	"CBCTF/internal/view"
)

func TestGetChallengeResp(t *testing.T) {
	challengeView := view.ChallengeView{
		Challenge: model.Challenge{
			BaseModel: model.BaseModel{ID: 3},
			RandID:    "rand-id",
			Name:      "Copied Challenge",
			Type:      model.PodsChallengeType,
			Category:  "web",
			Hidden:    true,
		},
		ContestName: "Example Contest",
		Flags:       []view.ChallengeFlagView{{ID: 1, Value: "flag{value}"}},
	}
	want := gin.H{
		"id":               "rand-id",
		"name":             "Copied Challenge",
		"description":      "",
		"category":         "web",
		"type":             model.PodsChallengeType,
		"contest_id":       uint(0),
		"contest_name":     "Example Contest",
		"hidden":           true,
		"attempt":          int64(0),
		"hints":            []string(nil),
		"tags":             []string(nil),
		"generator_image":  "",
		"flags":            []gin.H{{"id": uint(1), "value": "flag{value}"}},
		"docker_compose":   "",
		"network_policies": model.NetworkPolicies(nil),
		"file":             "",
	}
	if got := GetChallengeResp(challengeView); !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}
