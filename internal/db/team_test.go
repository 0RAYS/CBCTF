package db

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"github.com/spf13/viper"
	"testing"
)

func InitTeamTest() {
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
	user2, ok, msg := CreateUser(ctx, "user2", "password", "user2@0rays.club")
	log.Logger.Debug(user2.ID, ok, msg)
	contest1, ok, msg := CreateContest(ctx, "contest1")
	log.Logger.Debug(contest1.ID, ok, msg)
	contest2, ok, msg := CreateContest(ctx, "contest2")
	log.Logger.Debug(contest2.ID, ok, msg)
	team1, ok, msg := CreateTeam(ctx, "team1", user1.ID, contest1.ID)
	log.Logger.Debug(team1.ID, ok, msg)
}

func TestCreateTeam(t *testing.T) {
	InitTeamTest()
	var ctx context.Context
	test, ok, msg := CreateTeam(ctx, "team1", 1, 1)
	if ok {
		t.Fatal("Should not create duplicated team")
	}
	log.Logger.Debug(test, msg)
	test, ok, msg = CreateTeam(ctx, "team2", 1, 1)
	if ok {
		t.Fatal("Team member should not be repeated")
	}
	log.Logger.Debug(test, msg)
	test, ok, msg = CreateTeam(ctx, "team2", 2, 1)
	if !ok {
		t.Fatal("Should create team successfully")
	}
	log.Logger.Debug(test, msg)
	test, ok, msg = CreateTeam(ctx, "team2", 1, 2)
	if !ok {
		t.Fatal("Should create team successfully")
	}
	log.Logger.Debug(test, msg)
}

func TestGetTeamByID(t *testing.T) {
	InitTeamTest()
	var ctx context.Context
	_, ok, msg := GetTeamByID(ctx, 1)
	if !ok {
		t.Fatal("Should get team successfully")
	}
	log.Logger.Debug(msg)
	_, ok, msg = GetTeamByID(ctx, 0)
	if ok {
		t.Fatal("Should not get team successfully")
	}
	log.Logger.Debug(msg)

	// 不预加载
	team1, ok, msg := GetTeamByID(ctx, 1, false)
	if !ok {
		t.Fatal("Should get team successfully")
	}
	if len(team1.Users) != 0 {
		t.Fatal("Should not preload users")
	}
	log.Logger.Debug(team1, ok, msg)

	// 预加载但不递归
	team1, ok, msg = GetTeamByID(ctx, 1, true, false)
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
	team1, ok, msg = GetTeamByID(ctx, 1, true, true)
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
	var ctx context.Context
	ok, msg := DeleteTeam(ctx, 0)
	if ok {
		t.Fatal("Should not delete team successfully")
	}
	log.Logger.Debug(msg)

	team1, ok, msg := GetTeamByID(ctx, 1)
	if !ok {
		t.Fatal("Should get team successfully")
	}
	contest1, ok, msg := GetContestByID(ctx, 1)
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

	ok, msg = DeleteTeam(ctx, 1)
	if !ok {
		t.Fatal("Should delete team successfully")
	}
	log.Logger.Debug(msg)

	user1, ok, msg := GetUserByID(ctx, 1)
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
	var ctx context.Context
	ok, msg := JoinTeam(ctx, 1, 1, 1)
	if ok {
		t.Fatal("Should not join team successfully")
	}
	log.Logger.Debug(msg)
	ok, msg = JoinTeam(ctx, 2, 1, 1)
	if !ok {
		t.Fatal("Should join team successfully")
	}
	log.Logger.Debug(msg)
}

func TestLeaveTeam(t *testing.T) {
	InitTeamTest()
	var ctx context.Context
	ok, msg := JoinTeam(ctx, 2, 1, 1)
	if !ok {
		t.Fatal("Should join team successfully")
	}
	ok, msg = LeaveTeam(ctx, 1, 1, 1)
	if ok {
		t.Fatal("Should not leave team successfully")
	}
	log.Logger.Debug(msg)
	ok, msg = LeaveTeam(ctx, 2, 1, 1)
	if !ok {
		t.Fatal("Should leave team successfully")
	}
	log.Logger.Debug(msg)
	ok, msg = LeaveTeam(ctx, 1, 1, 1)
	if !ok {
		t.Fatal("Should leave team successfully")
	}
	log.Logger.Debug(msg)
	ok, msg = LeaveTeam(ctx, 1, 1, 1)
	if ok {
		t.Fatal("Should not leave team successfully")
	}
	log.Logger.Debug(msg)
}
