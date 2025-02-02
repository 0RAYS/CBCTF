package cmd

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 由 MySQL 迁移至 SQLite 时, 有严重bug, 不要使用
func migrateDB() {
	config.Init()
	log.Init()

	// 检查参数是否完整
	if *srcType == "" || *sqliteFile == "" || *mysqlDSN == "" {
		log.Logger.Errorf("Error: Both --src --sqlite and --mysql parameters are required")
	}
	var (
		src *gorm.DB
		dst *gorm.DB
		err error
	)
	switch *srcType {
	case "sqlite":
		src, err = gorm.Open(sqlite.Open(*sqliteFile), &gorm.Config{})
		if err != nil {
			log.Logger.Errorf("Failed to connect to SQLite: %v", err)
			return
		}
		dst, err = gorm.Open(mysql.Open(*mysqlDSN), &gorm.Config{})
		if err != nil {
			log.Logger.Errorf("Failed to connect to MySQL: %v", err)
			return
		}
	case "mysql":
		src, err = gorm.Open(mysql.Open(*mysqlDSN), &gorm.Config{})
		if err != nil {
			log.Logger.Errorf("Failed to connect to MySQL: %v", err)
			return
		}
		dst, err = gorm.Open(sqlite.Open(*sqliteFile), &gorm.Config{})
		if err != nil {
			log.Logger.Errorf("Failed to connect to SQLite: %v", err)
			return
		}
	default:
		log.Logger.Errorf("Error: Unsupported source database type: %s", *srcType)
		return
	}
	models := []interface{}{model.Admin{}, model.User{}, model.Team{},
		model.Contest{}, model.Avatar{}, model.IP{},
		model.Challenge{}, model.Usage{}, model.Flag{},
		model.Docker{}, model.Submission{},
	}
	err = dst.AutoMigrate(models...)
	if err != nil {
		log.Logger.Errorf("Failed to init data tables: %v", err)
		return
	}
	var (
		admins      []model.Admin
		users       []model.User
		teams       []model.Team
		contests    []model.Contest
		avatars     []model.Avatar
		ips         []model.IP
		challenges  []model.Challenge
		usages      []model.Usage
		flags       []model.Flag
		dockers     []model.Docker
		submissions []model.Submission
		errs        []error
	)
	errs = append(errs, src.Find(&admins).Error)
	errs = append(errs, src.Preload(clause.Associations).Find(&users).Error)
	errs = append(errs, src.Find(&teams).Error)
	errs = append(errs, src.Find(&contests).Error)
	errs = append(errs, src.Find(&avatars).Error)
	errs = append(errs, src.Find(&ips).Error)
	errs = append(errs, src.Find(&challenges).Error)
	errs = append(errs, src.Find(&usages).Error)
	errs = append(errs, src.Find(&flags).Error)
	errs = append(errs, src.Find(&dockers).Error)
	errs = append(errs, src.Find(&submissions).Error)
	for _, err := range errs {
		if err != nil {
			log.Logger.Errorf("Failed to fetch data: %v", err)
			return
		}
	}
	for _, v := range admins {
		err = dst.Create(&v).Error
		if err != nil {
			log.Logger.Errorf("Failed to insert data: %v", err)
		}
	}
	// Preload(clause.Associations)
	for _, v := range users {
		err = dst.Create(&v).Error
		if err != nil {
			log.Logger.Errorf("Failed to insert data: %v", err)
		}
	}
	for _, v := range avatars {
		err = dst.Create(&v).Error
		if err != nil {
			log.Logger.Errorf("Failed to insert data: %v", err)
		}
	}
	for _, v := range ips {
		err = dst.Create(&v).Error
		if err != nil {
			log.Logger.Errorf("Failed to insert data: %v", err)
		}
	}
	for _, v := range challenges {
		err = dst.Create(&v).Error
		if err != nil {
			log.Logger.Errorf("Failed to insert data: %v", err)
		}
	}
	for _, v := range usages {
		err = dst.Create(&v).Error
		if err != nil {
			log.Logger.Errorf("Failed to insert data: %v", err)
		}
	}
	for _, v := range flags {
		err = dst.Create(&v).Error
		if err != nil {
			log.Logger.Errorf("Failed to insert data: %v", err)
		}
	}
	for _, v := range dockers {
		err = dst.Create(&v).Error
		if err != nil {
			log.Logger.Errorf("Failed to insert data: %v", err)
		}
	}
	for _, v := range submissions {
		err = dst.Create(&v).Error
		if err != nil {
			log.Logger.Errorf("Failed to insert data: %v", err)
		}
	}
	fmt.Println("Migration completed successfully!")
}
