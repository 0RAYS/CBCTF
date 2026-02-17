package model

type UserTeam struct {
	UserID uint `gorm:"index:idx_user_team,unique"`
	TeamID uint `gorm:"index:idx_user_team,unique"`
}

type UserContest struct {
	UserID    uint `gorm:"index:idx_user_contest,unique"`
	ContestID uint `gorm:"index:idx_user_contest,unique"`
}

type UserGroup struct {
	UserID  uint `gorm:"index:idx_user_group,unique"`
	GroupID uint `gorm:"index:idx_user_group,unique"`
}

type RolePermission struct {
	RoleID       uint `gorm:"index:idx_role_permission,unique"`
	PermissionID uint `gorm:"index:idx_role_permission,unique"`
}
