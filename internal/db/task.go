package db

import (
	"CBCTF/internal/i18n"
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
