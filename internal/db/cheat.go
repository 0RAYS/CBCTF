package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"fmt"
	"gorm.io/gorm"
)

// CreateCheat 创建单一事件作弊记录
func CreateCheat(tx *gorm.DB, userID, teamID, contestID uint, reason string, t string) (model.Cheat, bool, string) {
	cheat := model.InitCheat(userID, teamID, contestID, reason, t)
	return RecordCheat(tx, cheat)
}

// RecordCheat 可创建关联事件作弊记录
func RecordCheat(tx *gorm.DB, cheat model.Cheat) (model.Cheat, bool, string) {
	res := tx.Model(&model.Cheat{}).Create(&cheat)
	if res.Error != nil {
		log.Logger.Warningf("Failed to record cheat: %v", res.Error)
		return model.Cheat{}, false, "CreateCheatError"
	}
	return cheat, true, "Success"
}

// GetCheatsByColumn 根据列获取作弊记录
func GetCheatsByColumn(tx *gorm.DB, column string, id uint) ([]model.Cheat, bool, string) {
	var cheats []model.Cheat
	res := tx.Model(&model.Cheat{}).Where(fmt.Sprintf("%s = ?", column), id).Find(&cheats)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get cheats: %v", res.Error)
		return make([]model.Cheat, 0), true, "GetCheatsError"
	}
	return cheats, true, "Success"
}

// DeleteCheat 根据 ID 获取作弊记录
func DeleteCheat(tx *gorm.DB, id string) (bool, string) {
	res := tx.Model(&model.Cheat{}).Where("id = ?", id).Delete(&model.Cheat{})
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete cheat: %v", res.Error)
		return false, "DeleteCheatError"
	}
	return true, "Success"
}
