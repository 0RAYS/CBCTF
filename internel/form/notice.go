package form

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/model"
	"github.com/gin-gonic/gin"
	"slices"
)

var allowedNoticeType = []string{model.NoticeTypeNormal, model.NoticeTypeUpdate, model.NoticeTypeImportant}

// CreateNoticeForm 创建公告表单
type CreateNoticeForm struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Type    string `json:"type" binding:"required"`
}

func (f *CreateNoticeForm) Bind(c *gin.Context) (bool, string) {
	if err := c.ShouldBind(f); err != nil {
		return false, i18n.BadRequest
	}
	if !slices.Contains(allowedNoticeType, f.Type) {
		return false, i18n.InvalidNoticeType
	}
	return true, i18n.Success
}

// UpdateNoticeForm 更新公告表单
type UpdateNoticeForm struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Type    *string `json:"type"`
}

func (f *UpdateNoticeForm) Bind(c *gin.Context) (bool, string) {
	if err := c.ShouldBind(f); err != nil {
		return false, i18n.BadRequest
	}
	if f.Type != nil {
		if !slices.Contains(allowedNoticeType, *f.Type) {
			return false, i18n.InvalidNoticeType
		}
	}
	return true, i18n.Success
}
