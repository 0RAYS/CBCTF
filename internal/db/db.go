package db

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
)

var DB *gorm.DB

func Init() {
	var err error
	var level log.Level
	switch strings.ToLower(config.Env.GetString("gorm.log.level")) {
	case "info":
		level = log.Info
	case "warn":
		level = log.Warn
	case "error":
		level = log.Error
	case "silent":
		level = log.Silent
	default:
		level = log.Silent
	}
	DB, err = gorm.Open(sqlite.Open(config.Env.GetString("gorm.file")), &gorm.Config{Logger: log.NewGormLogger(level)})
	if err != nil {
		log.Logger.Panicf("failed to connect database: %v", err)
	}
	err = DB.AutoMigrate(&model.User{}, &model.Team{}, &model.Contest{}, &model.File{})
	if err != nil {
		log.Logger.Panicf("failed to migrate database: %v", err)
	}
	initAdmin()
}
