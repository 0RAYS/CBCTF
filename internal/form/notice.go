package form

type CreateNoticeForm struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateNoticeForm struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}
