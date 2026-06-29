package dto

// VerifyEmail 邮箱验证表单
type VerifyEmail struct {
	Token string `form:"token" json:"token" binding:"required"`
}
