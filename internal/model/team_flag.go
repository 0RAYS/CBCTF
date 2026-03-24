package model

// TeamFlag
// BelongsTo Team
// BelongsTo ContestFlag
// BelongsTo ChallengeFlag
type TeamFlag struct {
	TeamID          uint          `gorm:"index;uniqueIndex:idx_team_flags_team_contest_flag_active,where:deleted_at IS NULL;index:idx_team_flags_team_value_active,priority:1,where:deleted_at IS NULL;index:idx_team_flags_team_solved_contest_flag_active,priority:1,where:deleted_at IS NULL" json:"team_id"`
	Team            Team          `json:"-"`
	ContestFlagID   uint          `gorm:"index;uniqueIndex:idx_team_flags_team_contest_flag_active,where:deleted_at IS NULL;index:idx_team_flags_team_solved_contest_flag_active,priority:3,where:deleted_at IS NULL" json:"contest_flag_id"`
	ContestFlag     ContestFlag   `json:"-"`
	ChallengeFlagID uint          `gorm:"index" json:"challenge_flag_id"`
	ChallengeFlag   ChallengeFlag `json:"-"`
	Value           string        `gorm:"index:idx_team_flags_team_value_active,priority:2,where:deleted_at IS NULL" json:"value"`
	Solved          bool          `gorm:"index:idx_team_flags_team_solved_contest_flag_active,priority:2,where:deleted_at IS NULL" json:"solved"`
	BaseModel
}
