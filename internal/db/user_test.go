package db

import (
	"CBCTF/internal/config"
	"CBCTF/internal/constants"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"context"
	"testing"
	"time"
)

func InitUserTest() {
	config.Env = &config.Config{}
	config.Env.Gorm.Type = "sqlite"
	config.Env.Gorm.SQLite.File = ":memory:"
	config.Env.Gorm.Log.Level = "debug"
	config.Env.Log.Level = "debug"
	config.Env.Log.Save = false
	log.Init()
	Init()
	redis.Init()
	var ctx context.Context
	_, _, _ = CreateAdmin(ctx, "admin1", "password", "admin1@0rays.club")
	_, _, _ = CreateUser(ctx, constants.CreateUserForm{Name: "user1", Password: "password", Email: "user1@0rays.club"})
	_, _, _ = CreateContest(ctx, constants.CreateContestForm{Name: "contest1", Size: 1, Start: time.Now(), Duration: time.Duration(10)})
	_, _, _ = CreateTeam(ctx, constants.CreateTeamForm{Name: "team1"}, 1, 1)
}

func TestCreateUser(t *testing.T) {
	InitUserTest()
	var ctx context.Context
	if _, ok, _ := CreateUser(ctx, constants.CreateUserForm{Name: "test", Password: "password", Email: "test_email"}); ok {
		t.Fatalf("Should not create user with invalid email")
	}
	if _, ok, _ := CreateUser(ctx, constants.CreateUserForm{Name: "user1", Password: "password", Email: "test@0rays.club"}); ok {
		t.Fatalf("Should not create duplicated user")
	}
	if _, ok, _ := CreateUser(ctx, constants.CreateUserForm{Name: "test", Password: "password", Email: "user1@0rays.club"}); ok {
		t.Fatalf("Should not create duplicated email")
	}
	if _, ok, _ := CreateAdmin(ctx, "user1", "password", "test@0rays.club"); !ok {
		t.Fatalf("Failed to create admin which name is duplicated with user")
	}
	if _, ok, _ := CreateAdmin(ctx, "test", "password", "user1@0rays.club"); ok {
		t.Fatalf("Should not create admin which email is duplicated with user")
	}
	if user1, _, _ := GetUserByID(ctx, 1); user1.Password == "password" {
		t.Fatalf("Failed to hash password")
	}
}

func TestGetUserByID(t *testing.T) {
	InitUserTest()
	var ctx context.Context
	if _, ok, _ := GetUserByID(ctx, 0); ok {
		t.Fatalf("Should not get user with invalid id")
	}
	if _, ok, _ := GetUserByID(ctx, 1); !ok {
		t.Fatalf("Failed to get user by id")
	}
}

func TestDeleteUser(t *testing.T) {
	InitUserTest()
	var ctx context.Context
	if ok, _ := DeleteUser(ctx, 0); ok {
		t.Fatalf("Sest hould return false when delete invalid user")
	}

	user1, ok, msg := GetUserByID(ctx, 1)
	if !ok {
		t.Fatalf(msg)
	}
	contest1, ok, msg := GetContestByID(ctx, 1)
	if !ok {
		t.Fatalf(msg)
	}

	var tmp []model.Team
	if err := DB.WithContext(ctx).Model(&user1).Association("Teams").Find(&tmp); err != nil {
		t.Fatalf(err.Error())
	}
	if len(tmp) == 0 {
		t.Fatalf("Failed to find association between user and team")
	}
	log.Logger.Debug(tmp)
	var tmp2 []model.Contest
	if err := DB.WithContext(ctx).Model(&user1).Association("Contests").Find(&tmp2); err != nil {
		t.Fatalf(err.Error())
	}
	if len(tmp2) == 0 {
		t.Fatalf("Failed to find association between user and contest")
	}
	log.Logger.Debug(tmp2)

	if err := DB.WithContext(ctx).Model(&contest1).Association("Teams").Find(&tmp); err != nil {
		t.Fatalf(err.Error())
	}
	if len(tmp) == 0 {
		t.Fatalf("Failed to find association between contest and team")
	}
	log.Logger.Debug(tmp)

	if ok, _ := DeleteUser(ctx, 1); !ok {
		t.Fatalf("Failed to delete user")
	}

	if err := DB.WithContext(ctx).Model(&user1).Association("Teams").Find(&tmp); err != nil {
		t.Fatalf(err.Error())
	}
	if len(tmp) != 0 {
		t.Fatalf("Should not find association between user and team")
	}
	log.Logger.Debug(tmp)
	if err := DB.WithContext(ctx).Model(&user1).Association("Contests").Find(&tmp2); err != nil {
		t.Fatalf(err.Error())
	}
	if len(tmp2) != 0 {
		t.Fatalf("Should not find association between user and contest")
	}
	log.Logger.Debug(tmp2)

	if err := DB.WithContext(ctx).Model(&contest1).Association("Teams").Find(&tmp); err != nil {
		t.Fatalf(err.Error())
	}
	if len(tmp) != 0 {
		t.Fatalf("Should not find association between contest and team")
	}
	log.Logger.Debug(tmp)
}

func TestGetUsers(t *testing.T) {
	InitUserTest()
	test, _, _ := CreateUser(context.Background(), constants.CreateUserForm{Name: "test", Password: "password", Email: "test@0rays.club"})
	_, _ = UpdateUser(context.Background(), test.ID, map[string]interface{}{"hidden": true})
	var ctx context.Context
	users, count, ok, msg := GetUsers(ctx, 0, 0, true)
	log.Logger.Info(users, count, ok, msg)
	if len(users) != 2 {
		t.Fatalf("Failed to get all users")
	}
	users, count, ok, msg = GetUsers(ctx, 0, 0, false)
	log.Logger.Info(users, count, ok, msg)
	if len(users) != 1 {
		t.Fatalf("Failed to filter hidden users")
	}
}
