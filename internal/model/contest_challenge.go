package model

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
	Description  string        `json:"description"`
	Type         ChallengeType `json:"type"`
	Category     string        `json:"category"`
	Hidden       bool          `json:"hidden"`
	Attempt      int64         `json:"attempt"`
	Hints        StringList    `gorm:"default:null;type:json" json:"hints"`
	Tags         StringList    `gorm:"default:null;type:json" json:"tags"`
	DeletedSalt  string        `gorm:"default:'';type:varchar(36);index:idx_contest_challenge_deleted_salt,unique;" json:"-"`
	BaseModel
}

func (c ContestChallenge) TableName() string {
	return "contest_challenges"
}

func (c ContestChallenge) ModelName() string {
	return "ContestChallenge"
}

func (c ContestChallenge) GetBaseModel() BaseModel {
	return c.BaseModel
}

func (c ContestChallenge) UniqueFields() []string {
	return []string{"id"}
}

func (c ContestChallenge) QueryFields() []string {
	return []string{"id", "contest_id", "challenge_id", "name", "category", "type", "hidden"}
}
