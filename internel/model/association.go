package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserTeam struct {
	UserID uint `gorm:"index:idx_user_team,unique"`
	TeamID uint `gorm:"index:idx_user_team,unique"`
}

type UserContest struct {
	UserID    uint `gorm:"index:idx_user_contest,unique"`
	ContestID uint `gorm:"index:idx_user_contest,unique"`
}

var Associations = map[string][]string{
	"Admin": {
		"Notices",
	},
	"Answer": {
		"Team", "Flag",
	},
	"Challenge": {
		"Usages", "Submissions",
	},
	"Cheat": {
		"User", "Team", "Contest",
	},
	"Container": {
		"Usage", "Team", "User", "Traffics",
	},
	"Contest": {
		"Teams", "Users", "Notices", "Usages", "Flags", "Cheats", "Submissions",
	},
	"Device": {
		"User",
	},
	"File": {},
	"Flag": {
		"Contest", "Usage", "Answers", "Submissions",
	},
	"Notice": {
		"Contest", "Admin",
	},
	"Request": {},
	"Submission": {
		"Usage", "Contest", "Challenge", "Team", "User", "Flag",
	},
	"Team": {
		"Contest", "Users", "Answers", "Submissions", "Containers", "Cheats",
	},
	"Traffic": {
		"Container",
	},
	"Usage": {
		"Contest", "Challenge", "Flags", "Containers", "Submissions",
	},
	"User": {
		"Teams", "Contests", "Submissions", "Containers", "Devices", "Cheats",
	},
}

func GetPreload(tx *gorm.DB, preload bool, nestedL ...string) *gorm.DB {
	if preload {
		tx = tx.Preload(clause.Associations)
		for _, nested := range nestedL {
			tx = tx.Preload(nested)
		}
	}
	return tx
}
