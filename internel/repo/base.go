package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/utils"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetPreload(tx *gorm.DB, preloadL ...string) *gorm.DB {
	for _, nested := range preloadL {
		if nested == "all" {
			tx = tx.Preload(clause.Associations)
		} else {
			tx = tx.Preload(nested)
		}
	}
	return tx
}

type Repo[T any] struct {
	DB    *gorm.DB
	Model string
}

func (r *Repo[T]) Create(options interface{}) (T, bool, string) {
	m, err := utils.S2S[T](options)
	if err != nil {
		log.Logger.Warningf("Failed to convert options to %T: %s", new(T), err)
		return *new(T), false, fmt.Sprintf("Options2%sError", r.Model)
	}
	if res := r.DB.Model(new(T)).Create(&m); res.Error != nil {
		log.Logger.Warningf("Failed to create %T: %s", new(T), res.Error)
		return *new(T), false, fmt.Sprintf("Create%sError", r.Model)
	}
	return m, true, "Success"
}

func (r *Repo[T]) getByUniqueKey(key string, value interface{}, preloadL ...string) (T, bool, string) {
	switch key {
	case "id":
		value = value.(uint)
	default:
		return *new(T), false, "UnsupportedKey"
	}
	var m T
	res := r.DB.Model(new(T)).Where(key+" = ?", value)
	res = GetPreload(res, preloadL...).Limit(1).Find(&m)
	if res.RowsAffected == 0 {
		return m, false, fmt.Sprintf("%sNotFound", r.Model)
	}
	return m, true, "Success"
}

func (r *Repo[T]) GetByID(id uint, preloadL ...string) (T, bool, string) {
	return r.getByUniqueKey("id", id, preloadL...)
}

func (r *Repo[T]) Count() (int64, bool, string) {
	var count int64
	if res := r.DB.Model(new(T)).Count(&count); res.Error != nil {
		log.Logger.Warningf("Failed to count %T: %s", new(T), res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (r *Repo[T]) GetAll(limit, offset int, preloadL ...string) ([]T, int64, bool, string) {
	var (
		ms             = make([]T, 0)
		count, ok, msg = r.Count()
	)
	if !ok {
		return ms, count, false, msg
	}
	res := r.DB.Model(new(T))
	res = GetPreload(res, preloadL...).Limit(limit).Offset(offset).Find(&ms)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get all %T: %s", new(T), res.Error)
		return ms, count, false, fmt.Sprintf("Get%sError", r.Model)
	}
	return ms, count, true, "Success"
}

func (r *Repo[T]) Delete(idL ...uint) (bool, string) {
	if res := r.DB.Model(new(T)).Where("id IN ?", idL).Delete(new(T)); res.Error != nil {
		log.Logger.Warningf("Failed to delete %T: %s", new(T), res.Error)
		return false, fmt.Sprintf("Delete%sError", r.Model)
	}
	return true, "Success"
}
