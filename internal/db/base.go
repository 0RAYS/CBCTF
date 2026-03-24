package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"errors"
	"fmt"
	"slices"

	"github.com/jackc/pgx/v5/pgconn"
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
	m, ret := b.GetByUniqueField(key, value)
	return m.GetBaseModel().ID == id || !ret.OK
}

func (b *BaseRepo[M]) Insert(m M) (M, model.RetVal) {
	if res := b.DB.Model(new(M)).Create(&m); res.Error != nil {
		if attr, ok := duplicateKeyAttr(res.Error); ok {
			attr["Model"] = model.ModelName(m)
			return *new(M), model.RetVal{Msg: i18n.Model.DuplicateKeyValue, Attr: attr}
		}
		log.Logger.Warningf("Failed to create %T: %s", new(M), res.Error)
		return *new(M), model.RetVal{Msg: i18n.Model.CreateError, Attr: map[string]any{"Model": model.ModelName(m), "Error": res.Error.Error()}}
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
	if conditions := options.Conditions; len(conditions) > 0 {
		tx = tx.Where(conditions)
	}
	if search := options.Search; len(search) > 0 {
		for key, value := range search {
			tx = tx.Where(fmt.Sprintf("%s ILIKE ?", key), "%"+value+"%")
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
		log.Logger.Warningf("Failed to get %s: %s", model.ModelName(m), res.Error)
		return *new(M), model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": model.ModelName(m), "Error": res.Error.Error()}}
	}
	if res.RowsAffected == 0 {
		return *new(M), model.RetVal{Msg: i18n.Model.NotFound, Attr: map[string]any{"Model": model.ModelName(m)}}
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
	res := applyGetOptions(b.DB.Model(new(M)).Clauses(clause.Locking{Strength: "UPDATE"}), options).Limit(1).Find(&m)
	if res.Error != nil {
		log.Logger.Warningf("Failed to get %s for update: %s", model.ModelName(m), res.Error)
		return *new(M), model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": model.ModelName(m), "Error": res.Error.Error()}}
	}
	if res.RowsAffected == 0 {
		return *new(M), model.RetVal{Msg: i18n.Model.NotFound, Attr: map[string]any{"Model": model.ModelName(*new(M))}}
	}
	return m, model.SuccessRetVal()
}

func (b *BaseRepo[M]) GetByUniqueField(key string, value any, optionsL ...GetOptions) (M, model.RetVal) {
	if !slices.Contains(model.UniqueFields(*new(M)), key) {
		return *new(M), model.RetVal{Msg: i18n.Model.NotUniqueKey, Attr: map[string]any{"Model": model.ModelName(*new(M)), "Key": key}}
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
		log.Logger.Warningf("Failed to count %s: %s", model.ModelName(*new(M)), res.Error)
		return 0, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": model.ModelName(*new(M)), "Error": res.Error.Error()}}
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
	if res := tx.Limit(limit).Offset(offset).Find(&ms); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, count, model.RetVal{Msg: i18n.Model.NotFound, Attr: map[string]any{"Model": model.ModelName(*new(M))}}
		}
		log.Logger.Warningf("Failed to get %s: %s", model.ModelName(*new(M)), res.Error)
		return nil, count, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": model.ModelName(*new(M)), "Error": res.Error.Error()}}
	}
	return ms, count, model.SuccessRetVal()
}

func (b *BaseRepo[M]) Update(id uint, options UpdateOptions) model.RetVal {
	var count uint
	data := options.Convert2Map()
	if len(data) == 0 {
		return model.SuccessRetVal()
	}
	for {
		count++
		if count > 10 {
			log.Logger.Warningf("Failed to update %s: too many times failed due to optimistic lock", model.ModelName(*new(M)))
			return model.RetVal{Msg: i18n.Model.DeadLock, Attr: map[string]any{"Model": model.ModelName(*new(M))}}
		}
		m, ret := b.GetByID(id)
		if !ret.OK {
			return ret
		}
		res := b.DB.Model(&m).Where("id = ?", id).Updates(data)
		if res.Error != nil {
			if attr, ok := duplicateKeyAttr(res.Error); ok {
				attr["Model"] = model.ModelName(*new(M))
				return model.RetVal{Msg: i18n.Model.NotUniqueKey, Attr: attr}
			}
			log.Logger.Warningf("Failed to update %s: %s", model.ModelName(*new(M)), res.Error)
			return model.RetVal{Msg: i18n.Model.UpdateError, Attr: map[string]any{"Model": model.ModelName(*new(M)), "Error": res.Error.Error()}}
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
		log.Logger.Warningf("Failed to update %s: %s", model.ModelName(*new(M)), res.Error)
		return model.RetVal{Msg: i18n.Model.UpdateError, Attr: map[string]any{"Model": model.ModelName(*new(M)), "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func (b *BaseRepo[M]) Delete(idL ...uint) model.RetVal {
	res := b.DB.Model(new(M)).Where("id IN ?", idL).Delete(new(M))
	if res.Error != nil {
		log.Logger.Warningf("Failed to delete %s: %s", model.ModelName(*new(M)), res.Error)
		return model.RetVal{Msg: i18n.Model.DeleteError, Attr: map[string]any{"Model": model.ModelName(*new(M)), "Error": res.Error.Error()}}
	}
	return model.SuccessRetVal()
}

func duplicateKeyAttr(err error) (map[string]any, bool) {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		return nil, false
	}
	attr := map[string]any{"Error": pgErr.Error()}
	if pgErr.ColumnName != "" {
		attr["Key"] = pgErr.ColumnName
	}
	return attr, true
}
