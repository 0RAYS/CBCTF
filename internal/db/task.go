package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"sort"
	"time"

	"gorm.io/gorm"
)

type TaskRepo struct {
	BaseRepo[model.Task]
}

type CreateTaskOptions struct {
	TaskID      string
	Type        string
	Queue       string
	Status      string
	Payload     any
	Result      any
	Error       string
	RetryCount  int
	MaxRetry    int
	ProcessedAt time.Time
}

func (c CreateTaskOptions) Convert2Model() model.Model {
	return model.Task{
		TaskID:      c.TaskID,
		Type:        c.Type,
		Queue:       c.Queue,
		Status:      c.Status,
		Payload:     model.TaskPayload{V: c.Payload},
		Result:      model.TaskPayload{V: c.Result},
		Error:       c.Error,
		RetryCount:  c.RetryCount,
		MaxRetry:    c.MaxRetry,
		ProcessedAt: c.ProcessedAt,
	}
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
		return nil, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": model.ModelName(model.Task{}), "Error": err.Error()}}
	}
	sort.Strings(queues)
	return queues, model.SuccessRetVal()
}

func (t *TaskRepo) ListTypes() ([]string, model.RetVal) {
	types := make([]string, 0)
	if err := t.DB.Model(&model.Task{}).
		Distinct("type").
		Where("type <> ''").
		Order("type ASC").
		Pluck("type", &types).Error; err != nil {
		return nil, model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": model.ModelName(model.Task{}), "Error": err.Error()}}
	}
	sort.Strings(types)
	return types, model.SuccessRetVal()
}
