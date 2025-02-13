package db

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
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
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: log.NewGormLogger(level)})
	log.Logger.Infof("Connecting to MySQL database: %s:%d", config.Env.Gorm.MySQL.Host, config.Env.Gorm.MySQL.Port)
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

	err = DB.AutoMigrate(
		&model.Admin{}, &model.User{}, &model.Team{},
		&model.Contest{}, &model.Avatar{}, &model.IP{},
		&model.Challenge{}, &model.Usage{}, &model.Flag{},
		&model.Docker{}, &model.Submission{}, &model.Device{},
		&model.Traffic{},
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
	InitAdmin(tx)
	tx.Commit()
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
