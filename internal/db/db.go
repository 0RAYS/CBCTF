package db

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
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
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=30 TimeZone=Asia/Shanghai",
		config.Env.Gorm.Postgres.Host,
		config.Env.Gorm.Postgres.Port,
		config.Env.Gorm.Postgres.User,
		config.Env.Gorm.Postgres.Pwd,
		config.Env.Gorm.Postgres.DB,
		getPostgresSSLMode(config.Env.Gorm.Postgres.SSLMode),
	)
	log.Logger.Infof("Connecting to PostgreSQL database: %s:%d", config.Env.Gorm.Postgres.Host, config.Env.Gorm.Postgres.Port)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: log.NewGormLogger(level),
	})
	if err != nil {
		log.Logger.Fatalf("Failed to connect database: %s", err)
	}
	if sql, err := DB.DB(); err != nil {
		log.Logger.Fatalf("Failed to get database: %s", err)
	} else {
		sql.SetMaxIdleConns(config.Env.Gorm.Postgres.MaxIdleConns)
		sql.SetMaxOpenConns(config.Env.Gorm.Postgres.MaxOpenConns)
		sql.SetConnMaxIdleTime(time.Hour)
		sql.SetConnMaxLifetime(24 * time.Hour)
	}

	if err = DB.Use(prometheus.New(prometheus.Config{
		DBName:          config.Env.Gorm.Postgres.DB,
		RefreshInterval: 15,
		StartServer:     false,
	})); err != nil {
		log.Logger.Warningf("Failed to register prometheus: %s", err)
	}
	if err = DB.Exec(`CREATE EXTENSION IF NOT EXISTS pg_trgm`).Error; err != nil {
		log.Logger.Warningf("Failed to ensure pg_trgm extension: %s", err)
	}

	err = DB.AutoMigrate(
		&model.Challenge{}, &model.ChallengeFlag{}, &model.Cheat{}, &model.Victim{}, &model.Pod{},
		&model.ContestChallenge{}, &model.ContestFlag{}, &model.CronJob{}, &model.Device{}, &model.Email{},
		&model.Event{}, &model.File{}, &model.Generator{}, &model.Group{}, &model.Notice{}, &model.Oauth{},
		&model.Permission{}, &model.Request{}, &model.Role{}, &model.Setting{}, &model.Smtp{}, &model.Submission{},
		&model.Team{}, &model.TeamFlag{}, &model.Traffic{}, &model.User{}, &model.Webhook{}, &model.WebhookHistory{},
	)
	if err != nil {
		log.Logger.Fatalf("Failed to migrate database: %s", err)
	}
	err = DB.SetupJoinTable(&model.User{}, "Teams", &model.UserTeam{})
	if err != nil {
		log.Logger.Fatalf("Failed to setup join table: %s", err)
	}
	err = DB.SetupJoinTable(&model.User{}, "Contests", &model.UserContest{})
	if err != nil {
		log.Logger.Fatalf("Failed to setup join table: %s", err)
	}
	err = DB.SetupJoinTable(&model.User{}, "Groups", &model.UserGroup{})
	if err != nil {
		log.Logger.Fatalf("Failed to setup join table: %s", err)
	}
	err = DB.SetupJoinTable(&model.Role{}, "Permissions", &model.RolePermission{})
	if err != nil {
		log.Logger.Fatalf("Failed to setup join table: %s", err)
	}
	log.Logger.Info("Connected to database")

	if ret := InitSettingRepo(DB).InitSettings(); !ret.OK {
		log.Logger.Fatalf("Failed to init settings: %s %v", ret.Msg, ret.Attr)
	}
	if ret := InitPermissionRepo(DB).InitPermissions(); !ret.OK {
		log.Logger.Fatalf("Failed to init permissions: %s %v", ret.Msg, ret.Attr)
	}
	if ret := InitRoleRepo(DB).InitDefaultRoles(); !ret.OK {
		log.Logger.Fatalf("Failed to init default roles: %s %v", ret.Msg, ret.Attr)
	}
	if ret := InitGroupRepo(DB).InitDefaultGroups(); !ret.OK {
		log.Logger.Fatalf("Failed to init default groups: %s %v", ret.Msg, ret.Attr)
	}
	if ret := InitCronJobRepo(DB).InitCronJob(); !ret.OK {
		log.Logger.Fatalf("Failed to init cron jobs: %s %v", ret.Msg, ret.Attr)
	}
	if ret := InitUserRepo(DB).InitAdmin(); !ret.OK {
		log.Logger.Fatalf("Failed to init Admin: %v", ret)
	}
	InitOauthRepo(DB).RegisterDefault()
}

func getPostgresSSLMode(enabled bool) string {
	if enabled {
		return "require"
	}
	return "disable"
}

func Stop() {
	db, err := DB.DB()
	if err != nil {
		log.Logger.Warningf("Failed to stop PostgreSQL connection: %s", err)
		return
	}
	if err = db.Close(); err != nil {
		log.Logger.Warningf("Failed to stop PostgreSQL connection: %s", err)
	}
}
