package model

import "CBCTF/internal/i18n"

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
	BasicModel
}

func (s Submission) GetModelName() string {
	return "Submission"
}

func (s Submission) GetVersion() uint {
	return s.Version
}

func (s Submission) CreateErrorString() string {
	return i18n.CreateSubmissionError
}

func (s Submission) DeleteErrorString() string {
	return i18n.DeleteSubmissionError
}

func (s Submission) GetErrorString() string {
	return i18n.GetSubmissionError
}

func (s Submission) NotFoundErrorString() string {
	return i18n.SubmissionNotFound
}

func (s Submission) UpdateErrorString() string {
	return i18n.UpdateSubmissionError
}

func (s Submission) GetUniqueKey() []string {
	return []string{"id"}
}

func (s Submission) GetForeignKeys() []string {
	return []string{"id", "contest_challenge_id", "contest_id", "challenge_id", "team_id", "user_id", "contest_flag_id"}
}
