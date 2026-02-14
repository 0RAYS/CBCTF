package dto

// SubmitFlagForm for submit flag
type SubmitFlagForm struct {
	Flag string `form:"flag" json:"flag" binding:"required"`
}
