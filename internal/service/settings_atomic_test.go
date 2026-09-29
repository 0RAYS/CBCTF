package service

import (
	"errors"
	"testing"

	"gorm.io/gorm"

	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"
)

func TestSettingsReadFailureDoesNotPublishPartialConfiguration(t *testing.T) {
	previous := config.Env
	config.Env = &config.Config{Host: "old"}
	t.Cleanup(func() { config.Env = previous })
	database := queryFailureDB(t)
	reads := 0
	_ = database.Callback().Query().After("gorm:query").Register("test:settings_read", func(tx *gorm.DB) {
		reads++
		if reads == 2 {
			_ = tx.AddError(errors.New("database unavailable"))
			return
		}
		setting := tx.Statement.Dest.(*model.Setting)
		setting.Value = model.SettingValue{V: "new"}
		tx.RowsAffected = 1
	})
	ret := db.InitSettingRepo(database).ReadSettings()
	if ret.OK || config.Env.Host != "old" {
		t.Fatalf("partial configuration published: %+v %+v", ret, config.Env)
	}
}

func TestSettingsBatchWriteFailureRollsBack(t *testing.T) {
	previous := config.Env
	config.Env = &config.Config{Host: "old"}
	t.Cleanup(func() { config.Env = previous })
	database := queryFailureDB(t)
	pool, _ := database.DB()
	recorder := &transactionRecorder{ConnPool: pool}
	database.Statement.ConnPool = recorder
	_ = database.Callback().Query().After("gorm:query").Register("test:settings", func(tx *gorm.DB) {
		setting := tx.Statement.Dest.(*model.Setting)
		setting.ID = 1
		setting.Key = "test"
		tx.RowsAffected = 1
	})
	writes := 0
	_ = database.Callback().Update().After("gorm:update").Register("test:settings_write", func(tx *gorm.DB) {
		writes++
		tx.RowsAffected = 1
		if writes == 2 {
			_ = tx.AddError(errors.New("write failed"))
		}
	})
	ret := UpdateSystemSettings(database, dto.UpdateSettingForm{Host: new("new"), GinMode: new("release")})
	if ret.OK || recorder.rollbacks != 1 || recorder.commits != 0 || config.Env.Host != "old" {
		t.Fatalf("non-atomic settings update: %+v %+v", ret, recorder)
	}
}
