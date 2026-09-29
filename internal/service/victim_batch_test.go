package service

import (
	"testing"

	"gorm.io/gorm"

	"CBCTF/internal/dto"
	"CBCTF/internal/model"
)

func TestStopVictimsReportsNonStoppableMissingAndIdempotentItems(t *testing.T) {
	database := queryFailureDB(t)
	_ = database.Callback().Query().After("gorm:query").Register("test:victims", func(tx *gorm.DB) {
		switch dst := tx.Statement.Dest.(type) {
		case *int64:
			*dst = 2
		case *[]model.Victim:
			*dst = []model.Victim{{ID: 1, Status: model.WaitingVictimStatus}, {ID: 2, Status: model.TerminatingVictimStatus}}
		}
		tx.RowsAffected = 2
	})
	ret := StopVictims(database, 1, dto.StopVictimsForm{Victims: []uint{1, 2, 3}})
	batch, ok := ret.Data.(*model.BatchResult)
	if ret.OK || !ok || batch.Failed != 2 || batch.Skipped != 1 || batch.NotAttempted != 0 || batch.Requested != 3 {
		t.Fatalf("batch failures hidden: %+v %+v", ret, batch)
	}
}
