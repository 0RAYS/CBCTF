package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/task"

	"github.com/gin-gonic/gin"
)

func GetTasks(ctx *gin.Context) {
	var form dto.ListTasksForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
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
	repo := db.InitTaskRepo(db.DB)
	tasks, count, ret := repo.List(form.Limit, form.Offset, options)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0, len(tasks))
	for _, item := range tasks {
		data = append(data, resp.GetTaskResp(item))
	}
	queues, ret := repo.ListQueues()
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"tasks": data, "count": count, "queues": queues}))
}

func GetLiveTasks(ctx *gin.Context) {
	var form dto.ListLiveTasksForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	live, count, queues, types, ret := task.ListLiveTasks(form.Status, form.Queue, form.TaskID, form.Limit, form.Offset)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0, len(live))
	for _, item := range live {
		data = append(data, resp.GetLiveTaskResp(item))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{
		"tasks":  data,
		"count":  count,
		"queues": queues,
		"types":  types,
	}))
}
