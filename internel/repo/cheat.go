package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
)

type CheatRepo struct {
	Repo[model.Cheat]
}

type CreateCheatOptions struct {
	ID        string
	UserID    uint
	TeamID    uint
	ContestID uint
	Reason    string
	Type      string
	Checked   bool
	Cheats    model.Strings
}

type UpdateCheatOptions struct {
	Reason  *string `json:"reason"`
	Type    *string `json:"type"`
	Checked *bool   `json:"checked"`
}

func InitCheatRepo(tx *gorm.DB) *CheatRepo {
	return &CheatRepo{Repo: Repo[model.Cheat]{DB: tx, Model: "Cheat"}}
}

func (c *CheatRepo) Update(id uint, options UpdateCheatOptions) (bool, string) {
	var count int
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Cheat: too many times failed due to optimistic lock")
			return false, "DeadLock"
		}
		cheat, ok, msg := c.GetByID(id)
		if !ok {
			return ok, msg
		}
		data["version"] = cheat.Version + 1
		res := c.DB.Model(&model.Cheat{}).Omit("id", "created_at", "updated_at", "deleted_at").
			Where("id = ? AND version = ?", id, cheat.Version).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update Cheat: %s", res.Error)
			return false, "UpdateCheatError"
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, "Success"
}
