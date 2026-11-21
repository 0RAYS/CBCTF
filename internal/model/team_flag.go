package model

type TeamFlag struct {
	TeamID          uint          `gorm:"index;uniqueIndex:idx_team_flags_team_challenge_flag_active,where:deleted_at IS NULL;index:idx_team_flags_team_value_active,priority:1,where:deleted_at IS NULL;index:idx_team_flags_team_solved_challenge_flag_active,priority:1,where:deleted_at IS NULL" json:"team_id"`
	Team            Team          `json:"-"`
	ChallengeFlagID uint          `gorm:"index;uniqueIndex:idx_team_flags_team_challenge_flag_active,where:deleted_at IS NULL;index:idx_team_flags_team_solved_challenge_flag_active,priority:3,where:deleted_at IS NULL" json:"challenge_flag_id"`
	ChallengeFlag   ChallengeFlag `json:"-"`
	Value           string        `gorm:"index:idx_team_flags_team_value_active,priority:2,where:deleted_at IS NULL" json:"value"`
	Solved          bool          `gorm:"index:idx_team_flags_team_solved_challenge_flag_active,priority:2,where:deleted_at IS NULL" json:"solved"`
	BaseModel
}
