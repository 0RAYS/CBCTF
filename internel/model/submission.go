package model

import "CBCTF/internel/i18n"

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

func (r Submission) GetModelName() string {
	return "Submission"
}

func (r Submission) GetID() uint {
	return r.ID
}

func (r Submission) GetVersion() uint {
	return r.Version
}

func (r Submission) CreateErrorString() string {
	return i18n.CreateSubmissionError
}

func (r Submission) DeleteErrorString() string {
	return i18n.DeleteSubmissionError
}

func (r Submission) GetErrorString() string {
	return i18n.GetSubmissionError
}

func (r Submission) NotFoundErrorString() string {
	return i18n.SubmissionNotFound
}

func (r Submission) UpdateErrorString() string {
	return i18n.UpdateSubmissionError
}

func (r Submission) GetUniqueKey() []string {
	return []string{"id"}
}
