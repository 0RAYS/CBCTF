package dto

type PullImageTarget struct {
	Node   string `form:"node" json:"node" binding:"required,min=1"`
	Image  string `form:"image" json:"image" binding:"required,min=1"`
	Manual bool   `form:"manual" json:"manual"`
}

type PullImageForm struct {
	Targets    []PullImageTarget `form:"targets" json:"targets" binding:"required,dive"`
	PullPolicy string            `form:"pull_policy" json:"pull_policy" binding:"required,oneof=Always Never IfNotPresent"`
}
