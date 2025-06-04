package form

// CreateNoticeForm 创建公告表单
type CreateNoticeForm struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	Type    string `json:"type" binding:"required"`
}

// UpdateNoticeForm 更新公告表单
type UpdateNoticeForm struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Type    *string `json:"type"`
}
