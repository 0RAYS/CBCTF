package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
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
		"Usages",
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
		"Answers",
	},
	"Notice": {
		"Contest", "Admin",
	},
	"Request": {},
	"Submission": {
		"Usage", "Contest", "Challenge", "Team", "User",
	},
	"Team": {
		"Contest", "Users", "Answers", "Submissions", "Containers", "Cheats",
	},
	"Traffic": {
		"Container",
	},
	"Usage": {
		"Contest", "Challenge", "Flags", "Containers",
	},
	"User": {
		"Teams", "Contests", "Submissions", "Containers", "Devices", "Cheats",
	},
}

func generateAssociations(key string, depth int) []string {
	tmp := strings.TrimSuffix(key, "s")
	if depth < 2 {
		return []string{}
	}
	var result []string
	if associations, exists := Associations[tmp]; exists {
		for _, assoc := range associations {
			fullAssoc := key + "." + assoc
			result = append(result, fullAssoc)
			if depth > 2 {
				subAssociations := generateAssociations(assoc, depth-1)
				for _, sub := range subAssociations {
					result = append(result, fullAssoc+"."+sub[len(assoc)+1:])
				}
			}
		}
	}
	return result
}

func GetPreload(tx *gorm.DB, model string, preload bool, depth int) *gorm.DB {
	if preload {
		tx = tx.Preload(clause.Associations)
		if depth < 2 {
			return tx
		}
		depth++
		result := generateAssociations(model, depth)
		for _, r := range result {
			t := strings.Split(r, ".")[1:]
			if len(t) < 2 {
				continue
			}
			tx = tx.Preload(strings.Join(t, "."))
		}
	}
	return tx
}
