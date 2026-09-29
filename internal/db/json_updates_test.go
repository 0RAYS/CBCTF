package db

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"CBCTF/internal/model"
)

func jsonTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	database, err := gorm.Open(postgres.Open("host=127.0.0.1 dbname=test"), &gorm.Config{DryRun: true, DisableAutomaticPing: true, SkipDefaultTransaction: true, Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	pool, _ := database.DB()
	t.Cleanup(func() { _ = pool.Close() })
	if err = database.Use(JSONUpdates{}); err != nil {
		t.Fatal(err)
	}
	return database
}

func resolvedJSONArgs(t *testing.T, tx *gorm.DB) []string {
	t.Helper()
	if tx.Error != nil {
		t.Fatal(tx.Error)
	}
	var values []string
	for _, value := range tx.Statement.Vars {
		if valuer, ok := value.(driver.Valuer); ok {
			var err error
			value, err = valuer.Value()
			if err != nil {
				t.Fatal(err)
			}
		}
		if text, ok := value.(string); ok && json.Valid([]byte(text)) {
			values = append(values, text)
		}
	}
	return values
}

func TestSerializerMapUpdatesMatchStructWritesAndPreserveInput(t *testing.T) {
	database := jsonTestDB(t)
	for _, rules := range [][]string{{"one", "中文"}, {}, nil} {
		changes := map[string]any{"rules": rules}
		mapped := database.Model(&model.Contest{}).Where("id = ?", 1).Updates(changes)
		structured := database.Model(&model.Contest{}).Where("id = ?", 1).Select("rules").Updates(&model.Contest{Rules: rules})
		a, b := resolvedJSONArgs(t, mapped), resolvedJSONArgs(t, structured)
		if !reflect.DeepEqual(a, b) {
			t.Fatalf("map=%v struct=%v", a, b)
		}
		if !reflect.DeepEqual(changes["rules"], rules) {
			t.Fatal("map mutated; optimistic retries would double-encode")
		}
		if strings.Contains(mapped.Statement.SQL.String(), "($") {
			t.Fatalf("JSON list expanded to SQL tuple: %s", mapped.Statement.SQL.String())
		}
	}
	for _, tc := range []struct {
		model any
		key   string
		value any
		want  string
	}{
		{&model.Webhook{}, "headers", map[string]string{}, `{}`},
		{&model.Branding{}, "site_name", model.LocalizedText{ZhCN: "平台", En: "Platform"}, `{"zh_cn":"平台","en":"Platform"}`},
		{&model.ChallengeFlag{}, "binding", model.FlagBinding{PodKey: "web", Target: "FLAG"}, `{"pod_key":"web","container_key":"","type":"","target":"FLAG"}`},
	} {
		args := resolvedJSONArgs(t, database.Model(tc.model).Where("id = ?", 1).Updates(map[string]any{tc.key: tc.value}))
		if len(args) != 1 || args[0] != tc.want {
			t.Fatalf("%s: %v", tc.key, args)
		}
	}
}

func TestSerializerExpressionsErrorsAndScanReset(t *testing.T) {
	database := jsonTestDB(t)
	expression := gorm.Expr("COALESCE(rules, '[]'::jsonb) || ?::jsonb", `["new"]`)
	query := database.Model(&model.Contest{}).Where("id = ?", 1).Update("rules", expression)
	if query.Error != nil || !strings.Contains(query.Statement.SQL.String(), "COALESCE(rules") {
		t.Fatalf("expression was serialized: %s %v", query.Statement.SQL.String(), query.Error)
	}
	query = database.Model(&model.Contest{}).Where("id = ?", 1).Update("rules", 42)
	if query.Error == nil {
		t.Fatal("invalid JSON field type accepted")
	}
	query = database.Model(&model.Contest{}).Where("id = ?", 1).Select("name").Updates(map[string]any{"name": "new", "rules": 42})
	if query.Error != nil {
		t.Fatalf("omitted serializer field was evaluated: %v", query.Error)
	}
	parsed, err := schema.Parse(&model.Contest{}, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		t.Fatal(err)
	}
	field := parsed.LookUpField("rules")
	record := model.Contest{Rules: []string{"stale"}}
	dst := reflect.ValueOf(&record).Elem()
	for _, input := range []any{`["one","中文"]`, []byte(`[]`), nil} {
		if err = field.Serializer.Scan(context.Background(), field, dst, input); err != nil {
			t.Fatal(err)
		}
		if input == nil && record.Rules != nil {
			t.Fatal("SQL NULL did not clear reused model")
		}
	}
	if err = field.Serializer.Scan(context.Background(), field, dst, `broken`); err == nil {
		t.Fatal("malformed JSON silently accepted")
	}
}

func TestVictimNestedJSONCreateAndUpdate(t *testing.T) {
	database := jsonTestDB(t)
	spec := model.VictimSpec{FrpEnabled: true, Pods: []model.PodSpec{{Key: "web", Containers: []model.VictimContainerSpec{{Name: "app", Environment: map[string]string{"FLAG": "flag{test}"}, Command: []string{"sh", "-c", "sleep 1"}}}}}}
	want, _ := json.Marshal(spec)
	for _, tx := range []*gorm.DB{
		database.Create(&model.Victim{Spec: spec}),
		database.Model(&model.Victim{}).Where("id = ?", 1).Updates(map[string]any{"spec": spec}),
	} {
		found := false
		for _, arg := range resolvedJSONArgs(t, tx) {
			if arg == string(want) {
				found = true
			}
		}
		if !found {
			t.Fatalf("nested spec changed: %v", resolvedJSONArgs(t, tx))
		}
	}
}

func TestSettingUpdatePreservesScalarZeroValues(t *testing.T) {
	database := jsonTestDB(t)
	if err := database.Callback().Query().After("gorm:query").Register("test:setting", func(tx *gorm.DB) {
		setting := tx.Statement.Dest.(*model.Setting)
		setting.ID = 1
		setting.Key = "test"
		tx.RowsAffected = 1
	}); err != nil {
		t.Fatal(err)
	}
	if err := database.Callback().Update().After("gorm:update").Register("test:updated", func(tx *gorm.DB) { tx.RowsAffected = 1 }); err != nil {
		t.Fatal(err)
	}
	for _, value := range []any{false, 0, "", []string{}, map[string]string{}} {
		if ret := InitSettingRepo(database).Update("test", UpdateSettingOptions{Value: &model.SettingValue{V: value}}); !ret.OK {
			t.Fatalf("scalar update failed: %#v %+v", value, ret)
		}
	}
}

type jsonRoundTripRow struct {
	ID     uint              `gorm:"primaryKey"`
	List   []string          `gorm:"serializer:json;type:jsonb;default:'[]'"`
	Object map[string]string `gorm:"serializer:json;type:jsonb;default:'{}'"`
	Spec   model.VictimSpec  `gorm:"serializer:json;type:jsonb;default:'{}'"`
}

func TestJSONSerializerPostgresRoundTrip(t *testing.T) {
	dsn := os.Getenv("CBCTF_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("CBCTF_TEST_POSTGRES_DSN is not set")
	}
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	pool, _ := database.DB()
	defer pool.Close()
	if err = database.Use(JSONUpdates{}); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tx := database.WithContext(ctx).Begin()
	defer tx.Rollback()
	if err = tx.Exec(`CREATE TEMP TABLE serializer_roundtrip (id bigint PRIMARY KEY, list jsonb DEFAULT '[]', object jsonb DEFAULT '{}', spec jsonb DEFAULT '{}') ON COMMIT DROP`).Error; err != nil {
		t.Fatal(err)
	}
	row := jsonRoundTripRow{ID: 1, List: []string{"first"}, Object: map[string]string{"a": "b"}, Spec: model.VictimSpec{FrpEnabled: true}}
	if err = tx.Table("serializer_roundtrip").Create(&row).Error; err != nil {
		t.Fatal(err)
	}
	if err = tx.Table("serializer_roundtrip").Model(&row).Updates(map[string]any{"list": []string{}, "object": map[string]string{}, "spec": model.VictimSpec{}}).Error; err != nil {
		t.Fatal(err)
	}
	var read jsonRoundTripRow
	if err = tx.Table("serializer_roundtrip").First(&read, 1).Error; err != nil {
		t.Fatal(err)
	}
	if read.List == nil || len(read.List) != 0 || read.Object == nil || len(read.Object) != 0 || read.Spec.FrpEnabled {
		t.Fatalf("empty values did not round-trip: %+v", read)
	}
	if err = tx.Table("serializer_roundtrip").Model(&row).Update("list", nil).Error; err != nil {
		t.Fatal(err)
	}
	if err = tx.Table("serializer_roundtrip").First(&read, 1).Error; err != nil {
		t.Fatal(err)
	}
	if read.List != nil {
		t.Fatal("SQL NULL did not clear the old slice")
	}
	if err = tx.Table("serializer_roundtrip").Create(&jsonRoundTripRow{ID: 2}).Error; err != nil {
		t.Fatal(err)
	}
	read = jsonRoundTripRow{}
	if err = tx.Table("serializer_roundtrip").First(&read, 2).Error; err != nil {
		t.Fatal(err)
	}
	if read.List == nil || read.Object == nil {
		t.Fatalf("JSON defaults missing: %+v", read)
	}
}
