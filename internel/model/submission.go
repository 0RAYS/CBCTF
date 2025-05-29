package model

type Submission struct {
	UsageID     uint      `json:"usage_id"`
	Usage       Usage     `json:"-"`
	ContestID   uint      `json:"contest_id"`
	Contest     Contest   `json:"-"`
	ChallengeID string    `json:"challenge_id"`
	Challenge   Challenge `json:"-"`
	TeamID      uint      `json:"team_id"`
	Team        Team      `json:"-"`
	UserID      uint      `json:"user_id"`
	User        User      `json:"-"`
	FlagID      uint      `json:"flag_id"`
	Flag        Flag      `json:"-"`
	Value       string    `json:"value"`
	Solved      bool      `json:"solved"`
	Score       float64   `gorm:"default:0" json:"score"`
	BaseModel
}
