package dto

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

type ListTasksForm struct {
	ListModelsForm
	Queue  string `form:"queue" json:"queue"`
	Status string `form:"status" json:"status"`
	TaskID string `form:"task_id" json:"task_id"`
}

func (f *ListTasksForm) Validate(ctx *gin.Context) model.RetVal {
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 20
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
	return model.SuccessRetVal()
}

type ListLiveTasksForm struct {
	ListModelsForm
	Queue  string `form:"queue" json:"queue"`
	Status string `form:"status" json:"status"`
	TaskID string `form:"task_id" json:"task_id"`
}

func (f *ListLiveTasksForm) Validate(ctx *gin.Context) model.RetVal {
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 20
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
	if f.Status == "" {
		f.Status = "active"
	}
	return model.SuccessRetVal()
}
