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

func InitContestTest() {
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
	contest1, ok, msg := CreateContest(tx, constants.CreateContestForm{Name: "contest1", Size: 1, Start: time.Now(), Duration: time.Duration(10)})
	log.Logger.Debug(contest1.ID, ok, msg)
	team1, ok, msg := CreateTeam(tx, constants.CreateTeamForm{Name: "team1", Captcha: contest1.Captcha}, user1, contest1)
	log.Logger.Debug(team1.ID, ok, msg)
	tx.Commit()
}

func TestCreateContest(t *testing.T) {
	InitContestTest()
	defer os.Remove("test.db")
	defer Close()
	var ctx context.Context
	tx := DB.WithContext(ctx).Begin()
	test, ok, msg := CreateContest(tx, constants.CreateContestForm{Name: "contest1", Size: 1, Start: time.Now(), Duration: time.Duration(10)})
	if ok {
		tx.Commit()
		t.Fatal("Should not create duplicated admin")
	}
	tx.Rollback()
	log.Logger.Debug(test, msg)
}

func TestGetContestByID(t *testing.T) {
	InitContestTest()
	defer os.Remove("test.db")
	defer Close()
	test, ok, msg := GetContestByID(DB, 0)
	if ok {
		t.Fatal("Should not get contest with invalid id")
	}
	log.Logger.Debug(test, msg)
	// 递归预加载
	contest1, ok, msg := GetContestByID(DB, 1, true, true)
	if !ok {
		t.Fatal("Failed to get contest by id")
	}
	if len(contest1.Users) == 0 {
		t.Fatal("Failed to preload users")
	}
	log.Logger.Debug(contest1.Users[0])
	if len(contest1.Teams) == 0 {
		t.Fatal("Failed to preload teams")
	}
	log.Logger.Debug(contest1.Teams[0])
	if len(contest1.Teams[0].Users) == 0 {
		t.Fatal("Failed to preload teams.users")
	}
	log.Logger.Debug(contest1.Teams[0].Users[0])
	if len(contest1.Users[0].Contests) == 0 {
		t.Fatal("Failed to preload users.contests")
	}
	log.Logger.Debug(contest1.Users[0].Contests[0])
	if len(contest1.Users[0].Teams) == 0 {
		t.Fatal("Failed to preload users.teams")
	}
	log.Logger.Debug(contest1.Users[0].Teams[0])

	// 不预加载
	contest1, ok, msg = GetContestByID(DB, 1, false)
	if !ok {
		t.Fatal("Failed to get contest by id")
	}
	if len(contest1.Users) != 0 {
		t.Fatal("Should not preload users")
	}
	log.Logger.Debug(contest1.Users)
	if len(contest1.Teams) != 0 {
		t.Fatal("Should not preload teams")
	}
	log.Logger.Debug(contest1.Teams)

	// 预加载但不递归
	contest1, ok, msg = GetContestByID(DB, 1, true, false)
	if !ok {
		t.Fatal("Failed to get contest by id")
	}
	if len(contest1.Users) == 0 {
		t.Fatal("Failed to preload users")
	}
	log.Logger.Debug(contest1.Users[0])
	if len(contest1.Teams) == 0 {
		t.Fatal("Failed to preload teams")
	}
	log.Logger.Debug(contest1.Teams[0])
	if len(contest1.Teams[0].Users) != 0 {
		t.Fatal("Should not preload teams.users")
	}
	log.Logger.Debug(contest1.Teams[0].Users)
	if len(contest1.Users[0].Contests) != 0 {
		t.Fatal("Should not preload users.contests")
	}
	log.Logger.Debug(contest1.Users[0].Contests)
	if len(contest1.Users[0].Teams) != 0 {
		t.Fatal("Should not preload users.teams")
	}
	log.Logger.Debug(contest1.Users[0].Teams)
}

func TestDeleteContest(t *testing.T) {
	InitContestTest()
	defer os.Remove("test.db")
	defer Close()
	var ctx context.Context
	tx := DB.WithContext(ctx).Begin()
	if ok, _ := DeleteContest(tx, model.Contest{ID: 0}); ok {
		tx.Commit()
		t.Fatal("Should return true when delete invalid contest")
	}
	tx.Rollback()
	user1, ok, _ := GetUserByID(DB, 1)
	if !ok {
		t.Fatal("Failed to get user by id")
	}
	contest1, ok, _ := GetContestByID(DB, 1)
	if !ok {
		t.Fatal("Failed to get contest by id")
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
	if err := DB.WithContext(ctx).Model(&contest1).Association("Users").Find(&tmp2); err != nil {
		t.Fatal(err)
	}
	if len(tmp2) == 0 {
		t.Fatal("Failed to find association between contest and user")
	}
	log.Logger.Debug(tmp2)

	if err := DB.WithContext(ctx).Model(&user1).Association("Teams").Find(&tmp); err != nil {
		t.Fatal(err)
	}
	if len(tmp) == 0 {
		t.Fatal("Failed to find association between user and team")
	}
	log.Logger.Debug(tmp)

	tx = DB.WithContext(ctx).Begin()
	if ok, msg := DeleteContest(tx, contest1); !ok {
		t.Fatalf("Failed to delete contest by id: %s", msg)
	}
	tx.Commit()
	if err := DB.WithContext(ctx).Model(&contest1).Association("Teams").Find(&tmp); err != nil {
		t.Fatal(err)
	}
	if len(tmp) != 0 {
		t.Fatal("Should not find association between contest and team")
	}
	if err := DB.WithContext(ctx).Model(&contest1).Association("Users").Find(&tmp2); err != nil {
		t.Fatal(err)
	}
	log.Logger.Debug(tmp)
	if len(tmp2) != 0 {
		t.Fatal("Should not find association between contest and user")
	}
	if err := DB.WithContext(ctx).Model(&user1).Association("Teams").Find(&tmp); err != nil {
		t.Fatal(err)
	}
	log.Logger.Debug(tmp2)
	if len(tmp) != 0 {
		t.Fatal("Should not find association between user and team")
	}
	log.Logger.Debug(tmp)
}
