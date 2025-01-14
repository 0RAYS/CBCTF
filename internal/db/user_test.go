package db

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"context"
	"github.com/spf13/viper"
	"testing"
)

func InitUserTest() {
	config.Env = viper.New()
	config.Env.Set("gorm.file", ":memory:")
	config.Env.Set("gorm.log.level", "silent")
	config.Env.Set("log.level", "info")
	config.Env.Set("log.file", false)
	log.Init()
	Init()
	var ctx context.Context
	_, _, _ = CreateAdmin(ctx, "admin1", "password", "admin1@0rays.club")
	_, _, _ = CreateUser(ctx, "user1", "password", "user1@0rays.club")
}

func TestCreateUser(t *testing.T) {
	InitUserTest()
	var ctx context.Context
	if _, ok, _ := CreateUser(ctx, "test", "password", "test_email"); ok {
		t.Fatalf("Should not create user with invalid email")
	}
	if _, ok, _ := CreateUser(ctx, "user1", "password", "test@0rays.club"); ok {
		t.Fatalf("Should not create duplicated user")
	}
	if _, ok, _ := CreateUser(ctx, "test", "password", "user1@0rays.club"); ok {
		t.Fatalf("Should not create duplicated email")
	}
	if _, ok, _ := CreateAdmin(ctx, "user1", "password", "test@0rays.club"); !ok {
		t.Fatalf("Failed to create admin which name is duplicated with user")
	}
	if _, ok, _ := CreateAdmin(ctx, "test", "password", "user1@0rays.club"); !ok {
		t.Fatalf("Failed to create admin which email is duplicated with user")
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

// 2025-01-14 不完全测试
func TestDeleteUser(t *testing.T) {
	InitUserTest()
	var ctx context.Context
	if ok, _ := DeleteUser(ctx, 0); ok {
		t.Fatalf("Should return false when delete invalid user")
	}
	if ok, _ := DeleteUser(ctx, 1); !ok {
		t.Fatalf("Failed to delete user")
	}
}
