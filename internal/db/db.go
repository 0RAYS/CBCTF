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
	switch strings.ToUpper(config.Env.GetString("gorm.log.level")) {
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
	DB, err = gorm.Open(sqlite.Open(config.Env.GetString("gorm.file")), &gorm.Config{Logger: log.NewGormLogger(level)})
	if err != nil {
		log.Logger.Panicf("failed to connect database: %v", err)
	}
	err = DB.AutoMigrate(&model.Admin{}, &model.User{}, &model.Team{}, &model.Contest{}, &model.File{}, &model.IP{})
	if err != nil {
		log.Logger.Panicf("failed to migrate database: %v", err)
	}
	InitAdmin()
}
