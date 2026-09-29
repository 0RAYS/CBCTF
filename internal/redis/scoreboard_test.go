package redis

import (
	"context"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"CBCTF/internal/model"
)

type rankingSnapshotHook struct {
	payload  []byte
	commands []string
}

func (*rankingSnapshotHook) DialHook(next goredis.DialHook) goredis.DialHook { return next }
func (*rankingSnapshotHook) ProcessPipelineHook(next goredis.ProcessPipelineHook) goredis.ProcessPipelineHook {
	return next
}
func (h *rankingSnapshotHook) ProcessHook(_ goredis.ProcessHook) goredis.ProcessHook {
	return func(_ context.Context, cmd goredis.Cmder) error {
		h.commands = append(h.commands, cmd.Name())
		if cmd.Name() == "set" {
			h.payload = cmd.Args()[2].([]byte)
			cmd.(*goredis.StatusCmd).SetVal("OK")
		} else {
			cmd.(*goredis.StringCmd).SetVal(string(h.payload))
		}
		return nil
	}
}

func TestRankingSnapshotPreservesDatabaseOrderAndEmptyLists(t *testing.T) {
	old := RDB
	client := goredis.NewClient(&goredis.Options{})
	RDB = client
	t.Cleanup(func() { RDB = old; _ = client.Close() })
	hook := &rankingSnapshotHook{}
	client.AddHook(hook)
	teams := []model.Team{{ID: 2, Score: 100000, Last: time.Unix(1, 1)}, {ID: 1, Score: 100000, Last: time.Unix(1, 2)}}
	if ret := UpdateTeamRanking(context.Background(), 1, teams); !ret.OK {
		t.Fatal(ret)
	}
	got, ret := GetTeamRanking(context.Background(), 1, 0, 1)
	if !ret.OK || len(got) != 2 || got[0].ID != 2 || !got[0].Last.Equal(teams[0].Last) {
		t.Fatalf("ordering/precision lost: %+v %+v", got, ret)
	}
	if len(hook.commands) != 2 || hook.commands[0] != "set" || hook.commands[1] != "get" {
		t.Fatalf("non-atomic cache protocol: %v", hook.commands)
	}
	if ret := UpdateTeamRanking(context.Background(), 1, []model.Team{}); !ret.OK {
		t.Fatal(ret)
	}
	got, ret = GetTeamRanking(context.Background(), 1, 0, 10)
	if !ret.OK || len(got) != 0 {
		t.Fatalf("empty snapshot not preserved: %+v %+v", got, ret)
	}
}
