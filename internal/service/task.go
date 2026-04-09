package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

func ListTasks(tx *gorm.DB, form dto.ListTasksForm) ([]model.Task, int64, []string, model.RetVal) {
	options := db.GetOptions{
		Sort: []string{"processed_at DESC", "id DESC"},
	}
	conditions := make(map[string]any)
	if form.Queue != "" {
		conditions["queue"] = form.Queue
	}
	if form.Status != "" {
		conditions["status"] = form.Status
	}
	if len(conditions) > 0 {
		options.Conditions = conditions
	}
	if form.TaskID != "" {
		options.Search = map[string]string{"task_id": form.TaskID}
	}
	repo := db.InitTaskRepo(tx)
	tasks, count, ret := repo.List(form.Limit, form.Offset, options)
	if !ret.OK {
		return nil, 0, nil, ret
	}
	queues, ret := repo.ListQueues()
	if !ret.OK {
		return nil, 0, nil, ret
	}
	return tasks, count, queues, model.SuccessRetVal()
}
