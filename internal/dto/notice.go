package dto

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

// CreateNoticeForm 创建公告表单
type CreateNoticeForm struct {
	Title   string `form:"title" json:"title" binding:"required"`
	Content string `form:"content" json:"content" binding:"required"`
	Type    string `form:"type" json:"type" binding:"required,oneof=normal update important"`
}

func (f *CreateNoticeForm) Bind(c *gin.Context) model.RetVal {
	if err := c.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

// UpdateNoticeForm 更新公告表单
type UpdateNoticeForm struct {
	Title   *string `form:"title" json:"title" binding:"omitempty,min=1"`
	Content *string `form:"content" json:"content" binding:"omitempty,min=1"`
	Type    *string `form:"type" json:"type" binding:"omitempty,oneof=normal update important"`
}

func (f *UpdateNoticeForm) Bind(c *gin.Context) model.RetVal {
	if err := c.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
