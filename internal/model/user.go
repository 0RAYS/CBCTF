package model

import "CBCTF/internal/i18n"

// User
// ManyToMany Teams
// ManyToMany Contests
// HasMany Devices
// HasMany Submissions
type User struct {
	Teams       []*Team      `gorm:"many2many:user_teams;" json:"-"`
	Contests    []*Contest   `gorm:"many2many:user_contests;" json:"-"`
	Devices     []Device     `json:"-"`
	Submissions []Submission `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name        string       `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Password    string       `gorm:"not null" json:"-"`
	Email       string       `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Country     string       `gorm:"default:'CN'" json:"country"`
	Avatar      AvatarURL    `json:"avatar"`
	Desc        string       `json:"desc"`
	Verified    bool         `gorm:"default:false" json:"verified"`
	Hidden      bool         `gorm:"default:false" json:"hidden"`
	Banned      bool         `gorm:"default:false" json:"banned"`
	Score       float64      `gorm:"default:0" json:"score"`
	Solved      int64        `gorm:"default:0" json:"solved"`
	BasicModel
}

func (u User) GetModelName() string {
	return "User"
}

func (u User) GetVersion() uint {
	return u.Version
}

func (u User) CreateErrorString() string {
	return i18n.CreateUserError
}

func (u User) DeleteErrorString() string {
	return i18n.DeleteUserError
}

func (u User) GetErrorString() string {
	return i18n.GetUserError
}

func (u User) NotFoundErrorString() string {
	return i18n.UserNotFound
}

func (u User) UpdateErrorString() string {
	return i18n.UpdateUserError
}

func (u User) GetUniqueKey() []string {
	return []string{"id", "name", "email"}
}
