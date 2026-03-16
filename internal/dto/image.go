package dto

type PullImageForm struct {
	Images     []string `form:"images" json:"images" binding:"required,dive,min=1"`
	PullPolicy string   `form:"pull_policy" json:"pull_policy" binding:"required,oneof=Always Never IfNotPresent"`
}
