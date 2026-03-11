package model

// Submission 提交记录
// BelongsTo ContestChallenge
// BelongsTo Contest
// BelongsTo Challenge
// BelongsTo Team
// BelongsTo User
// BelongsTo ContestFlag
type Submission struct {
	ContestChallengeID uint             `json:"contest_challenge_id"`
	ContestChallenge   ContestChallenge `json:"-"`
	ContestID          uint             `json:"contest_id"`
	Contest            Contest          `json:"-"`
	ChallengeID        uint             `json:"challenge_id"`
	Challenge          Challenge        `json:"-"`
	TeamID             uint             `json:"team_id"`
	Team               Team             `json:"-"`
	UserID             uint             `json:"user_id"`
	User               User             `json:"-"`
	ContestFlagID      uint             `json:"contest_flag_id"`
	ContestFlag        ContestFlag      `json:"-"`
	Value              string           `json:"value"`
	Solved             bool             `json:"solved"`
	Score              float64          `gorm:"default:0" json:"score"`
	IP                 string           `json:"ip"`
	BaseModel
}

func (s Submission) TableName() string {
	return "submissions"
}

func (s Submission) ModelName() string {
	return "Submission"
}

func (s Submission) GetBaseModel() BaseModel {
	return s.BaseModel
}

func (s Submission) UniqueFields() []string {
	return []string{"id"}
}

func (s Submission) QueryFields() []string {
	return []string{"id", "value", "solved", "ip", "user_id", "team_id", "contest_id", "challenge_id", "contest_challenge_id"}
}
