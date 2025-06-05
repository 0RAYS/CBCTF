package form

type CreateChallengeForm struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Desc     string `form:"desc" json:"desc"`
	Category string `form:"category" json:"category"`
	Type     string `form:"type" json:"type"`
}
