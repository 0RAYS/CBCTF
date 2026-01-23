package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"slices"

	"github.com/gin-gonic/gin"
)

var allowedNoticeType = []string{model.NoticeTypeNormal, model.NoticeTypeUpdate, model.NoticeTypeImportant}

// CreateNoticeForm 创建公告表单
type CreateNoticeForm struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Type    string `json:"type" binding:"required"`
}

func (f *CreateNoticeForm) Bind(c *gin.Context) model.RetVal {
	if err := c.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	if !slices.Contains(allowedNoticeType, f.Type) {
		return model.RetVal{Msg: i18n.Model.Notice.InvalidType}
	}
	return model.SuccessRetVal()
}

// UpdateNoticeForm 更新公告表单
type UpdateNoticeForm struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Type    *string `json:"type"`
}

func (f *UpdateNoticeForm) Bind(c *gin.Context) model.RetVal {
	if err := c.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	if f.Type != nil {
		if !slices.Contains(allowedNoticeType, *f.Type) {
			return model.RetVal{Msg: i18n.Model.Notice.InvalidType}
		}
	}
	return model.SuccessRetVal()
}
