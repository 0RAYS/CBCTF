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

func (t TeamFlag) TableName() string {
	return "team_flags"
}

func (t TeamFlag) ModelName() string {
	return "TeamFlag"
}

func (t TeamFlag) GetBaseModel() BaseModel {
	return t.BaseModel
}

func (t TeamFlag) UniqueFields() []string {
	return []string{"id"}
}

func (t TeamFlag) QueryFields() []string {
	return []string{}
}
