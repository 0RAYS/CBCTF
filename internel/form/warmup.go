package form

type WarmUpImageForm struct {
	Images     []string `form:"images" json:"images" binding:"required"`
	PullPolicy string   `form:"pull_policy" json:"pull_policy" binding:"required"`
}
