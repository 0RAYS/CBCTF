package dto

type StartGeneratorsForm struct {
	Challenges []string `form:"challenges" json:"challenges" binding:"required,dive,uuid"`
}

type StopGeneratorsForm struct {
	Generators []uint `form:"generators" json:"generators" binding:"required,dive,gt=0"`
}
