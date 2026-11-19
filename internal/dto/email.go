package dto

type VerifyEmail struct {
	Token string `form:"token" json:"token" binding:"required"`
}
