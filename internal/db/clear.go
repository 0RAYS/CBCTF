package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"fmt"
	"gorm.io/gorm"
)

// associations 排除过使用了 Gorm 关联关系的其余关联关系, 不包含 model.Docker, 由定时任务删除
var associations = map[string][]interface{}{
	"user_id":      {model.Submission{}},
	"team_id":      {model.Submission{}, model.Flag{}},
	"challenge_id": {model.Flag{}, model.Submission{}, model.Usage{}},
	"contest_id":   {model.Flag{}, model.Submission{}, model.Usage{}, model.Notice{}},
	"usage_id":     {model.Submission{}},
}

// ClearByID 清除所有与指定 ID 相关的数据
func ClearByID(tx *gorm.DB, column string, id interface{}) bool {
	var ok bool
	switch column {
	case "user_id", "team_id", "contest_id", "docker_id", "usage_id":
		id, ok = id.(uint)
	case "challenge_id":
		id, ok = id.(string)
	}
	if !ok {
		log.Logger.Warningf("Invalid type of id: %v", id)
		return false
	}
	for _, v := range associations[column] {
		if err := tx.Model(&v).Where(fmt.Sprintf("%s = ?", column), id).Delete(&v).Error; err != nil {
			log.Logger.Warningf("Failed to delete %s by %s: %v", v, column, err)
			return false
		}
	}
	return true
}
