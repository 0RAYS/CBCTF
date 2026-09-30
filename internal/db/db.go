package db

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"CBCTF/internal/config"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

var (
	// DB and HTTPDB share the request pool.
	DB     *gorm.DB
	HTTPDB *gorm.DB
	TaskDB *gorm.DB
	// Advisory locks use a separate pool so callbacks can still acquire query connections.
	WorkloadLockDB *gorm.DB
	CronDB         *gorm.DB
)

type Tx = gorm.DB

func WithTransactionDB(root *gorm.DB, fn func(tx *Tx) model.RetVal) model.RetVal {
	if root == nil {
		return model.RetVal{Msg: i18n.DB.Unavailable}
	}
	var ret model.RetVal
	err := root.Transaction(func(tx *gorm.DB) error {
		ret = fn(tx)
		if !ret.OK {
			return errors.New(ret.Msg)
		}
		return nil
	})
	if err != nil && (ret.OK || ret.Msg == "") {
		log.Logger.Warningf("Failed to execute transaction: %s", err)
		return model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	return ret
}

func Init() {
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

	DB = openPostgresPool("http", config.Env.Gorm.Postgres.MaxOpenConns, config.Env.Gorm.Postgres.MaxIdleConns, level)
	HTTPDB = DB
	TaskDB = openPostgresPool(
		"task",
		backgroundPoolLimit(config.Env.Gorm.Postgres.MaxOpenConns, 4),
		backgroundPoolLimit(config.Env.Gorm.Postgres.MaxIdleConns, 4),
		level,
	)
	WorkloadLockDB = openPostgresPool(
		"workload-lock",
		backgroundPoolLimit(config.Env.Gorm.Postgres.MaxOpenConns, 4),
		backgroundPoolLimit(config.Env.Gorm.Postgres.MaxIdleConns, 4),
		level,
	)
	CronDB = openPostgresPool(
		"cron",
		backgroundPoolLimit(config.Env.Gorm.Postgres.MaxOpenConns, 10),
		backgroundPoolLimit(config.Env.Gorm.Postgres.MaxIdleConns, 10),
		level,
	)

	if err := DB.Exec(`CREATE EXTENSION IF NOT EXISTS pg_trgm`).Error; err != nil {
		log.Logger.Warningf("Failed to ensure pg_trgm extension: %s", err)
	}

	err := DB.AutoMigrate(
		&model.Branding{}, &model.Challenge{}, &model.ChallengeFlag{}, &model.Cheat{}, &model.Victim{}, &model.Pod{},
		&model.Contest{}, &model.ContestChallenge{}, &model.ContestFlag{}, &model.CronJob{}, &model.Email{},
		&model.Event{}, &model.File{}, &model.Generator{}, &model.Group{}, &model.Notice{}, &model.Oauth{},
		&model.Permission{}, &model.Request{}, &model.Role{}, &model.Setting{}, &model.Smtp{}, &model.Submission{},
		&model.Task{}, &model.Team{}, &model.TeamFlag{}, &model.Traffic{}, &model.User{}, &model.Webhook{}, &model.WebhookHistory{},
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
	if ret := InitBrandingRepo(DB).InitDefault(); !ret.OK {
		log.Logger.Fatalf("Failed to init branding: %s %v", ret.Msg, ret.Attr)
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
	if ret := InitOauthRepo(DB).RegisterDefault(); !ret.OK {
		log.Logger.Fatalf("Failed to init OAuth providers: %s %v", ret.Msg, ret.Attr)
	}
}

func openPostgresPool(name string, maxOpenConns, maxIdleConns int, level log.Level) *gorm.DB {
	sslMode := "disable"
	if config.Env.Gorm.Postgres.SSLMode {
		sslMode = "require"
	}
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s connect_timeout=30 TimeZone=Asia/Shanghai",
		config.Env.Gorm.Postgres.Host,
		config.Env.Gorm.Postgres.Port,
		config.Env.Gorm.Postgres.User,
		config.Env.Gorm.Postgres.Pwd,
		config.Env.Gorm.Postgres.DB,
		sslMode,
	)
	log.Logger.Infof(
		"Connecting to PostgreSQL database pool %q: %s:%d max_open=%d max_idle=%d",
		name, config.Env.Gorm.Postgres.Host, config.Env.Gorm.Postgres.Port, maxOpenConns, maxIdleConns,
	)
	pool, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: log.NewGormLogger(level),
	})
	if err != nil {
		log.Logger.Fatalf("Failed to connect database pool %q: %s", name, err)
	}
	if err = pool.Use(JSONUpdates{}); err != nil {
		log.Logger.Fatalf("Failed to configure JSON serializers for pool %q: %s", name, err)
	}
	sqlDB, err := pool.DB()
	if err != nil {
		log.Logger.Fatalf("Failed to get database pool %q: %s", name, err)
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxIdleTime(time.Hour)
	sqlDB.SetConnMaxLifetime(24 * time.Hour)
	return pool
}

func backgroundPoolLimit(limit, divisor int) int {
	if limit <= 0 {
		return limit
	}
	return max(1, limit/divisor)
}

func Stop() {
	for name, pool := range map[string]*gorm.DB{
		"workload-lock": WorkloadLockDB,
		"cron":          CronDB,
		"task":          TaskDB,
		"http":          HTTPDB,
	} {
		if pool == nil {
			continue
		}
		sqlDB, err := pool.DB()
		if err != nil {
			log.Logger.Warningf("Failed to stop PostgreSQL pool %q: %s", name, err)
			continue
		}
		if err = sqlDB.Close(); err != nil {
			log.Logger.Warningf("Failed to stop PostgreSQL pool %q: %s", name, err)
		}
	}
}
