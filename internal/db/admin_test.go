package db

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"context"
	"github.com/spf13/viper"
	"testing"
)

func InitAdminTest() {
	config.Env = viper.New()
	config.Env.Set("gorm.file", ":memory:")
	config.Env.Set("gorm.log.level", "silent")
	config.Env.Set("log.level", "info")
	config.Env.Set("log.file", false)
	config.Env.Set("upload.max", 8)
	log.Init()
	Init()
	var ctx context.Context
	admin1, ok, msg := CreateAdmin(ctx, "admin1", "password", "admin1@0rays.club")
	log.Logger.Info(admin1.ID, ok, msg)
	user1, ok, msg := CreateUser(ctx, "user1", "password", "user1@0rays.club")
	log.Logger.Info(user1.ID, ok, msg)
}

func TestCreateAdmin(t *testing.T) {
	InitAdminTest()
	var ctx context.Context
	if _, ok, _ := CreateAdmin(ctx, "test", "password", "test_email"); ok {
		t.Fatal("Should not create admin with invalid email")
	}
	if _, ok, _ := CreateAdmin(ctx, "admin1", "password", "test@0rays.club"); ok {
		t.Fatal("Should not create duplicated admin")
	}
	if _, ok, _ := CreateAdmin(ctx, "test", "password", "admin1@0rays.club"); ok {
		t.Fatal("Should not create duplicated email")
	}
	if _, ok, _ := CreateUser(ctx, "admin1", "password", "test@0rays.club"); !ok {
		t.Fatal("Failed to create user which name is duplicated with admin")
	}
	if _, ok, _ := CreateUser(ctx, "test", "password", "admin1@0rays.club"); ok {
		t.Fatal("Should not create user which email is duplicated with admin")
	}
	if admin1, _, _ := GetAdminByID(ctx, 1); admin1.Password == "password" {
		t.Fatal("Failed to hash password")
	}
}

func TestGetAdminByID(t *testing.T) {
	InitAdminTest()
	var ctx context.Context
	if _, ok, _ := GetAdminByID(ctx, 0); ok {
		t.Fatal("Should not get admin with invalid id")
	}
	if _, ok, _ := GetAdminByID(ctx, 1); !ok {
		t.Fatal("Failed to get admin by id")
	}
}

func TestDeleteAdmin(t *testing.T) {
	InitAdminTest()
	var ctx context.Context
	if ok, _ := DeleteAdmin(ctx, 0); !ok {
		t.Fatal("Should return true when delete invalid admin")
	}
	if ok, msg := DeleteAdmin(ctx, 1); !ok {
		t.Fatalf("Failed to delete admin by id: %s", msg)
	}
}
