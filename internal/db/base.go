package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"errors"
	"fmt"
	"slices"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BasicRepo[M model.Model] struct {
	DB *gorm.DB
}

type CreateOptions interface {
	Convert2Model() model.Model
}

// GetOptions 设置 Preload 时, 须确保对应关系的外键被 Select
type GetOptions struct {
	Conditions       map[string]any
	SearchConditions map[string]any
	Selects          []string
	Preloads         map[string]GetOptions
	Deleted          bool
}

type CountOptions struct {
	Conditions       map[string]any
	SearchConditions map[string]any
	Deleted          bool
}

type UpdateOptions interface {
	Convert2Map() map[string]any
}

type DiffUpdateOptions interface {
	Convert2Expr() map[string]any
}

func (b *BasicRepo[M]) Create(options CreateOptions) (M, bool, string) {
	m := options.Convert2Model().(M)
	if res := b.DB.Model(new(M)).Create(&m); res.Error != nil {
		log.Logger.Warningf("Failed to create %T: %s", new(M), res.Error)
		return *new(M), false, M.CreateErrorString(*new(M))
	}
	return m, true, i18n.Success
}

func ApplyGetOptions(tx *gorm.DB, options GetOptions) *gorm.DB {
	if columns := options.Selects; len(columns) > 0 {
		if !slices.Contains(columns, "id") {
			columns = append([]string{"id"}, columns...)
		}
		tx = tx.Select(columns)
	}
	if conditions := options.Conditions; len(conditions) > 0 {
		tx = tx.Where(conditions)
	}
	if search := options.SearchConditions; len(search) > 0 {
		for key, value := range search {
			tx = tx.Where(fmt.Sprintf("%s LIKE ?", key), "%"+value.(string)+"%")
		}
	}
	if options.Deleted {
		tx = tx.Unscoped()
	}
	if preloads := options.Preloads; preloads != nil {
		for rel, subOptions := range preloads {
			if rel == "all" {
				tx = tx.Preload(clause.Associations)
				continue
			}
			tx = tx.Preload(rel, func(tx *gorm.DB) *gorm.DB {
				return ApplyGetOptions(tx, subOptions)
			})
		}
	}
	return tx
}

func (b *BasicRepo[M]) Get(options GetOptions) (M, bool, string) {
	var m M
	res := ApplyGetOptions(b.DB.Model(new(M)), options).Limit(1).Find(&m)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get %s: %s", m.GetModelName(), res.Error)
		return *new(M), false, m.GetErrorString()
	}
	if res.RowsAffected == 0 {
		return *new(M), false, m.NotFoundErrorString()
	}
	return m, true, i18n.Success
}

func (b *BasicRepo[M]) GetByID(id uint, optionsL ...GetOptions) (M, bool, string) {
	options := GetOptions{}
	if len(optionsL) > 0 {
		options = optionsL[0]
	}
	return b.GetByUniqueKey("id", id, options)
}

func (b *BasicRepo[M]) GetByUniqueKey(key string, value any, optionsL ...GetOptions) (M, bool, string) {
	if !slices.Contains(M.GetUniqueKey(*new(M)), key) {
		return *new(M), false, i18n.UnsupportedKey
	}
	options := GetOptions{}
	if len(optionsL) > 0 {
		options = optionsL[0]
	}
	if options.Conditions == nil {
		options.Conditions = make(map[string]any)
	}
	options.Conditions[key] = value
	return b.Get(options)
}

func (b *BasicRepo[M]) Count(optionsL ...CountOptions) (int64, bool, string) {
	var count int64
	res := b.DB.Model(new(M))
	if len(optionsL) > 0 {
		options := optionsL[0]
		if conditions := options.Conditions; len(conditions) > 0 {
			res = res.Where(conditions)
		}
		if search := options.SearchConditions; len(search) > 0 {
			for key, value := range search {
				res = res.Where(fmt.Sprintf("%s LIKE ?", key), "%"+value.(string)+"%")
			}
		}
		if options.Deleted {
			res = res.Unscoped()
		}
	}
	if res = res.Count(&count); res.Error != nil {
		log.Logger.Warningf("Failed to count %s: %s", M.GetModelName(*new(M)), res.Error)
		return 0, false, M.GetErrorString(*new(M))
	}
	return count, true, i18n.Success
}

func (b *BasicRepo[M]) CountAssociation(m M, association string) int64 {
	return b.DB.Model(&m).Association(association).Count()
}

func (b *BasicRepo[M]) List(limit, offset int, optionsL ...GetOptions) ([]M, int64, bool, string) {
	options := GetOptions{}
	if len(optionsL) > 0 {
		options = optionsL[0]
	}
	var (
		ms             = make([]M, 0)
		count, ok, msg = b.Count(CountOptions{
			Conditions: options.Conditions,
			Deleted:    options.Deleted,
		})
	)
	if !ok {
		return nil, count, false, msg
	}
	if res := ApplyGetOptions(b.DB.Model(new(M)), options).Order("id").Limit(limit).Offset(offset).Find(&ms); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, count, false, M.NotFoundErrorString(*new(M))
		}
		log.Logger.Warningf("Failed to get %s: %s", M.GetModelName(*new(M)), res.Error)
		return nil, count, false, M.GetErrorString(*new(M))
	}
	return ms, count, true, i18n.Success
}

func (b *BasicRepo[M]) Update(id uint, options UpdateOptions) (bool, string) {
	var count uint
	data := options.Convert2Map()
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update %s: too many times failed due to optimistic lock", M.GetModelName(*new(M)))
			return false, i18n.DeadLock
		}
		m, ok, msg := b.GetByID(id, GetOptions{Selects: []string{"id", "version"}})
		if !ok {
			return ok, msg
		}
		res := b.DB.Model(&m).Where("id = ?", id).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update %s: %s", M.GetModelName(*new(M)), res.Error)
			return false, M.UpdateErrorString(*new(M))
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return true, i18n.Success
}

func (b *BasicRepo[M]) DiffUpdate(id uint, options DiffUpdateOptions) (bool, string) {
	data := options.Convert2Expr()
	res := b.DB.Model(new(M)).Where("id = ?", id).Updates(data)
	if res.Error != nil {
		log.Logger.Warningf("Failed to update %s: %s", M.GetModelName(*new(M)), res.Error)
		return false, M.UpdateErrorString(*new(M))
	}
	return true, i18n.Success
}

func (b *BasicRepo[M]) Delete(idL ...uint) (bool, string) {
	res := b.DB.Model(new(M)).Where("id IN ?", idL).Delete(new(M))
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete %s: %s", M.GetModelName(*new(M)), res.Error)
		return false, M.DeleteErrorString(*new(M))
	}
	return true, i18n.Success
}
