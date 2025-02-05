package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"fmt"
	"gorm.io/gorm"
)

// associations 排除过使用了 Gorm 关联关系的其余关联关系
var associations = map[string][]interface{}{
	"user_id":      {model.Submission{}},
	"team_id":      {model.Docker{}, model.Submission{}, model.Flag{}, model.Scoreboard{}},
	"challenge_id": {model.Flag{}, model.Submission{}, model.Usage{}, model.Docker{}},
	"contest_id":   {model.Docker{}, model.Flag{}, model.Submission{}, model.Usage{}, model.Scoreboard{}},
}

// ClearByID 清除所有与指定 ID 相关的数据
func ClearByID(tx *gorm.DB, column string, id interface{}) {
	var ok bool
	switch column {
	case "user_id", "team_id", "contest_id", "docker_id":
		id, ok = id.(uint)
	case "challenge_id":
		id, ok = id.(string)
	}
	if !ok {
		log.Logger.Warningf("Invalid type of id: %v", id)
		return
	}
	for _, v := range associations[column] {
		if err := tx.Model(&v).Where(fmt.Sprintf("%s = ?", column), id).Delete(&v).Error; err != nil {
			log.Logger.Warningf("Failed to delete %s by %s: %v", v, column, err)

			return
		}
	}
}
