package model

type Admin struct {
	Notices  []Notice `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name     string   `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Password string   `gorm:"not null" json:"-"`
	Email    string   `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Avatar   string   `json:"avatar"`
	Verified bool     `gorm:"default:false" json:"verified"`
	BaseModel
}
