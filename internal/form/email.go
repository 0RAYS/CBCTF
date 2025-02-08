package form

type VerifyEmail struct {
	ID    string `form:"id" binding:"required"`
	Token string `form:"token" binding:"required"`
}
