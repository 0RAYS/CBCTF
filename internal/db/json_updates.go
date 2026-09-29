package db

import (
	"database/sql/driver"
	"maps"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// JSONUpdates applies a model field's declared serializer to map updates too.
// GORM serializes struct writes, but Updates(map) binds raw values; slices would
// otherwise expand into SQL tuples and embedded structs reach the driver raw.
// A single adapter keeps all pools and all repositories on the model tags.
type JSONUpdates struct{}

func (JSONUpdates) Name() string { return "cbctf:json_updates" }

func (JSONUpdates) Initialize(database *gorm.DB) error {
	return database.Callback().Update().Before("gorm:update").Register("cbctf:json_updates", serializeJSONUpdates)
}

func serializeJSONUpdates(tx *gorm.DB) {
	changes, ok := tx.Statement.Dest.(map[string]any)
	if !ok || tx.Statement.Schema == nil || tx.Error != nil {
		return
	}
	var prepared map[string]any
	selected, restricted := tx.Statement.SelectAndOmitColumns(false, true)
	for key, value := range changes {
		field := tx.Statement.Schema.LookUpField(key)
		if field == nil || field.Serializer == nil || !field.Updatable {
			continue
		}
		if included, exists := selected[field.DBName]; (exists && !included) || (!exists && restricted) {
			continue
		}
		// Expressions are explicit SQL; driver-valuers include values already
		// wrapped by GORM for struct writes. Never encode either a second time.
		if _, ok = value.(clause.Expression); ok {
			continue
		}
		if _, ok = value.(driver.Valuer); ok {
			continue
		}
		record := reflect.New(tx.Statement.Schema.ModelType).Elem()
		if err := field.Set(tx.Statement.Context, record, value); err != nil {
			_ = tx.AddError(err)
			return
		}
		bound, _ := field.ValueOf(tx.Statement.Context, record)
		if prepared == nil {
			prepared = maps.Clone(changes)
		}
		prepared[key] = bound
	}
	if prepared != nil {
		tx.Statement.Dest = prepared
	}
}
