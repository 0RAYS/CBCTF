package db

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
)

type BaseRepo[M model.Model] struct {
	DB *gorm.DB
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

func (b *BaseRepo[M]) IsUniqueKeyValue(id uint, key string, value any) (bool, model.RetVal) {
	m, ret := b.GetByUniqueField(key, value)
	if ret.OK {
		return m.GetBaseModel().ID == id, model.SuccessRetVal()
	}
	if ret.Msg == i18n.Model.NotFound {
		return true, model.SuccessRetVal()
	}
	return false, ret
}

func (b *BaseRepo[M]) IsUniqueConditionValue(id uint, conditions map[string]any) (bool, model.RetVal) {
	m, ret := b.Get(GetOptions{Conditions: conditions})
	if ret.OK {
		return m.GetBaseModel().ID == id, model.SuccessRetVal()
	}
	if ret.Msg == i18n.Model.NotFound {
		return true, model.SuccessRetVal()
	}
	return false, ret
}

// conditionTarget renders the lookup conditions as a stable target string so
// that a failure reports what was actually looked up.
func conditionTarget(conditions map[string]any) string {
	if len(conditions) == 0 {
		return ""
	}
	keys := make([]string, 0, len(conditions))
	for key := range conditions {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", key, conditions[key]))
	}
	return strings.Join(parts, ",")
}

func (b *BaseRepo[M]) Create(m M) (M, model.RetVal) {
	for _, key := range model.UniqueFields(m) {
		value := utils.GetFieldByJSONTag(m, key)
		unique, ret := b.IsUniqueKeyValue(0, key, value)
		if !ret.OK {
			return *new(M), ret
		}
		if !unique {
			return *new(M), model.FailRetVal(i18n.Model.DuplicateKeyValue, model.Name(m), fmt.Sprintf("%s=%v", key, value), nil).
				WithAttrs(map[string]any{"Key": key})
		}
	}
	for _, fields := range model.UniqueIndexes(m) {
		conditions := make(map[string]any, len(fields))
		for _, field := range fields {
			conditions[field] = utils.GetFieldByJSONTag(m, field)
		}
		unique, ret := b.IsUniqueConditionValue(0, conditions)
		if !ret.OK {
			return *new(M), ret
		}
		if !unique {
			return *new(M), model.FailRetVal(i18n.Model.DuplicateKeyValue, model.Name(m), conditionTarget(conditions), nil).
				WithAttrs(map[string]any{"Key": strings.Join(fields, ",")})
		}
	}
	target := ""
	if keys := model.UniqueFields(m); len(keys) > 0 {
		conditions := make(map[string]any, len(keys))
		for _, key := range keys {
			conditions[key] = utils.GetFieldByJSONTag(m, key)
		}
		target = conditionTarget(conditions)
	}
	if res := b.DB.Model(new(M)).Create(&m); res.Error != nil {
		log.Logger.Warningf("Failed to create %s: %s (target=%s)", model.Name(m), res.Error, target)
		return *new(M), model.FailRetVal(i18n.Model.CreateError, model.Name(m), target, res.Error)
	}
	return m, model.SuccessRetVal()
}

func applyGetOptions(tx *gorm.DB, options GetOptions) *gorm.DB {
	if conditions := options.Conditions; len(conditions) > 0 {
		tx = tx.Where(conditions)
	}
	if search := options.Search; len(search) > 0 {
		for key, value := range search {
			tx = tx.Where(fmt.Sprintf("%s ILIKE ?", key), "%"+value+"%")
		}
	}
	if preloads := options.Preloads; len(preloads) > 0 {
		for rel, subOptions := range preloads {
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
	return applyGetOptions(tx, GetOptions{
		Conditions: options.Conditions,
		Search:     options.Search,
		Deleted:    options.Deleted,
	})
}

func (b *BaseRepo[M]) Get(options GetOptions) (M, model.RetVal) {
	var m M
	target := conditionTarget(options.Conditions)
	res := applyGetOptions(b.DB.Model(new(M)), options).Limit(1).Find(&m)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get %s: %s (target=%s)", model.Name(m), res.Error, target)
		return *new(M), model.FailRetVal(i18n.Model.GetError, model.Name(m), target, res.Error)
	}
	if res.RowsAffected == 0 {
		return *new(M), model.FailRetVal(i18n.Model.NotFound, model.Name(m), target, nil)
	}
	return m, model.SuccessRetVal()
}

func (b *BaseRepo[M]) GetByID(id uint, options ...GetOptions) (M, model.RetVal) {
	return b.GetByUniqueField("id", id, options...)
}

func (b *BaseRepo[M]) GetByIDForUpdate(id uint, optionsL ...GetOptions) (M, model.RetVal) {
	options := GetOptions{}
	if len(optionsL) > 0 {
		options = optionsL[0]
	}
	if options.Conditions == nil {
		options.Conditions = make(map[string]any)
	}
	options.Conditions["id"] = id

	var m M
	target := conditionTarget(options.Conditions)
	res := applyGetOptions(b.DB.Model(new(M)).Clauses(clause.Locking{Strength: "UPDATE"}), options).Limit(1).Find(&m)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get %s for update: %s (target=%s)", model.Name(m), res.Error, target)
		return *new(M), model.FailRetVal(i18n.Model.GetError, model.Name(m), target, res.Error)
	}
	if res.RowsAffected == 0 {
		return *new(M), model.FailRetVal(i18n.Model.NotFound, model.Name(*new(M)), target, nil)
	}
	return m, model.SuccessRetVal()
}

func (b *BaseRepo[M]) GetByUniqueField(key string, value any, optionsL ...GetOptions) (M, model.RetVal) {
	if !slices.Contains(model.UniqueFields(*new(M)), key) {
		return *new(M), model.FailRetVal(i18n.Model.NotUniqueKey, model.Name(*new(M)), key, nil).
			WithAttrs(map[string]any{"Key": key})
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
	target := ""
	if len(optionsL) > 0 {
		res = applyCountOptions(res, optionsL[0])
		target = conditionTarget(optionsL[0].Conditions)
	}
	if res = res.Count(&count); res.Error != nil {
		log.Logger.Warningf("Failed to count %s: %s (target=%s)", model.Name(*new(M)), res.Error, target)
		return 0, model.FailRetVal(i18n.Model.GetError, model.Name(*new(M)), target, res.Error)
	}
	return count, model.SuccessRetVal()
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
			Search:     options.Search,
			Deleted:    options.Deleted,
		})
	)
	if !ret.OK {
		return nil, count, ret
	}
	tx := applyGetOptions(b.DB.Model(new(M)), options)
	if len(options.Sort) == 0 {
		tx = tx.Order("id")
	}
	target := conditionTarget(options.Conditions)
	if res := tx.Limit(limit).Offset(offset).Find(&ms); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, count, model.FailRetVal(i18n.Model.NotFound, model.Name(*new(M)), target, nil)
		}
		log.Logger.Warningf("Failed to get %s: %s (target=%s)", model.Name(*new(M)), res.Error, target)
		return nil, count, model.FailRetVal(i18n.Model.GetError, model.Name(*new(M)), target, res.Error)
	}
	return ms, count, model.SuccessRetVal()
}

func (b *BaseRepo[M]) FindAll(optionsL ...GetOptions) ([]M, model.RetVal) {
	options := GetOptions{}
	if len(optionsL) > 0 {
		options = optionsL[0]
	}
	ms := make([]M, 0)
	tx := applyGetOptions(b.DB.Model(new(M)), options)
	if len(options.Sort) == 0 {
		tx = tx.Order("id")
	}
	target := conditionTarget(options.Conditions)
	if res := tx.Find(&ms); res.Error != nil {
		log.Logger.Warningf("Failed to get %s: %s (target=%s)", model.Name(*new(M)), res.Error, target)
		return nil, model.FailRetVal(i18n.Model.GetError, model.Name(*new(M)), target, res.Error)
	}
	return ms, model.SuccessRetVal()
}

func (b *BaseRepo[M]) Update(id uint, options UpdateOptions) model.RetVal {
	var count uint
	data := options.Convert2Map()
	if len(data) == 0 {
		return model.SuccessRetVal()
	}
	target := fmt.Sprintf("id=%d", id)
	for _, key := range model.UniqueFields(*new(M)) {
		if value, ok := data[key]; ok {
			unique, ret := b.IsUniqueKeyValue(id, key, value)
			if !ret.OK {
				return ret
			}
			if !unique {
				return model.FailRetVal(i18n.Model.NotUniqueKey, model.Name(*new(M)), fmt.Sprintf("%s=%v", key, value), nil).
					WithAttrs(map[string]any{"Key": key})
			}
		}
	}
	var current M
	var hasCurrent bool
	getCurrent := func() (M, model.RetVal) {
		if hasCurrent {
			return current, model.SuccessRetVal()
		}
		m, ret := b.GetByID(id)
		if !ret.OK {
			return *new(M), ret
		}
		current = m
		hasCurrent = true
		return current, model.SuccessRetVal()
	}
	for _, fields := range model.UniqueIndexes(*new(M)) {
		changed := false
		conditions := make(map[string]any, len(fields))
		for _, field := range fields {
			if value, ok := data[field]; ok {
				conditions[field] = value
				changed = true
				continue
			}
			m, ret := getCurrent()
			if !ret.OK {
				return ret
			}
			conditions[field] = utils.GetFieldByJSONTag(m, field)
		}
		if !changed {
			continue
		}
		unique, ret := b.IsUniqueConditionValue(id, conditions)
		if !ret.OK {
			return ret
		}
		if !unique {
			return model.FailRetVal(i18n.Model.NotUniqueKey, model.Name(*new(M)), conditionTarget(conditions), nil).
				WithAttrs(map[string]any{"Key": strings.Join(fields, ",")})
		}
	}
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update %s: too many times failed due to optimistic lock (target=%s)", model.Name(*new(M)), target)
			return model.FailRetVal(i18n.Model.DeadLock, model.Name(*new(M)), target, nil)
		}
		m, ret := b.GetByID(id)
		if !ret.OK {
			return ret
		}
		res := b.DB.Model(&m).Where("id = ?", id).Updates(data)
		if res.Error != nil {
			log.Logger.Warningf("Failed to update %s: %s (target=%s)", model.Name(*new(M)), res.Error, target)
			return model.FailRetVal(i18n.Model.UpdateError, model.Name(*new(M)), target, res.Error)
		}
		if res.RowsAffected == 0 {
			continue
		}
		break
	}
	return model.SuccessRetVal()
}

func (b *BaseRepo[M]) DiffUpdate(id uint, options DiffUpdateOptions) model.RetVal {
	target := fmt.Sprintf("id=%d", id)
	data := options.Convert2Expr()
	res := b.DB.Model(new(M)).Where("id = ?", id).Updates(data)
	if res.Error != nil {
		log.Logger.Warningf("Failed to update %s: %s (target=%s)", model.Name(*new(M)), res.Error, target)
		return model.FailRetVal(i18n.Model.UpdateError, model.Name(*new(M)), target, res.Error)
	}
	return model.SuccessRetVal()
}

func (b *BaseRepo[M]) Delete(idL ...uint) model.RetVal {
	if len(idL) == 0 {
		return model.SuccessRetVal()
	}
	target := fmt.Sprintf("id=%v", idL)
	res := b.DB.Model(new(M)).Where("id IN ?", idL).Delete(new(M))
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete %s: %s (target=%s)", model.Name(*new(M)), res.Error, target)
		return model.FailRetVal(i18n.Model.DeleteError, model.Name(*new(M)), target, res.Error)
	}
	return model.SuccessRetVal()
}

func (b *BaseRepo[M]) DeleteByFieldID(field string, idL ...uint) model.RetVal {
	if len(idL) == 0 {
		return model.SuccessRetVal()
	}
	var targetIDL []uint
	target := fmt.Sprintf("%s=%v", field, idL)
	if res := b.DB.Model(new(M)).Where(field+" IN ?", idL).Pluck("id", &targetIDL); res.Error != nil {
		log.Logger.Warningf("Failed to get %s by %s %v: %s", model.Name(*new(M)), field, idL, res.Error)
		return model.FailRetVal(i18n.Model.DeleteError, model.Name(*new(M)), target, res.Error)
	}
	return b.Delete(targetIDL...)
}
