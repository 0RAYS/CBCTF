package dto

type CreateNoticeForm struct {
	Title   string `form:"title" json:"title" binding:"required"`
	Content string `form:"content" json:"content" binding:"required"`
	Type    string `form:"type" json:"type" binding:"required"`
}

type UpdateNoticeForm struct {
	Title   *string `form:"title" json:"title" binding:"omitempty,min=1"`
	Content *string `form:"content" json:"content" binding:"omitempty,min=1"`
	Type    *string `form:"type" json:"type" binding:"omitempty,min=1"`
}
