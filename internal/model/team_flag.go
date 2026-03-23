package model

// TeamFlag
// BelongsTo Team
// BelongsTo ContestFlag
// BelongsTo ChallengeFlag
type TeamFlag struct {
	TeamID          uint          `json:"team_id"`
	Team            Team          `json:"-"`
	ContestFlagID   uint          `json:"contest_flag_id"`
	ContestFlag     ContestFlag   `json:"-"`
	ChallengeFlagID uint          `json:"challenge_flag_id"`
	ChallengeFlag   ChallengeFlag `json:"-"`
	Value           string        `json:"value"`
	Solved          bool          `json:"solved"`
	BaseModel
}
