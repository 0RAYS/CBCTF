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

func InitUserTest() {
	config.Env = &config.Config{}
	config.Env.Gorm.Log.Level = "debug"
	config.Env.Log.Level = "debug"
	config.Env.Log.Save = false
	log.Init()
	InitTest()
	redis.Init()
	var ctx context.Context
	tx := DB.WithContext(ctx).Begin()
	_, _, _ = CreateAdmin(tx, "admin1", "password", "admin1@0rays.club")
	user1, _, _ := CreateUser(tx, constants.CreateUserForm{Name: "user1", Password: "password", Email: "user1@0rays.club"})
	contest1, _, _ := CreateContest(tx, constants.CreateContestForm{Name: "contest1", Size: 1, Start: time.Now(), Duration: time.Duration(10), Hidden: false})
	_, _, _ = CreateTeam(tx, constants.CreateTeamForm{Name: "team1"}, user1, contest1)
	tx.Commit()
}

func TestCreateUser(t *testing.T) {
	InitUserTest()
	defer os.Remove("test.db")
	defer Close()
	var ctx context.Context
	tx := DB.WithContext(ctx).Begin()
	if _, ok, _ := CreateUser(tx, constants.CreateUserForm{Name: "test", Password: "password", Email: "test_email"}); ok {
		t.Fatalf("Should not create user with invalid email")
	}
	if _, ok, _ := CreateUser(tx, constants.CreateUserForm{Name: "user1", Password: "password", Email: "test@0rays.club"}); ok {
		t.Fatalf("Should not create duplicated user")
	}
	if _, ok, _ := CreateUser(tx, constants.CreateUserForm{Name: "test", Password: "password", Email: "user1@0rays.club"}); ok {
		t.Fatalf("Should not create duplicated email")
	}
	if _, ok, _ := CreateAdmin(tx, "user1", "password", "test@0rays.club"); !ok {
		t.Fatalf("Failed to create admin which name is duplicated with user")
	}
	if _, ok, _ := CreateAdmin(tx, "test", "password", "user1@0rays.club"); ok {
		t.Fatalf("Should not create admin which email is duplicated with user")
	}
	if user1, _, _ := GetUserByID(ctx, 1); user1.Password == "password" {
		t.Fatalf("Failed to hash password")
	}
	tx.Commit()
}

func TestGetUserByID(t *testing.T) {
	InitUserTest()
	defer os.Remove("test.db")
	defer Close()
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
	defer os.Remove("test.db")
	defer Close()
	var ctx context.Context
	tx := DB.WithContext(ctx).Begin()
	if ok, _ := DeleteUser(tx, ctx, 0); ok {
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

	if ok, _ := DeleteUser(tx, ctx, 1); !ok {
		t.Fatalf("Failed to delete user")
	}
	tx.Commit()

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
	defer os.Remove("test.db")
	defer Close()
	var ctx context.Context
	tx := DB.WithContext(ctx).Begin()
	test, _, _ := CreateUser(tx, constants.CreateUserForm{Name: "test", Password: "password", Email: "test@0rays.club"})
	_, _ = UpdateUser(tx, test.ID, map[string]interface{}{"hidden": true})
	tx.Commit()
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
