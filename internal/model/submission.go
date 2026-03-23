package model

// Submission 提交记录
// BelongsTo ContestChallenge
// BelongsTo Contest
// BelongsTo Challenge
// BelongsTo Team
// BelongsTo User
// BelongsTo ContestFlag
type Submission struct {
	ContestChallengeID uint             `gorm:"index:idx_submissions_team_challenge_solved,priority:2;index" json:"contest_challenge_id"`
	ContestChallenge   ContestChallenge `json:"-"`
	ContestID          uint             `gorm:"index" json:"contest_id"`
	Contest            Contest          `json:"-"`
	ChallengeID        uint             `gorm:"index" json:"challenge_id"`
	Challenge          Challenge        `json:"-"`
	TeamID             uint             `gorm:"index:idx_submissions_team_challenge_solved,priority:1;index" json:"team_id"`
	Team               Team             `json:"-"`
	UserID             uint             `gorm:"index" json:"user_id"`
	User               User             `json:"-"`
	ContestFlagID      uint             `gorm:"index:idx_submissions_contest_flag_solved,priority:1;index" json:"contest_flag_id"`
	ContestFlag        ContestFlag      `json:"-"`
	Value              string           `json:"value"`
	Solved             bool             `gorm:"index:idx_submissions_team_challenge_solved,priority:3;index:idx_submissions_contest_flag_solved,priority:2" json:"solved"`
	Score              float64          `gorm:"default:0" json:"score"`
	IP                 string           `json:"ip"`
	BaseModel
}
