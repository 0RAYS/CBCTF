package model

type UserTeam struct {
	UserID uint `gorm:"primaryKey;autoIncrement:false"`
	TeamID uint `gorm:"primaryKey;autoIncrement:false"`
}

type UserContest struct {
	UserID    uint `gorm:"primaryKey;autoIncrement:false"`
	ContestID uint `gorm:"primaryKey;autoIncrement:false"`
}

type UserGroup struct {
	UserID  uint `gorm:"primaryKey;autoIncrement:false"`
	GroupID uint `gorm:"primaryKey;autoIncrement:false"`
}

type RolePermission struct {
	RoleID       uint `gorm:"primaryKey;autoIncrement:false"`
	PermissionID uint `gorm:"primaryKey;autoIncrement:false"`
}
