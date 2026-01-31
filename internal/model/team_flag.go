package model

// TeamFlag
// BelongsTo Team
// BelongsTo ContestFlag
type TeamFlag struct {
	TeamID          uint        `json:"team_id"`
	ContestFlagID   uint        `json:"contest_flag_id"`
	ContestFlag     ContestFlag `json:"-"`
	ChallengeFlagID uint        `json:"challenge_flag_id"`
	Value           string      `json:"value"`
	Solved          bool        `json:"solved"`
	BaseModel
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
