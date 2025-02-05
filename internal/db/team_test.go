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

func InitTeamTest() {
	config.Env = &config.Config{}
	config.Env.Gorm.Log.Level = "debug"
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
	contest1, ok, msg := CreateContest(tx, constants.CreateContestForm{Name: "contest1", Size: 4, Start: time.Now(), Duration: time.Duration(10), Hidden: false})
	log.Logger.Debug(contest1.ID, ok, msg)
	contest2, ok, msg := CreateContest(tx, constants.CreateContestForm{Name: "contest2", Size: 4, Start: time.Now(), Duration: time.Duration(10), Hidden: false})
	log.Logger.Debug(contest2.ID, ok, msg)
	team1, ok, msg := CreateTeam(tx, constants.CreateTeamForm{Name: "team1", Captcha: contest1.Captcha}, user1, contest1)
	log.Logger.Debug(team1.ID, ok, msg)
	tx.Commit()
}

func TestCreateTeam(t *testing.T) {
	InitTeamTest()
	defer os.Remove("test.db")
	defer Close()
	var ctx context.Context
	user1, ok, msg := GetUserByID(DB, 1)
	user2, ok, msg := GetUserByID(DB, 2)
	contest1, ok, msg := GetContestByID(DB, 1)
	contest2, ok, msg := GetContestByID(DB, 2)
	tx := DB.WithContext(ctx).Begin()
	test, ok, msg := CreateTeam(tx, constants.CreateTeamForm{Name: "team1"}, user1, contest1)
	if ok {
		t.Fatal("Should not create duplicated team")
	}
	log.Logger.Debug(test, msg)
	test, ok, msg = CreateTeam(tx, constants.CreateTeamForm{Name: "team2"}, user1, contest1)
	if ok {
		t.Fatal("Team member should not be repeated")
	}
	log.Logger.Debug(test, msg)
	test, ok, msg = CreateTeam(tx, constants.CreateTeamForm{Name: "team2"}, user2, contest1)
	if !ok {
		t.Fatal("Should create team successfully", msg)
	}
	log.Logger.Debug(test, msg)
	test, ok, msg = CreateTeam(tx, constants.CreateTeamForm{Name: "team2"}, user1, contest2)
	if !ok {
		t.Fatal("Should create team successfully")
	}
	log.Logger.Debug(test, msg)
	tx.Commit()
}

func TestGetTeamByID(t *testing.T) {
	InitTeamTest()
	defer os.Remove("test.db")
	defer Close()
	_, ok, msg := GetTeamByID(DB, 1)
	if !ok {
		t.Fatal("Should get team successfully")
	}
	log.Logger.Debug(msg)
	_, ok, msg = GetTeamByID(DB, 0)
	if ok {
		t.Fatal("Should not get team successfully")
	}
	log.Logger.Debug(msg)

	// 不预加载
	team1, ok, msg := GetTeamByID(DB, 1, false)
	if !ok {
		t.Fatal("Should get team successfully")
	}
	if len(team1.Users) != 0 {
		t.Fatal("Should not preload users")
	}
	log.Logger.Debug(team1, ok, msg)

	// 预加载但不递归
	team1, ok, msg = GetTeamByID(DB, 1, true, false)
	if !ok {
		t.Fatal("Should get team successfully")
	}
	if len(team1.Users) == 0 {
		t.Fatal("Failed to preload users")
	}
	log.Logger.Debug(team1.Users)
	if len(team1.Users[0].Teams) != 0 {
		t.Fatal("Should not preload users.teams")
	}
	log.Logger.Debug(team1.Users[0].Teams, ok, msg)

	// 递归预加载
	team1, ok, msg = GetTeamByID(DB, 1, true, true)
	if !ok {
		t.Fatal("Should get team successfully")
	}
	if len(team1.Users) == 0 {
		t.Fatal("Failed to preload users")
	}
	log.Logger.Debug(team1.Users)
	if len(team1.Users[0].Teams) == 0 {
		t.Fatal("Failed to preload users.teams")
	}
	log.Logger.Debug(team1.Users[0].Teams)
}

func TestDeleteTeam(t *testing.T) {
	InitTeamTest()
	defer os.Remove("test.db")
	defer Close()
	var ctx context.Context
	tx := DB.WithContext(ctx).Begin()
	ok, msg := DeleteTeam(tx, model.Team{ID: 0})
	if ok {
		tx.Commit()
		t.Fatal("Should not delete team successfully")
	}
	log.Logger.Debug(msg)
	tx.Rollback()

	team1, ok, msg := GetTeamByID(DB, 1)
	if !ok {
		t.Fatal("Should get team successfully")
	}
	contest1, ok, msg := GetContestByID(DB, 1)
	if !ok {
		t.Fatal("Should get contest successfully")
	}
	var tmp []model.Team
	if err := DB.WithContext(ctx).Model(&contest1).Association("Teams").Find(&tmp); err != nil {
		t.Fatal(err)
	}
	if len(tmp) == 0 {
		t.Fatal("Failed to find association between contest and team")
	}
	log.Logger.Debug(tmp)

	var tmp2 []model.User
	if err := DB.WithContext(ctx).Model(&team1).Association("Users").Find(&tmp2); err != nil {
		t.Fatal(err)
	}
	if len(tmp2) == 0 {
		t.Fatal("Failed to find association between contest and user")
	}
	log.Logger.Debug(tmp2)

	if err := DB.WithContext(ctx).Model(&contest1).Association("Users").Find(&tmp2); err != nil {
		t.Fatal(err)
	}
	if len(tmp2) == 0 {
		t.Fatal("Failed to find association between user and team")
	}
	log.Logger.Debug(tmp)
	team1, _, _ = GetTeamByID(DB, 1)
	tx = DB.WithContext(ctx).Begin()
	ok, msg = DeleteTeam(tx, team1)
	if !ok {
		tx.Rollback()
		t.Fatal("Should delete team successfully")
	}
	tx.Commit()
	log.Logger.Debug(msg)

	user1, ok, msg := GetUserByID(DB, 1)
	if !ok {
		t.Fatal("Failed to get user by id")
	}

	if err := DB.WithContext(ctx).Model(&contest1).Association("Teams").Find(&tmp); err != nil {
		t.Fatal(err)
	}
	if len(tmp) != 0 {
		t.Fatal("Should not find association between contest and team")
	}
	log.Logger.Debug(tmp)

	if err := DB.WithContext(ctx).Model(&user1).Association("Teams").Find(&tmp); err != nil {
		t.Fatal(err)
	}
	if len(tmp) != 0 {
		t.Fatal("Should not find association between contest and user")
	}
	log.Logger.Debug(tmp)

	if err := DB.WithContext(ctx).Model(&contest1).Association("Users").Find(&tmp2); err != nil {
		t.Fatal(err)
	}
	if len(tmp2) != 0 {
		t.Fatal("Should not find association between user and team")
	}
	log.Logger.Debug(tmp2)
}

func TestJoinTeam(t *testing.T) {
	InitTeamTest()
	defer os.Remove("test.db")
	defer Close()
	var ctx context.Context
	user1, ok, msg := GetUserByID(DB, 1)
	user2, ok, msg := GetUserByID(DB, 2)
	team1, ok, msg := GetTeamByID(DB, 1)
	contest1, ok, msg := GetContestByID(DB, 1)
	tx := DB.WithContext(ctx).Begin()
	ok, msg = JoinTeam(tx, user1, team1, contest1)
	if ok {
		tx.Commit()
		t.Fatal("Should not join team successfully", msg)
	}
	log.Logger.Debug(msg)
	ok, msg = JoinTeam(tx, user2, team1, contest1)
	if !ok {
		tx.Rollback()
		t.Fatal("Should join team successfully", msg)
	}
	tx.Commit()
	log.Logger.Debug(msg)
}

func TestLeaveTeam(t *testing.T) {
	InitTeamTest()
	defer os.Remove("test.db")
	defer Close()
	var ctx context.Context
	user1, ok, msg := GetUserByID(DB, 1)
	user2, ok, msg := GetUserByID(DB, 2)
	team1, ok, msg := GetTeamByID(DB, 1)
	contest1, ok, msg := GetContestByID(DB, 1)
	tx := DB.WithContext(ctx).Begin()
	ok, msg = JoinTeam(tx, user2, team1, contest1)
	if !ok {
		t.Fatal("Should join team successfully")
	}
	tx.Commit()
	tx = DB.WithContext(ctx).Begin()
	ok, msg = LeaveTeam(tx, user1, team1, contest1)
	if ok {
		t.Fatal("Should not leave team successfully")
	}
	log.Logger.Debug(msg)
	ok, msg = LeaveTeam(tx, user2, team1, contest1)
	if !ok {
		t.Fatal("Should leave team successfully")
	}
	log.Logger.Debug(msg)
	ok, msg = LeaveTeam(tx, user1, team1, contest1)
	if ok {
		t.Fatal("Should not leave team successfully")
	}
	tx.Commit()
	log.Logger.Debug(msg)
}
