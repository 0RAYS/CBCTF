package service

import (
	"CBCTF/internal/dto"
	"CBCTF/internal/model"
	"testing"

	"gorm.io/gorm"
)

func TestGeneratorBatchesDoNotReportMissingTargetsAsSuccess(t *testing.T) {
	database := queryFailureDB(t)
	_ = database.Callback().Query().After("gorm:query").Register("test:missing", func(tx *gorm.DB) {
		if count, ok := tx.Statement.Dest.(*int64); ok {
			*count = 0
		}
	})
	for _, ret := range []model.RetVal{
		StartGenerators(database, 1, dto.StartGeneratorsForm{Challenges: []string{"missing", "missing"}}),
		StopGenerators(database, 1, dto.StopGeneratorsForm{Generators: []uint{1, 2}}),
	} {
		batch, ok := ret.Data.(*model.BatchResult)
		if ret.OK || !ok || batch.Requested != 2 || batch.Failed != 2 {
			t.Fatalf("missing targets hidden: %+v %+v", ret, batch)
		}
	}
}
