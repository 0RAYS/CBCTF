package model

type UserTeam struct {
	UserID uint `gorm:"index:idx_user_team,unique"`
	TeamID uint `gorm:"index:idx_user_team,unique"`
}

type UserContest struct {
	UserID    uint `gorm:"index:idx_user_contest,unique"`
	ContestID uint `gorm:"index:idx_user_contest,unique"`
}
