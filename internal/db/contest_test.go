package db

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"github.com/spf13/viper"
	"testing"
)

func InitContestTest() {
	config.Env = viper.New()
	config.Env.Set("gorm.file", ":memory:")
	config.Env.Set("gorm.log.level", "silent")
	config.Env.Set("log.level", "debug")
	config.Env.Set("log.file", false)
	log.Init()
	Init()
	var ctx context.Context

	user1, ok, msg := CreateUser(ctx, "user1", "password", "user1@0rays.club")
	log.Logger.Debug(user1.ID, ok, msg)
	contest1, ok, msg := CreateContest(ctx, "contest1")
	log.Logger.Debug(contest1.ID, ok, msg)
	team1, ok, msg := CreateTeam(ctx, "team1", user1.ID, contest1.ID)
	log.Logger.Debug(team1.ID, ok, msg)
}

func TestCreateContest(t *testing.T) {
	InitContestTest()
	var ctx context.Context
	test, ok, msg := CreateContest(ctx, "contest1")
	if ok {
		t.Fatal("Should not create duplicated admin")
	}
	log.Logger.Debug(test, msg)
}

func TestGetContestByID(t *testing.T) {
	InitContestTest()
	var ctx context.Context
	test, ok, msg := GetContestByID(ctx, 0)
	if ok {
		t.Fatal("Should not get contest with invalid id")
	}
	log.Logger.Debug(test, msg)
	// 测试递归预加载
	contest1, ok, msg := GetContestByID(ctx, 1, true, true)
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

	// 测试不预加载
	contest1, ok, msg = GetContestByID(ctx, 1, false)
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

	// 测试预加载但不递归
	contest1, ok, msg = GetContestByID(ctx, 1, true, false)
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
	var ctx context.Context
	if ok, _ := DeleteContest(ctx, 0); ok {
		t.Fatal("Should return true when delete invalid contest")
	}
	user1, ok, _ := GetUserByID(ctx, 1)
	if !ok {
		t.Fatal("Failed to get user by id")
	}
	contest1, ok, _ := GetContestByID(ctx, 1)
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

	if ok, msg := DeleteContest(ctx, 1); !ok {
		t.Fatalf("Failed to delete contest by id: %s", msg)
	}
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
