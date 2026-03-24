package model

// TeamFlag
// BelongsTo Team
// BelongsTo ContestFlag
// BelongsTo ChallengeFlag
type TeamFlag struct {
	TeamID          uint          `gorm:"index;uniqueIndex:idx_team_flags_team_contest_flag_active,where:deleted_at IS NULL" json:"team_id"`
	Team            Team          `json:"-"`
	ContestFlagID   uint          `gorm:"index;uniqueIndex:idx_team_flags_team_contest_flag_active,where:deleted_at IS NULL" json:"contest_flag_id"`
	ContestFlag     ContestFlag   `json:"-"`
	ChallengeFlagID uint          `gorm:"index" json:"challenge_flag_id"`
	ChallengeFlag   ChallengeFlag `json:"-"`
	Value           string        `json:"value"`
	Solved          bool          `json:"solved"`
	BaseModel
}
