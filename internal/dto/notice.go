package dto

// CreateNoticeForm 创建公告表单
type CreateNoticeForm struct {
	Title   string `form:"title" json:"title" binding:"required"`
	Content string `form:"content" json:"content" binding:"required"`
	Type    string `form:"type" json:"type" binding:"required,oneof=normal update important"`
}

// UpdateNoticeForm 更新公告表单
type UpdateNoticeForm struct {
	Title   *string `form:"title" json:"title" binding:"omitempty,min=1"`
	Content *string `form:"content" json:"content" binding:"omitempty,min=1"`
	Type    *string `form:"type" json:"type" binding:"omitempty,oneof=normal update important"`
}
