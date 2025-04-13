package model

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
	"Container": {
		"Pod",
	},
	"Cheat": {
		"User", "Team", "Contest",
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
	"Pod": {
		"Victim", "Containers", "Traffics",
	},
	"Request": {},
	"Submission": {
		"Usage", "Contest", "Challenge", "Team", "User", "Flag",
	},
	"Team": {
		"Contest", "Users", "Answers", "Submissions", "Victims", "Cheats",
	},
	"Traffic": {
		"Victim", "Pod",
	},
	"Usage": {
		"Contest", "Challenge", "Flags", "Victims", "Submissions",
	},
	"User": {
		"Teams", "Contests", "Submissions", "Victims", "Devices", "Cheats",
	},
	"Victim": {
		"Usage", "Team", "User", "Pods", "Traffics",
	},
}
