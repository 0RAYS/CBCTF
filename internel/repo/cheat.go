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
	UserID    uint
	TeamID    uint
	ContestID uint
	Reason    string
	Type      string
	Checked   bool
}

type UpdateCheatOptions struct {
	Reason  *string `json:"reason"`
	Type    *string `json:"type"`
	Checked *bool   `json:"checked"`
}

func InitCheatRepo(tx *gorm.DB) *CheatRepo {
	return &CheatRepo{Repo: Repo[model.Cheat]{DB: tx, Model: "Cheat"}}
}

//func (c *CheatRepo) Create(options CreateCheatOptions) (model.Cheat, bool, string) {
//	cheat, err := utils.S2S[model.Cheat](options)
//	if err != nil {
//		log.Logger.Warningf("Failed to convert options to model.Cheat: %s", err)
//		return model.Cheat{}, false, "Options2ModelError"
//	}
//	if res := c.DB.Model(&model.Cheat{}).Create(&cheat); res.Error != nil {
//		log.Logger.Warningf("Failed to create Cheat: %s", res.Error)
//		return model.Cheat{}, false, "CreateCheatError"
//	}
//	return cheat, true, "Success"
//}

//func (c *CheatRepo) getByUniqueKey(key string, value interface{}, preload bool, depth int) (model.Cheat, bool, string) {
//	switch key {
//	case "id":
//		value = value.(uint)
//	default:
//		return model.Cheat{}, false, "UnsupportedKey"
//	}
//	var cheat model.Cheat
//	res := c.DB.Model(&model.Cheat{}).Where(key+" = ?", value)
//	res = model.GetPreload(res, model.Cheat{}, preload, depth).Find(&cheat).Limit(1)
//	if res.RowsAffected == 0 {
//		return model.Cheat{}, false, "CheatNotFound"
//	}
//	return cheat, true, "Success"
//}

//func (c *CheatRepo) GetByID(id uint, preload bool, depth int) (model.Cheat, bool, string) {
//	return c.getByUniqueKey("id", id, preload, depth)
//}

//func (c *CheatRepo) Count() (int64, bool, string) {
//	var count int64
//	res := c.DB.Model(&model.Cheat{}).Count(&count)
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to count Cheats: %s", res.Error)
//		return 0, false, "CountModelError"
//	}
//	return count, true, "Success"
//}

//func (c *CheatRepo) GetAll(limit, offset int, preload bool, depth int) ([]model.Cheat, int64, bool, string) {
//	var (
//		cheats         = make([]model.Cheat, 0)
//		count, ok, msg = c.Count()
//	)
//	if !ok {
//		return cheats, count, false, msg
//	}
//	res := c.DB.Model(&model.Cheat{})
//	res = model.GetPreload(res, model.Cheat{}, preload, depth).Find(&cheats).Limit(limit).Offset(offset)
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to get Cheats: %s", res.Error)
//		return cheats, count, false, "GetCheatError"
//	}
//	return cheats, count, true, "Success"
//}

func (c *CheatRepo) Update(id uint, options UpdateCheatOptions) (bool, string) {
	var count int
	data := utils.UpdateOptions2Map(options)
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update Cheat: too many times failed due to optimistic lock")
			return false, "DeadLock"
		}
		cheat, ok, msg := c.GetByID(id, false, 0)
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

//func (c *CheatRepo) Delete(idL ...uint) (bool, string) {
//	res := c.DB.Model(&model.Cheat{}).Where("id IN ?", idL).Delete(&model.Cheat{})
//	if res.Error != nil {
//		log.Logger.Warningf("Failed to delete Cheat: %s", res.Error)
//		return false, "DeleteCheatError"
//	}
//	return true, "Success"
//}
