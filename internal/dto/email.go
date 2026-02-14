package dto

// VerifyEmail 邮箱验证表单
type VerifyEmail struct {
	ID    string `form:"id" json:"id" binding:"required"`
	Token string `form:"token" json:"token" binding:"required"`
}
