package model

import "CBCTF/internal/i18n"

// ContestChallenge 被赛事引用的挑战
// BelongsTo Contest
// BelongsTo Challenge
// HasMany ContestFlag
// HasMany Submission
type ContestChallenge struct {
	ContestID    uint          `gorm:"index:idx_contest_challenge_deleted_salt,unique;" json:"contest_id"`
	Contest      Contest       `json:"-"`
	ChallengeID  uint          `gorm:"index:idx_contest_challenge_deleted_salt,unique;" json:"challenge_id"`
	Challenge    Challenge     `json:"-"`
	ContestFlags []ContestFlag `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Submissions  []Submission  `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name         string        `json:"name"`
	Desc         string        `json:"desc"`
	Type         string        `json:"type"`
	Category     string        `json:"category"`
	Hidden       bool          `json:"hidden"`
	Attempt      int64         `json:"attempt"`
	Hints        StringList    `gorm:"default:null;type:json" json:"hints"`
	Tags         StringList    `gorm:"default:null;type:json" json:"tags"`
	DeletedSalt  string        `gorm:"default:'';type:varchar(36);index:idx_contest_challenge_deleted_salt,unique;" json:"-"`
	BaseModel
}

func (c ContestChallenge) GetModelName() string {
	return "ContestChallenge"
}

func (c ContestChallenge) GetBaseModel() BaseModel {
	return c.BaseModel
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

func (c ContestChallenge) GetAllowedQueryFields() []string {
	return []string{}
}
