package resp

import (
	"CBCTF/internal/model"
	"CBCTF/internal/view"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestGetContestRespCollections(t *testing.T) {
	start := time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		name    string
		contest model.Contest
		want    map[string]string
	}{
		{
			name: "nil",
			want: map[string]string{"rules": "[]", "prizes": "[]", "timelines": "[]"},
		},
		{
			name:    "empty",
			contest: model.Contest{Rules: model.StringList{}, Prizes: model.Prizes{}, Timelines: model.Timelines{}},
			want:    map[string]string{"rules": "[]", "prizes": "[]", "timelines": "[]"},
		},
		{
			name: "populated",
			contest: model.Contest{
				Rules:     model.StringList{"First rule", "Second rule"},
				Prizes:    model.Prizes{{Amount: "$0", Description: "Recognition"}},
				Timelines: model.Timelines{{Date: start, Title: "Start", Description: "Challenges open"}},
			},
			want: map[string]string{
				"rules":     `["First rule","Second rule"]`,
				"prizes":    `[{"amount":"$0","description":"Recognition"}]`,
				"timelines": `[{"date":"2026-09-10T00:00:00Z","title":"Start","description":"Challenges open"}]`,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			input := view.ContestView{Contest: tc.contest}
			before, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}
			encoded, err := json.Marshal(GetContestResp(input, false))
			if err != nil {
				t.Fatal(err)
			}
			var fields map[string]json.RawMessage
			if err := json.Unmarshal(encoded, &fields); err != nil {
				t.Fatal(err)
			}
			for field, want := range tc.want {
				if string(fields[field]) != want {
					t.Errorf("%s = %s, want %s", field, fields[field], want)
				}
			}
			after, err := json.Marshal(input)
			if err != nil {
				t.Fatal(err)
			}
			if string(before) != string(after) {
				t.Fatal("response builder mutated the original view")
			}
		})
	}
}

func TestGetContestRespMetadata(t *testing.T) {
	start := time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC)
	input := view.ContestView{
		Contest: model.Contest{
			BaseModel: model.BaseModel{ID: 7},
			Name:      "Example", Description: "Description", Start: start, Duration: time.Hour,
			Size: 3, Prefix: "CTF", Victims: 2, Picture: "/contest.png", Hidden: true, Blood: false,
			Captcha: "private-code",
		},
		TeamCount: 4, UserCount: 12, NoticeCount: 2, Highest: 1000, SolvedCount: 5,
	}
	for _, admin := range []bool{false, true} {
		for _, statsReady := range []bool{false, true} {
			input.StatsReady = statsReady
			want := gin.H{
				"id": uint(7), "name": "Example", "description": "Description", "start": start,
				"duration": int64(3600), "size": 3, "prefix": "CTF", "victims": int64(2),
				"picture": model.FileURL("/contest.png"), "hidden": true, "blood": false,
				"teams": int64(4), "users": int64(12), "notices": int64(2),
				"rules": model.StringList{}, "prizes": model.Prizes{}, "timelines": model.Timelines{},
			}
			if statsReady {
				want["highest"] = float64(1000)
				want["solved"] = int64(5)
			}
			if admin {
				want["captcha"] = "private-code"
			}
			if got := GetContestResp(input, admin); !reflect.DeepEqual(got, want) {
				t.Errorf("admin=%v statsReady=%v: got %#v, want %#v", admin, statsReady, got, want)
			}
		}
	}
}

func TestGetContestRespEmptyRulesStorageRoundTrip(t *testing.T) {
	stored, err := (model.StringList{}).Value()
	if err != nil {
		t.Fatal(err)
	}
	if stored != nil {
		t.Fatalf("empty StringList storage value = %v, want SQL NULL", stored)
	}
	rules := model.StringList{"previous rule"}
	if err := rules.Scan(stored); err != nil {
		t.Fatal(err)
	}
	if rules != nil {
		t.Fatalf("scanned SQL NULL = %#v, want nil slice", rules)
	}
	input := view.ContestView{Contest: model.Contest{Rules: rules}}
	encoded, err := json.Marshal(GetContestResp(input, false)["rules"])
	if err != nil {
		t.Fatal(err)
	}
	if string(encoded) != "[]" {
		t.Fatalf("response rules = %s, want []", encoded)
	}
	if input.Contest.Rules != nil {
		t.Fatal("response builder mutated the original nil rules")
	}
}
