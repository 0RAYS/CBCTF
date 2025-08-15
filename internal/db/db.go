package db

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
)

var DB *gorm.DB

func Init() {
	var err error
	var level log.Level
	switch config.Env.Gorm.Log.Level {
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
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: log.NewGormLogger(level),
	})
	if err != nil {
		log.Logger.Fatalf("Failed to connect database: %s", err.Error())
	}
	if sql, err := DB.DB(); err != nil {
		log.Logger.Fatalf("Failed to get database: %s", err.Error())
	} else {
		sql.SetMaxIdleConns(config.Env.Gorm.MySQL.MaxIdleConns)
		sql.SetMaxOpenConns(config.Env.Gorm.MySQL.MaxOpenConns)
		sql.SetConnMaxIdleTime(time.Hour)
		sql.SetConnMaxLifetime(24 * time.Hour)
	}

	if DB.Use(prometheus.New(prometheus.Config{
		DBName:          config.Env.Gorm.MySQL.DB,
		RefreshInterval: 15,
		StartServer:     false,
	})) != nil {
		log.Logger.Warningf("Failed to register prometheus: %s", err.Error())
	}

	// 指定数据表的存储引擎, 需要支持回滚操作
	err = DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&model.Admin{}, &model.Challenge{}, &model.ChallengeFlag{}, &model.Contest{}, &model.ContestChallenge{},
		&model.ContestFlag{}, &model.Device{}, &model.Docker{}, &model.Event{}, &model.File{},
		&model.Notice{}, &model.Oauth{}, &model.Request{}, &model.Submission{}, &model.Team{}, &model.TeamFlag{},
		&model.User{}, &model.Victim{}, &model.Pod{}, &model.Container{}, &model.Cheat{}, &model.Traffic{},
	)
	if err != nil {
		log.Logger.Fatalf("Failed to migrate database: %s", err.Error())
	}
	err = DB.SetupJoinTable(&model.User{}, "Teams", &model.UserTeam{})
	if err != nil {
		log.Logger.Fatalf("Failed to setup join table: %s", err.Error())
	}
	err = DB.SetupJoinTable(&model.User{}, "Contests", &model.UserContest{})
	if err != nil {
		log.Logger.Fatalf("Failed to setup join table: %s", err.Error())
	}
	log.Logger.Info("Connected to database")

	if ok, msg := InitAdminRepo(DB).InitAdmin(); !ok {
		log.Logger.Fatalf("Failed to init Admin: %s", msg)
	}
	InitOauthRepo(DB).RegisterDefault()
}
