package repo

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Basic[M model.Model] struct {
	DB *gorm.DB
}

type CreateOptions interface {
	Convert2Model() model.Model
}

type UpdateOptions interface {
	Convert2Map() map[string]any
}

func (r *Basic[M]) Create(options CreateOptions) (M, bool, string) {
	m := options.Convert2Model()
	if res := r.DB.Model(new(M)).Create(&m); res.Error != nil {
		log.Logger.Warningf("Failed to create %T: %s", new(M), res.Error)
		return *new(M), false, M.CreateErrorString(new(M))
	}
	return m.(M), true, i18n.Success
}

func (r *Basic[M]) GetWithConditions(conditions map[string]any, preloadL ...string) (M, bool, string) {
	var m M
	res := r.DB.Model(new(M))
	for key, value := range conditions {
		res = res.Where(fmt.Sprintf("%s = ?", key), value)
	}
	res = preload(res, preloadL...).Limit(1).Find(&m)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get %s: %s", m.GetModelName(), res.Error)
		return m, false, m.GetErrorString()
	}
	if res.RowsAffected == 0 {
		return m, false, m.NotFoundErrorString()
	}
	return m, true, i18n.Success
}

func (r *Basic[M]) getUniqueByKey(key string, value any, preloadL ...string) (M, bool, string) {
	if !utils.In(key, M.GetUniqueKey(new(M))) {
		return *new(M), false, i18n.UnsupportedKey
	}
	return r.GetWithConditions(map[string]any{key: value}, preloadL...)
}

func (r *Basic[M]) GetByID(id uint, preloadL ...string) (M, bool, string) {
	return r.getUniqueByKey("id", id, preloadL...)
}

func (r *Basic[M]) CountWithConditions(conditions map[string]any) (int64, bool, string) {
	var count int64
	res := r.DB.Model(new(M))
	for key, value := range conditions {
		res = res.Where(fmt.Sprintf("%s = ?", key), value)
	}
	if res = res.Count(&count); res.Error != nil {
		log.Logger.Warningf("Failed to count %T: %s", new(M), res.Error)
		return 0, false, M.GetErrorString(new(M))
	}
	return count, true, i18n.Success
}

func (r *Basic[M]) Count() (int64, bool, string) {
	return r.CountWithConditions(nil)
}

func (r *Basic[M]) ListWithConditions(limit, offset int, conditions map[string]any, preloadL ...string) ([]M, int64, bool, string) {
	var (
		models         = make([]M, 0)
		count, ok, msg = r.CountWithConditions(conditions)
	)
	if !ok {
		return models, count, false, msg
	}
	res := r.DB.Model(new(M))
	for key, value := range conditions {
		res = res.Where(fmt.Sprintf("%s = ?", key), value)
	}
	res = preload(res, preloadL...).Order("created_at ASC").Limit(limit).Offset(offset).Find(&models)
	if res.Error != nil {
		log.Logger.Errorf("Failed to get %s: %s", M.GetModelName(new(M)), res.Error)
		return models, count, false, M.GetErrorString(new(M))
	}
	return models, count, true, i18n.Success
}

func (r *Basic[M]) List(limit, offset int, preloadL ...string) ([]M, int64, bool, string) {
	return r.ListWithConditions(limit, offset, nil, preloadL...)
}

func (r *Basic[M]) Update(id uint, options UpdateOptions) (bool, string) {
	var count uint
	data := options.Convert2Map()
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update %s: too many times failed due to optimistic lock", M.GetModelName(new(M)))
			return false, i18n.DeadLock
		}
		m, ok, msg := r.GetByID(id)
		if !ok {
			return ok, msg
		}
		data["version"] = m.GetVersion() + 1
		res := r.DB.Model(new(M)).Where("id = ? AND version = ?", id, m.GetVersion()).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update %s: %s", M.GetModelName(new(M)), res.Error)
			return false, M.UpdateErrorString(new(M))
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, i18n.Success
}

func (r *Basic[M]) Delete(idL ...uint) (bool, string) {
	if res := r.DB.Model(new(M)).Where("id IN ?", idL).Delete(new(M)); res.Error != nil {
		log.Logger.Warningf("Failed to delete %s: %s", M.GetModelName(new(M)), res.Error)
		return false, M.DeleteErrorString(new(M))
	}
	return true, i18n.Success
}

func preload(tx *gorm.DB, preloadL ...string) *gorm.DB {
	for _, nested := range preloadL {
		if nested == "all" {
			tx = tx.Preload(clause.Associations)
		} else {
			tx = tx.Preload(nested)
		}
	}
	return tx
}
