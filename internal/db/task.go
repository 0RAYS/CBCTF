package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"sort"

	"gorm.io/gorm"
)

type TaskRepo struct {
	BaseRepo[model.Task]
}

func InitTaskRepo(tx *gorm.DB) *TaskRepo {
	return &TaskRepo{
		BaseRepo: BaseRepo[model.Task]{
			DB: tx,
		},
	}
}

// Create stores task execution history without generic uniqueness checks.
// Task records are append-only observability data and have no natural unique key.
func (t *TaskRepo) Create(task model.Task) model.RetVal {
	if res := t.DB.Model(&model.Task{}).Create(&task); res.Error != nil {
		log.Logger.Warningf("Failed to create Task: %s", res.Error)
		return model.RetVal{Msg: i18n.Model.CreateError, Attr: map[string]any{"Model": model.Name(model.Task{}), "Error": res.Error.Error()}}
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
		return nil, model.RetVal{Msg: i18n.Model.Task.GetError, Attr: map[string]any{"Error": err.Error()}}
	}
	sort.Strings(queues)
	return queues, model.SuccessRetVal()
}
