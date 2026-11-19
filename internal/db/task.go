package db

import (
	"fmt"
	"sort"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
)

type TaskRepo struct {
	BaseRepo[model.Task]
}

func InitTaskRepo(tx *gorm.DB) *TaskRepo {
	return &TaskRepo{
		DB: tx,
	}
}

func (t *TaskRepo) CreateBatch(tasks ...model.Task) model.RetVal {
	if len(tasks) == 0 {
		return model.SuccessRetVal()
	}
	for i := range tasks {
		if tasks[i].RecordID == "" {
			tasks[i].RecordID = utils.UUID()
		}
	}
	if res := t.DB.Model(&model.Task{}).Omit("id").Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "record_id"}}, DoNothing: true}).CreateInBatches(tasks, 200); res.Error != nil {
		log.Logger.Warningf("Failed to create Tasks: %s", res.Error)
		return model.FailRetVal(i18n.Model.CreateError, model.Name(model.Task{}), fmt.Sprintf("count=%d", len(tasks)), res.Error)
	}
	return model.SuccessRetVal()
}

func (t *TaskRepo) ListQueues() ([]string, model.RetVal) {
	queues := make([]string, 0)
	if err := t.DB.Model(&model.Task{}).
		Distinct("queue").
		Where("queue <> ''").
		Order("queue ASC").
		Pluck("queue", &queues).Error; err != nil {
		return nil, model.FailRetVal(i18n.Model.Task.GetError, "Task", "", err)
	}
	sort.Strings(queues)
	return queues, model.SuccessRetVal()
}
