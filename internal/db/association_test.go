package db

import (
	"CBCTF/internal/config"
	"CBCTF/internal/constants"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"context"
	"os"
	"testing"
	"time"
)

func InitAssociationTest() {
	config.Env = &config.Config{}
	config.Env.Gorm.Log.Level = "info"
	config.Env.Log.Level = "debug"
	config.Env.Log.Save = false
	log.Init()
	InitTest()
	redis.Init()
	var ctx context.Context
	tx := DB.WithContext(ctx).Begin()
	user1, ok, msg := CreateUser(tx, constants.CreateUserForm{Name: "user1", Password: "password", Email: "user1@0rays.club"})
	log.Logger.Debug(user1.ID, ok, msg)
	user2, ok, msg := CreateUser(tx, constants.CreateUserForm{Name: "user2", Password: "password", Email: "user2@0rays.club"})
	log.Logger.Debug(user2.ID, ok, msg)

	contest1, ok, msg := CreateContest(tx, constants.CreateContestForm{Name: "contest1", Size: 1, Start: time.Now(), Duration: time.Duration(10)})
	log.Logger.Debug(contest1.ID, ok, msg)
	contest2, ok, msg := CreateContest(tx, constants.CreateContestForm{Name: "contest2", Size: 1, Start: time.Now(), Duration: time.Duration(10)})
	log.Logger.Debug(contest2.ID, ok, msg)
	team1, ok, msg := CreateTeam(tx, constants.CreateTeamForm{Name: "team1"}, user1, contest1)
	log.Logger.Debug(team1.ID, ok, msg)
	tx.Commit()
}

func TestAppendUserToTeam(t *testing.T) {
	InitAssociationTest()
	defer os.Remove("test.db")
	defer Close()
	var ctx context.Context
	user2, ok, msg := GetUserByID(DB, 2)
	log.Logger.Debug(user2, ok, msg)
	team1, ok, msg := GetTeamByID(DB, 1)
	log.Logger.Debug(team1, ok, msg)
	tx := DB.WithContext(ctx).Begin()
	if err := AppendUserToTeam(tx, user2, team1); err != nil {
		tx.Rollback()
		log.Logger.Error(err)
	}
	tx.Commit()
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
	defer os.Remove("test.db")
	defer Close()
	var ctx context.Context
	user2, ok, msg := GetUserByID(DB, 2)
	log.Logger.Debug(user2, ok, msg)
	contest1, ok, msg := GetContestByID(DB, 1)
	log.Logger.Debug(contest1, ok, msg)
	tx := DB.WithContext(ctx).Begin()
	if err := AppendUserToContest(tx, user2, contest1); err != nil {
		tx.Rollback()
		log.Logger.Error(err)
	}
	tx.Commit()
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
	defer os.Remove("test.db")
	defer Close()
	var ctx context.Context
	team1, ok, msg := GetTeamByID(DB, 1)
	log.Logger.Debug(team1, ok, msg)
	contest2, ok, msg := GetContestByID(DB, 2)
	log.Logger.Debug(contest2, ok, msg)
	tx := DB.WithContext(ctx).Begin()
	if err := AppendTeamToContest(tx, team1, contest2); err != nil {
		tx.Rollback()
		log.Logger.Error(err)
	}
	tx.Commit()
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
	defer os.Remove("test.db")
	defer Close()
	var ctx context.Context
	user1, ok, msg := GetUserByID(DB, 1)
	log.Logger.Debug(user1, ok, msg)
	team1, ok, msg := GetTeamByID(DB, 1)
	log.Logger.Debug(team1, ok, msg)
	tx := DB.WithContext(ctx).Begin()
	if err := DeleteUserFromTeam(tx, user1, team1); err != nil {
		tx.Rollback()
		log.Logger.Error(err)
	}
	tx.Commit()
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
	defer os.Remove("test.db")
	defer Close()
	var ctx context.Context
	user1, ok, msg := GetUserByID(DB, 1)
	log.Logger.Debug(user1, ok, msg)
	contest1, ok, msg := GetContestByID(DB, 1)
	log.Logger.Debug(contest1, ok, msg)
	tx := DB.WithContext(ctx).Begin()
	if err := DeleteUserFromContest(tx, user1, contest1); err != nil {
		tx.Rollback()
		log.Logger.Error(err)
	}
	tx.Commit()
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
	defer os.Remove("test.db")
	defer Close()
	var ctx context.Context
	team1, ok, msg := GetTeamByID(DB, 1)
	log.Logger.Debug(team1, ok, msg)
	contest1, ok, msg := GetContestByID(DB, 1)
	log.Logger.Debug(contest1, ok, msg)
	tx := DB.WithContext(ctx).Begin()
	if err := DeleteTeamFromContest(tx, team1, contest1); err != nil {
		tx.Rollback()
		log.Logger.Error(err)
	}
	tx.Commit()
	var tmp []model.Team
	if err := DB.WithContext(ctx).Model(&contest1).Association("Teams").Find(&tmp); err != nil {
		log.Logger.Error(err)
	}
	log.Logger.Debug(tmp)
	if len(tmp) != 0 {
		t.Fatalf("Failed to delete team from contest")
	}
}
