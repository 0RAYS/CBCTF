package dto

import (
	"CBCTF/internal/model"
)

type CreateWebhookForm struct {
	Name    string           `form:"name" json:"name" binding:"required"`
	URL     string           `form:"url" json:"url" binding:"required,url"`
	Method  string           `form:"method" json:"method" binding:"required,oneof=POST GET"`
	Headers model.StringMap  `form:"headers" json:"headers"`
	Timeout int64            `form:"timeout" json:"timeout" binding:"gte=0"`
	Retry   int              `form:"retry" json:"retry" binding:"gte=0"`
	Events  model.StringList `form:"events" json:"events"`
}

type UpdateWebhookForm struct {
	Name    *string           `form:"name" json:"name" binding:"omitempty,min=1"`
	URL     *string           `form:"url" json:"url" binding:"omitempty,url"`
	Method  *string           `form:"method" json:"method" binding:"omitempty,oneof=POST GET"`
	Headers *model.StringMap  `form:"headers" json:"headers"`
	Timeout *int64            `form:"timeout" json:"timeout" binding:"omitempty,gte=0"`
	Retry   *int              `form:"retry" json:"retry" binding:"omitempty,gte=0"`
	On      *bool             `form:"on" json:"on"`
	Events  *model.StringList `form:"events" json:"events"`
}
