package model

import "CBCTF/internel/i18n"

// ContestChallenge 被赛事引用的挑战
// BelongsTo Contest
// BelongsTo Challenge
// HasMany ContestFlag
type ContestChallenge struct {
	ContestID    uint          `gorm:"index:idx_contest_challenge,unique;" json:"contest_id"`
	Contest      Contest       `json:"-"`
	ChallengeID  uint          `json:"index:idx_contest_challenge,unique;challenge_id"`
	Challenge    Challenge     `json:"-"`
	ContestFlags []ContestFlag `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Type         string        `json:"type"`
	Name         string        `json:"name"`
	Desc         string        `json:"desc"`
	Hidden       bool          `json:"hidden"`
	Attempt      int64         `json:"attempt"`
	Basic
}

func (c ContestChallenge) GetModelName() string {
	return "ContestChallenge"
}

func (c ContestChallenge) GetID() uint {
	return c.ID
}

func (c ContestChallenge) GetVersion() uint {
	return c.Version
}

func (c ContestChallenge) CreateErrorString() string {
	return i18n.CreateContestChallengeError
}

func (c ContestChallenge) DeleteErrorString() string {
	return i18n.DeleteContestChallengeError
}

func (c ContestChallenge) GetErrorString() string {
	return i18n.GetContestChallengeError
}

func (c ContestChallenge) NotFoundErrorString() string {
	return i18n.ContestChallengeNotFound
}

func (c ContestChallenge) UpdateErrorString() string {
	return i18n.UpdateContestChallengeError
}

func (c ContestChallenge) GetUniqueKey() []string {
	return []string{"id"}
}
