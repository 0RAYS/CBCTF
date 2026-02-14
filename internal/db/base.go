package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"errors"
	"fmt"
	"slices"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseRepo[M model.Model] struct {
	DB *gorm.DB
}

type CreateOptions interface {
	Convert2Model() model.Model
}

type GetOptions struct {
	Conditions map[string]any
	Search     map[string]string
	Preloads   map[string]GetOptions
	Deleted    bool
	Sort       []string
}

type CountOptions struct {
	Conditions map[string]any
	Search     map[string]string
	Deleted    bool
}

type UpdateOptions interface {
	Convert2Map() map[string]any
}

type DiffUpdateOptions interface {
	Convert2Expr() map[string]any
}

func (b *BaseRepo[M]) IsUniqueKeyValue(id uint, key string, value any) bool {
	m, ret := b.GetByUniqueKey(key, value)
	return m.GetBaseModel().ID == id || !ret.OK
}

func (b *BaseRepo[M]) Insert(m M) (M, model.RetVal) {
	for _, key := range m.UniqueFields() {
		value := utils.GetFieldByJSONTag(m, key)
		if !b.IsUniqueKeyValue(0, key, value) {
			return *new(M), model.RetVal{Msg: i18n.Model.DuplicateKeyValue, Attr: map[string]any{"Model": m.ModelName(), "Key": key}}
		}
	}
	if res := b.DB.Model(new(M)).Create(&m); res.Error != nil {
		log.Logger.Warningf("Failed to create %T: %s", new(M), res.Error)
		return *new(M), model.RetVal{Msg: i18n.Model.CreateError, Attr: map[string]any{"Model": m.ModelName(), "Error": res.Error.Error()}}
	}
	return m, model.SuccessRetVal()
}

func (b *BaseRepo[M]) Create(options CreateOptions) (M, model.RetVal) {
	return b.Insert(options.Convert2Model().(M))
}

func applyGetOptions(tx *gorm.DB, options GetOptions) *gorm.DB {
	if conditions := options.Conditions; len(conditions) > 0 {
		tx = tx.Where(conditions)
	}
	if search := options.Search; len(search) > 0 {
		for key, value := range search {
			tx = tx.Where(fmt.Sprintf("%s LIKE ?", key), "%"+value+"%")
		}
	}
	if preloads := options.Preloads; preloads != nil {
		for rel, subOptions := range preloads {
			if rel == "all" {
				tx = tx.Preload(clause.Associations)
				continue
			}
			tx = tx.Preload(rel, func(tx *gorm.DB) *gorm.DB {
				return applyGetOptions(tx, subOptions)
			})
		}
	}
	if options.Deleted {
		tx = tx.Unscoped()
	}
	if columns := options.Sort; len(columns) > 0 {
		for _, order := range columns {
			tx = tx.Order(order)
		}
	}
	return tx
}

func applyCountOptions(tx *gorm.DB, options CountOptions) *gorm.DB {
	if conditions := options.Conditions; len(conditions) > 0 {
		tx = tx.Where(conditions)
	}
	if search := options.Search; len(search) > 0 {
		for key, value := range search {
			tx = tx.Where(fmt.Sprintf("%s LIKE ?", key), "%"+value+"%")
		}
	}
	if options.Deleted {
		tx = tx.Unscoped()
	}
	return tx
}

func (b *BaseRepo[M]) Get(options GetOptions) (M, model.RetVal) {
	var m M
	res := applyGetOptions(b.DB.Model(new(M)), options).Limit(1).Find(&m)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get %s: %s", m.ModelName(), res.Error)
		return *new(M), model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": m.ModelName(), "Error": res.Error.Error()}}
	}
	if res.RowsAffected == 0 {
		return *new(M), model.RetVal{Msg: i18n.Model.NotFound, Attr: map[string]any{"Model": m.ModelName()}}
	}
	return m, model.SuccessRetVal()
}

func (b *BaseRepo[M]) GetByID(id uint, options ...GetOptions) (M, model.RetVal) {
	return b.GetByUniqueKey("id", id, options...)
}

func (b *BaseRepo[M]) GetByUniqueKey(key string, value any, optionsL ...GetOptions) (M, model.RetVal) {
	if !slices.Contains(M.UniqueFields(*new(M)), key) {
		return *new(M), model.RetVal{Msg: i18n.Model.NotUniqueKey, Attr: map[string]any{"Model": M.ModelName(*new(M)), "Key": key}}
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

func (b *BaseRepo[M]) Count(optionsL ...CountOptions) (int64, model.RetVal) {
	var count int64
	res := b.DB.Model(new(M))
	if len(optionsL) > 0 {
		res = applyCountOptions(res, optionsL[0])
	}
	if res = res.Count(&count); res.Error != nil {
		log.Logger.Warningf("Failed to count %s: %s", M.ModelName(*new(M)), res.Error)
		return 0, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": M.ModelName(*new(M)), "Error": res.Error.Error()}}
	}
	return count, model.SuccessRetVal()
}

func (b *BaseRepo[M]) CountAssociation(m M, association string) int64 {
	return b.DB.Model(&m).Association(association).Count()
}

func (b *BaseRepo[M]) List(limit, offset int, optionsL ...GetOptions) ([]M, int64, model.RetVal) {
	options := GetOptions{}
	if len(optionsL) > 0 {
		options = optionsL[0]
	}
	var (
		ms         = make([]M, 0)
		count, ret = b.Count(CountOptions{
			Conditions: options.Conditions,
			Deleted:    options.Deleted,
		})
	)
	if !ret.OK {
		return nil, count, ret
	}
	if res := applyGetOptions(b.DB.Model(new(M)), options).Order("id").Limit(limit).Offset(offset).Find(&ms); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, count, model.RetVal{Msg: i18n.Model.NotFound, Attr: map[string]any{"Model": M.ModelName(*new(M))}}
		}
		log.Logger.Warningf("Failed to get %s: %s", M.ModelName(*new(M)), res.Error)
		return nil, count, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": M.ModelName(*new(M)), "Error": res.Error.Error()}}
	}
	return ms, count, model.SuccessRetVal()
}

func (b *BaseRepo[M]) Update(id uint, options UpdateOptions) model.RetVal {
	var count uint
	data := options.Convert2Map()
	if len(data) == 0 {
		return model.SuccessRetVal()
	}
	for _, key := range M.UniqueFields(*new(M)) {
		if value, ok := data[key]; ok && !b.IsUniqueKeyValue(id, key, value) {
			return model.RetVal{Msg: i18n.Model.NotUniqueKey, Attr: map[string]any{"Model": M.ModelName(*new(M)), "Key": key}}
		}
	}
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update %s: too many times failed due to optimistic lock", M.ModelName(*new(M)))
			return model.RetVal{Msg: i18n.Model.DeadLock, Attr: map[string]any{"Model": M.ModelName(*new(M))}}
		}
		m, ret := b.GetByID(id)
		if !ret.OK {
			return ret
		}
		res := b.DB.Model(&m).Where("id = ?", id).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update %s: %s", M.ModelName(*new(M)), res.Error)
			return model.RetVal{Msg: i18n.Model.UpdateError, Attr: map[string]any{"Model": M.ModelName(*new(M)), "Error": res.Error.Error()}}
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return model.SuccessRetVal()
}

func (b *BaseRepo[M]) DiffUpdate(id uint, options DiffUpdateOptions) model.RetVal {
	data := options.Convert2Expr()
	res := b.DB.Model(new(M)).Where("id = ?", id).Updates(data)
	if res.Error != nil {
		log.Logger.Warningf("Failed to update %s: %s", M.ModelName(*new(M)), res.Error)
		return model.RetVal{Msg: i18n.Model.UpdateError, Attr: map[string]any{"Model": M.ModelName(*new(M)), "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func (b *BaseRepo[M]) Delete(idL ...uint) model.RetVal {
	res := b.DB.Model(new(M)).Where("id IN ?", idL).Delete(new(M))
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete %s: %s", M.ModelName(*new(M)), res.Error)
		return model.RetVal{Msg: i18n.Model.DeleteError, Attr: map[string]any{"Model": M.ModelName(*new(M)), "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}
