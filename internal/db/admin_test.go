package db

import (
	"CBCTF/internal/config"
	"CBCTF/internal/constants"
	"CBCTF/internal/log"
	"CBCTF/internal/redis"
	"context"
	"os"
	"testing"
)

func InitAdminTest() {
	config.Env = &config.Config{}
	config.Env.Gorm.Log.Level = "debug"
	config.Env.Log.Level = "debug"
	config.Env.Log.Save = false
	log.Init()
	InitTest()
	redis.Init()
	var ctx context.Context
	tx := DB.WithContext(ctx).Begin()
	admin1, ok, msg := CreateAdmin(tx, "admin1", "password", "admin1@0rays.club")
	log.Logger.Info(admin1.ID, ok, msg)
	user1, ok, msg := CreateUser(tx, constants.CreateUserForm{Name: "user1", Password: "password", Email: "user1@0rays.club"})
	log.Logger.Info(user1.ID, ok, msg)
	tx.Commit()
}

func TestCreateAdmin(t *testing.T) {
	InitAdminTest()
	defer os.Remove("test.db")
	defer Close()
	var ctx context.Context
	tx := DB.WithContext(ctx).Begin()
	if _, ok, _ := CreateAdmin(tx, "test", "password", "test_email"); ok {
		t.Fatal("Should not create admin with invalid email")
	}
	if _, ok, _ := CreateAdmin(tx, "admin1", "password", "test@0rays.club"); ok {
		t.Fatal("Should not create duplicated admin")
	}
	if _, ok, _ := CreateAdmin(tx, "test", "password", "admin1@0rays.club"); ok {
		t.Fatal("Should not create duplicated email")
	}
	if _, ok, _ := CreateUser(tx, constants.CreateUserForm{Name: "admin1", Password: "password", Email: "test@0rays.club"}); !ok {
		t.Fatal("Failed to create user which name is duplicated with admin")
	}
	if _, ok, _ := CreateUser(tx, constants.CreateUserForm{Name: "test", Password: "password", Email: "admin1@0rays.club"}); ok {
		t.Fatal("Should not create user which email is duplicated with admin")
	}
	if admin1, _, _ := GetAdminByID(DB, 1); admin1.Password == "password" {
		t.Fatal("Failed to hash password")
	}
	tx.Commit()
}

func TestGetAdminByID(t *testing.T) {
	InitAdminTest()
	defer os.Remove("test.db")
	defer Close()
	if _, ok, _ := GetAdminByID(DB, 0); ok {
		t.Fatal("Should not get admin with invalid id")
	}
	if _, ok, _ := GetAdminByID(DB, 1); !ok {
		t.Fatal("Failed to get admin by id")
	}
}

func TestDeleteAdmin(t *testing.T) {
	InitAdminTest()
	defer os.Remove("test.db")
	defer Close()
	var ctx context.Context
	tx := DB.WithContext(ctx).Begin()
	if ok, _ := DeleteAdmin(tx, 0); !ok {
		t.Fatal("Should return true when delete invalid admin")
	}
	if ok, msg := DeleteAdmin(tx, 1); !ok {
		t.Fatalf("Failed to delete admin by id: %s", msg)
	}
	tx.Commit()
}
