package model

// Admin 系统管理员
type Admin struct {
	Name     string  `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Password string  `gorm:"not null" json:"-"`
	Email    string  `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Picture  FileURL `json:"picture"`
	Verified bool    `gorm:"default:false" json:"verified"`
	BaseModel
}

func (a Admin) ModelName() string {
	return "Admin"
}

func (a Admin) GetBaseModel() BaseModel {
	return a.BaseModel
}

func (a Admin) UniqueFields() []string {
	return []string{"id", "name", "email"}
}

func (a Admin) QueryFields() []string {
	return []string{"id", "name", "email", "verified"}
}
