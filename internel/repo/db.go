package repo

import (
	"CBCTF/internel/config"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

var DB *gorm.DB

func Init() {
	var err error
	var level log.Level
	switch strings.ToUpper(config.Env.Gorm.Log.Level) {
	case "INFO":
		level = log.Info
	case "WARNING":
		level = log.Warn
	case "ERROR":
		level = log.Error
	case "SILENT":
		level = log.Silent
	default:
		level = log.Silent
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Env.Gorm.MySQL.User,
		config.Env.Gorm.MySQL.Pwd,
		config.Env.Gorm.MySQL.Host,
		config.Env.Gorm.MySQL.Port,
		config.Env.Gorm.MySQL.DB,
	)
	log.Logger.Infof("Connecting to MySQL database: %s:%d", config.Env.Gorm.MySQL.Host, config.Env.Gorm.MySQL.Port)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: log.NewGormLogger(level)})
	if err != nil {
		log.Logger.Fatalf("failed to connect database: %v", err)
	}
	if sql, err := DB.DB(); err != nil {
		log.Logger.Fatalf("failed to get database: %v", err)
	} else {
		sql.SetMaxIdleConns(config.Env.Gorm.MySQL.MaxIdleConns)
		sql.SetMaxOpenConns(config.Env.Gorm.MySQL.MaxOpenConns)
		sql.SetConnMaxLifetime(30 * time.Second)
	}

	// 指定数据表的存储引擎, 需要支持回滚操作
	err = DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&model.Admin{}, &model.Answer{}, &model.Challenge{}, &model.Cheat{}, &model.Contest{},
		&model.Device{}, &model.File{}, &model.Flag{}, &model.Notice{}, &model.Request{}, &model.Submission{},
		&model.Team{}, &model.Usage{}, &model.User{},
	)
	if err != nil {
		log.Logger.Fatalf("failed to migrate database: %v", err)
	}
	err = DB.SetupJoinTable(&model.User{}, "Teams", &model.UserTeam{})
	if err != nil {
		log.Logger.Fatalf("failed to setup join table: %v", err)
	}
	err = DB.SetupJoinTable(&model.User{}, "Contests", &model.UserContest{})
	if err != nil {
		log.Logger.Fatalf("failed to setup join table: %v", err)
	}

	log.Logger.Info("Connected to database")
	tx := DB.Begin()
	if ok, msg := initAdmin(tx); !ok {
		tx.Rollback()
		log.Logger.Fatalf("Failed to init Admin: %s", msg)
	}
	tx.Commit()
}

func initAdmin(tx *gorm.DB) (bool, string) {
	repo := InitAdminRepo(tx)
	if count, _, _ := repo.Count(); count == 0 {
		pwd := utils.UUID()
		_, ok, msg := repo.Create(CreateAdminOptions{
			Name:     "admin",
			Password: utils.HashPassword(pwd),
			Email:    "admin@0rays.club",
		})
		if !ok {
			return ok, msg
		}
		log.Logger.Infof("Init Admin: Admin{ name: admin, password: %s, email: admin@0rays.club}", pwd)
	}
	return true, "Success"
}

func Close() {
	if DB != nil {
		db, err := DB.DB()
		if err != nil {
			log.Logger.Errorf("failed to get database: %v", err)
		} else {
			_ = db.Close()
		}
	}
	log.Logger.Info("Database connection closed")
}
