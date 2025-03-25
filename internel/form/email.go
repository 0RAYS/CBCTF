package form

// VerifyEmail 邮箱验证表单
type VerifyEmail struct {
	ID    string `form:"id" binding:"required"`
	Token string `form:"token" binding:"required"`
}
