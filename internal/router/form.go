package router

import "time"

type LoginForm struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterForm struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type CreateUserForm RegisterForm

type CreateTeamForm struct {
	Name      string `json:"name" binding:"required"`
	ContestID uint   `json:"contest" binding:"required"`
	UserID    uint   `json:"user" binding:"required"`
}

type CreateContestForm struct {
	Name string `json:"name" binding:"required"`
}

type GetUsersForm struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

type GetTeamsForm struct {
	Offset    int  `form:"offset"`
	Limit     int  `form:"limit"`
	ContestID uint `form:"contest"`
}

type GetContestsForm GetUsersForm

type GetFilesForm GetUsersForm

// UpdateUserForm 使用指针类型主要为了判断是否赋值
type UpdateUserForm struct {
	Name     *string `json:"name"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
	Website  *string `json:"website"`
	Country  *string `json:"country"`
	Type     *string `json:"type"`
	Verified *bool   `json:"verified"`
}

// UpdateTeamForm 使用指针类型主要为了判断是否赋值
type UpdateTeamForm struct {
	Name    *string `json:"name"`
	Desc    *string `json:"desc"`
	Captcha *string `json:"captcha"`
	Banned  *bool   `json:"banned"`
	Hidden  *bool   `json:"hidden"`
}

// UpdateContestForm 使用指针类型主要为了判断是否赋值
type UpdateContestForm struct {
	Name     *string    `json:"name"`
	Desc     *string    `json:"desc"`
	Captcha  *string    `json:"captcha"`
	Size     *uint      `json:"size"`
	Start    *time.Time `json:"start"`
	Duration *uint64    `json:"duration"`
}
