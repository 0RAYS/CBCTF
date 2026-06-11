package model

type UserTeam struct {
	UserID uint `gorm:"primaryKey;autoIncrement:false;index"`
	TeamID uint `gorm:"primaryKey;autoIncrement:false;index"`
}

type UserContest struct {
	UserID    uint `gorm:"primaryKey;autoIncrement:false;index"`
	ContestID uint `gorm:"primaryKey;autoIncrement:false;index"`
}

type UserGroup struct {
	UserID  uint `gorm:"primaryKey;autoIncrement:false;index"`
	GroupID uint `gorm:"primaryKey;autoIncrement:false;index"`
}

type RolePermission struct {
	RoleID       uint `gorm:"primaryKey;autoIncrement:false;index"`
	PermissionID uint `gorm:"primaryKey;autoIncrement:false;index"`
}
