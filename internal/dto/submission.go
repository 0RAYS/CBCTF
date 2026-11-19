package dto

type SubmitFlagForm struct {
	Flag string `form:"flag" json:"flag" binding:"required"`
}
