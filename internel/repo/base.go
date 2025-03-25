package repo

import (
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
)

type Repo[T any] struct {
	DB    *gorm.DB
	Model string
}

func (r *Repo[T]) Create(options interface{}) (T, bool, string) {
	m, err := utils.S2S[T](options)
	if err != nil {
		log.Logger.Warningf("Failed to convert options to %T: %s", new(T), err)
		return *new(T), false, "Options2ModelError"
	}
	if res := r.DB.Model(new(T)).Create(&m); res.Error != nil {
		log.Logger.Warningf("Failed to create %T: %s", new(T), res.Error)
		return *new(T), false, "CreateModelError"
	}
	return m, true, "Success"
}

func (r *Repo[T]) getByUniqueKey(key string, value interface{}, preload bool, depth int) (T, bool, string) {
	switch key {
	case "id":
		value = value.(uint)
	default:
		return *new(T), false, "UnsupportedKey"
	}
	var m T
	res := r.DB.Model(new(T)).Where(key+" = ?", value)
	res = model.GetPreload(res, r.Model, preload, depth).Find(&m).Limit(1)
	if res.RowsAffected == 0 {
		return m, false, "ModelNotFound"
	}
	return m, true, "Success"
}

func (r *Repo[T]) GetByID(id uint, preload bool, depth int) (T, bool, string) {
	return r.getByUniqueKey("id", id, preload, depth)
}

func (r *Repo[T]) Count() (int64, bool, string) {
	var count int64
	if res := r.DB.Model(new(T)).Count(&count); res.Error != nil {
		log.Logger.Warningf("Failed to count %T: %s", new(T), res.Error)
		return 0, false, "CountModelError"
	}
	return count, true, "Success"
}

func (r *Repo[T]) GetAll(limit, offset int, preload bool, depth int) ([]T, int64, bool, string) {
	var (
		ms             = make([]T, 0)
		count, ok, msg = r.Count()
	)
	if !ok {
		return ms, count, false, msg
	}
	res := r.DB.Model(new(T))
	res = model.GetPreload(res, r.Model, preload, depth).Find(&ms).Limit(limit).Offset(offset)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get all %T: %s", new(T), res.Error)
		return ms, count, false, "GetModelError"
	}
	return ms, count, true, "Success"
}

func (r *Repo[T]) Delete(idL ...uint) (bool, string) {
	if res := r.DB.Model(new(T)).Where("id IN ?", idL).Delete(new(T)); res.Error != nil {
		log.Logger.Warningf("Failed to delete %T: %s", new(T), res.Error)
		return false, "DeleteModelError"
	}
	return true, "Success"
}
