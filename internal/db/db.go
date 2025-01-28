package db

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
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
	switch strings.ToLower(config.Env.Gorm.Type) {
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(config.Env.Gorm.SQLite.File), &gorm.Config{Logger: log.NewGormLogger(level)})
		log.Logger.Infof("Connecting to SQLite database: %s", config.Env.Gorm.SQLite.File)
	case "mysql":
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
	default:
		log.Logger.Fatalf("Unsupported database type: %s", config.Env.Gorm.Type)
	}
	if err != nil {
		log.Logger.Fatalf("failed to connect database: %v", err)
	}
	err = DB.AutoMigrate(
		&model.Admin{}, &model.User{}, &model.Team{},
		&model.Contest{}, &model.Avatar{}, &model.IP{},
		&model.Challenge{}, &model.Usage{},
	)
	if err != nil {
		log.Logger.Fatalf("failed to migrate database: %v", err)
	}
	log.Logger.Info("Connected to database")
	InitAdmin()
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
