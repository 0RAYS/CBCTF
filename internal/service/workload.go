package service

import (
	"fmt"

	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
)

// stopWorkloads stops the workloads with the given ids. The workloads that are
// already stopping are skipped and the ids that match nothing are reported.
func stopWorkloads[T model.Model](
	tx *gorm.DB,
	kind string,
	ids []uint,
	contestID uint,
	list func(db.GetOptions) ([]T, int64, model.RetVal),
	status func(T) string,
	stop func(T) model.RetVal,
) model.RetVal {
	if len(ids) == 0 {
		return model.SuccessRetVal()
	}
	options := db.GetOptions{Conditions: map[string]any{"id": ids}, Deleted: true}
	if contestID > 0 {
		options.Conditions["contest_id"] = contestID
	}
	workloads, _, ret := list(options)
	if !ret.OK {
		return ret
	}
	byID := make(map[uint]T, len(workloads))
	for _, workload := range workloads {
		byID[workload.GetBaseModel().ID] = workload
	}
	batch := model.NewBatch(len(ids))
	for _, id := range ids {
		if tx.Statement.Context.Err() != nil {
			return batch.Result(tx.Statement.Context)
		}
		key := fmt.Sprint(id)
		workload, exists := byID[id]
		if !exists {
			batch.Fail(key, "lookup", model.FailRetVal(i18n.Model.NotFound, kind, "id="+key, nil))
			continue
		}
		if model.IsStopping(status(workload)) {
			batch.Skip(key, "already_stopping")
			continue
		}
		ret = stop(workload)
		if ret.OK {
			batch.Success(key, "queued")
		} else {
			batch.Fail(key, "stop", ret)
			if model.BatchDependencyFailed(ret) {
				return batch.Result(tx.Statement.Context)
			}
		}
	}
	return batch.Result(tx.Statement.Context)
}
