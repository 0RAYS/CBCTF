package db

import (
	"strings"
	"testing"

	"CBCTF/internal/traffic"
	"gorm.io/gorm"
)

func TestPartialTrafficSnapshotCannotRetractDatabaseEvidence(t *testing.T) {
	database := jsonTestDB(t)
	var sql string
	_ = database.Callback().Create().After("gorm:create").Register("test:partial_traffic", func(tx *gorm.DB) { sql = tx.Statement.SQL.String() })
	ret := InitTrafficRepo(database).ReplaceAnalysis(1, []string{}, []traffic.TrafficAccess{}, &traffic.AnalysisReport{Partial: true}, false)
	if !ret.OK || !strings.Contains(sql, "jsonb_array_elements") || !strings.Contains(sql, "traffics.accesses") || !strings.Contains(sql, "EXCLUDED.accesses") {
		t.Fatalf("partial scan replaces trusted observations: %s %+v", sql, ret)
	}
}
