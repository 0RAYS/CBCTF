package service

import (
	"errors"
	"testing"

	"CBCTF/internal/model"
	"CBCTF/internal/resp"

	"gorm.io/gorm"
)

func TestFailedOptionalCountsAreNotDatabaseZeroes(t *testing.T) {
	database := queryFailureDB(t)
	_ = database.Callback().Query().Before("gorm:query").Register("test:count_failure", func(tx *gorm.DB) { _ = tx.AddError(errors.New("count unavailable")) })
	group := resp.GetGroupResp(BuildGroupView(database, model.Group{ID: 1, Name: "still-visible"}))
	if group["name"] != "still-visible" || group["users"] != nil || group["unavailable"] == nil {
		t.Fatalf("count failure turned into zero or lost object: %+v", group)
	}
	team := resp.GetTeamResp(BuildTeamView(database, model.Team{ID: 1, Name: "team"}), false)
	if team["users"] != nil || team["unavailable"] == nil {
		t.Fatalf("team count failure hidden: %+v", team)
	}
}
