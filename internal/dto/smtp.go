package dto

type CreateSmtpForm struct {
	Address string `form:"address" json:"address" binding:"required,email"`
	Host    string `form:"host" json:"host" binding:"required,hostname"`
	Port    int    `form:"port" json:"port" binding:"required,gte=0,lte=65535"`
	Pwd     string `form:"pwd" json:"pwd" binding:"required"`
}

type UpdateSmtpForm struct {
	Address *string `form:"address" json:"address" binding:"omitempty,email"`
	Host    *string `form:"host" json:"host" binding:"omitempty,hostname"`
	Port    *int    `form:"port" json:"port" binding:"omitempty,gte=0,lte=65535"`
	Pwd     *string `form:"pwd" json:"pwd"`
	On      *bool   `form:"on" json:"on"`
}
