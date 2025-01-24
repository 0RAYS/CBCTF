package db

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"context"
	"testing"
	"time"
)

func InitAssociationTest() {
	config.Env.Gorm.SQLite.File = ":memory:"
	config.Env.Gorm.Log.Level = "debug"
	config.Env.Log.Level = "debug"
	config.Env.Log.Save = false
	log.Init()
	Init()
	redis.Init()
	var ctx context.Context
	user1, ok, msg := CreateUser(ctx, "user1", "password", "user1@0rays.club", "", "", false, false, false)
	log.Logger.Debug(user1.ID, ok, msg)
	user2, ok, msg := CreateUser(ctx, "user2", "password", "user2@0rays.club", "", "", false, false, false)
	log.Logger.Debug(user2.ID, ok, msg)

	contest1, ok, msg := CreateContest(ctx, "contest1", "test", "", 1, time.Now(), time.Duration(10), false)
	log.Logger.Debug(contest1.ID, ok, msg)
	contest2, ok, msg := CreateContest(ctx, "contest2", "test", "", 1, time.Now(), time.Duration(10), false)
	log.Logger.Debug(contest2.ID, ok, msg)

	team1, ok, msg := CreateTeam(ctx, "team1", user1.ID, contest1.ID)
	log.Logger.Debug(team1.ID, ok, msg)
}

func TestAppendUserToTeam(t *testing.T) {
	InitAssociationTest()
	var ctx context.Context
	user2, ok, msg := GetUserByID(ctx, 2)
	log.Logger.Debug(user2, ok, msg)
	team1, ok, msg := GetTeamByID(ctx, 1)
	log.Logger.Debug(team1, ok, msg)
	if err := AppendUserToTeam(ctx, user2, team1); err != nil {
		log.Logger.Error(err)
	}
	var tmp []model.User
	if err := DB.WithContext(ctx).Model(&team1).Association("Users").Find(&tmp); err != nil {
		log.Logger.Error(err)
	}
	log.Logger.Debug(tmp)
	if len(tmp) != 2 {
		t.Fatalf("Failed to append user to team")
	}
}

func TestAppendUserToContest(t *testing.T) {
	InitAssociationTest()
	var ctx context.Context
	user2, ok, msg := GetUserByID(ctx, 2)
	log.Logger.Debug(user2, ok, msg)
	contest1, ok, msg := GetContestByID(ctx, 1)
	log.Logger.Debug(contest1, ok, msg)
	if err := AppendUserToContest(ctx, user2, contest1); err != nil {
		log.Logger.Error(err)
	}
	var tmp []model.User
	if err := DB.WithContext(ctx).Model(&contest1).Association("Users").Find(&tmp); err != nil {
		log.Logger.Error(err)
	}
	log.Logger.Debug(tmp)
	if len(tmp) != 2 {
		t.Fatalf("Failed to append user to contest")
	}
}

func TestAppendTeamToContest(t *testing.T) {
	InitAssociationTest()
	var ctx context.Context
	team1, ok, msg := GetTeamByID(ctx, 1)
	log.Logger.Debug(team1, ok, msg)
	contest2, ok, msg := GetContestByID(ctx, 2)
	log.Logger.Debug(contest2, ok, msg)
	if err := AppendTeamToContest(ctx, team1, contest2); err != nil {
		log.Logger.Error(err)
	}
	var tmp []model.Team
	if err := DB.WithContext(ctx).Model(&contest2).Association("Teams").Find(&tmp); err != nil {
		log.Logger.Error(err)
	}
	log.Logger.Debug(tmp)
	if len(tmp) != 1 {
		t.Fatalf("Failed to append team to contest")
	}
}

func TestDeleteUserFromTeam(t *testing.T) {
	InitAssociationTest()
	var ctx context.Context
	user1, ok, msg := GetUserByID(ctx, 1)
	log.Logger.Debug(user1, ok, msg)
	team1, ok, msg := GetTeamByID(ctx, 1)
	log.Logger.Debug(team1, ok, msg)
	if err := DeleteUserFromTeam(ctx, user1, team1); err != nil {
		log.Logger.Error(err)
	}
	var tmp []model.User
	if err := DB.WithContext(ctx).Model(&team1).Association("Users").Find(&tmp); err != nil {
		log.Logger.Error(err)
	}
	log.Logger.Debug(tmp)
	if len(tmp) != 0 {
		t.Fatalf("Failed to delete user from team")
	}
}

func TestDeleteUserFromContest(t *testing.T) {
	InitAssociationTest()
	var ctx context.Context
	user1, ok, msg := GetUserByID(ctx, 1)
	log.Logger.Debug(user1, ok, msg)
	contest1, ok, msg := GetContestByID(ctx, 1)
	log.Logger.Debug(contest1, ok, msg)
	if err := DeleteUserFromContest(ctx, user1, contest1); err != nil {
		log.Logger.Error(err)
	}
	var tmp []model.User
	if err := DB.WithContext(ctx).Model(&contest1).Association("Users").Find(&tmp); err != nil {
		log.Logger.Error(err)
	}
	log.Logger.Debug(tmp)
	if len(tmp) != 0 {
		t.Fatalf("Failed to delete user from contest")
	}
}

func TestDeleteTeamFromContest(t *testing.T) {
	InitAssociationTest()
	var ctx context.Context
	team1, ok, msg := GetTeamByID(ctx, 1)
	log.Logger.Debug(team1, ok, msg)
	contest1, ok, msg := GetContestByID(ctx, 1)
	log.Logger.Debug(contest1, ok, msg)
	if err := DeleteTeamFromContest(ctx, team1, contest1); err != nil {
		log.Logger.Error(err)
	}
	var tmp []model.Team
	if err := DB.WithContext(ctx).Model(&contest1).Association("Teams").Find(&tmp); err != nil {
		log.Logger.Error(err)
	}
	log.Logger.Debug(tmp)
	if len(tmp) != 0 {
		t.Fatalf("Failed to delete team from contest")
	}
}
