package task

import (
	"CBCTF/internal/model"
	"testing"
)

func TestTaskRecordRetryPreservesIdentityAndNewArrivals(t *testing.T) {
	previous := DrainTaskRecordPool()
	t.Cleanup(func() { DrainTaskRecordPool(); RestoreTaskRecords(previous) })
	appendTaskRecord(model.Task{TaskID: "first"})
	batch := DrainTaskRecordPool()
	if len(batch) != 1 || batch[0].RecordID == "" {
		t.Fatal("missing stable identity")
	}
	appendTaskRecord(model.Task{TaskID: "second"})
	RestoreTaskRecords(batch)
	got := DrainTaskRecordPool()
	if len(got) != 2 || got[0].RecordID != batch[0].RecordID || got[1].TaskID != "second" {
		t.Fatalf("retry lost records: %+v", got)
	}
}
