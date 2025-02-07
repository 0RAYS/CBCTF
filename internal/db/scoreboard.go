package db

import "gorm.io/gorm"

func UpdateScoreboard(tx *gorm.DB, contestID uint) (bool, string) {

	return true, "Success"
}
