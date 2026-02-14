package dto

type GetAssetForm struct {
	Filename string `form:"filename" json:"filename" binding:"required"`
}
