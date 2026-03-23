package model

// ContestChallenge 被赛事引用的挑战
// BelongsTo Contest
// BelongsTo Challenge
// HasMany ContestFlag
// HasMany Submission
type ContestChallenge struct {
	ContestID    uint          `gorm:"index;uniqueIndex:idx_contest_challenges_unique_active,where:deleted_at IS NULL" json:"contest_id"`
	Contest      Contest       `json:"-"`
	ChallengeID  uint          `gorm:"index;uniqueIndex:idx_contest_challenges_unique_active,where:deleted_at IS NULL" json:"challenge_id"`
	Challenge    Challenge     `json:"-"`
	ContestFlags []ContestFlag `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Submissions  []Submission  `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Type         ChallengeType `gorm:"index" json:"type"`
	Category     string        `gorm:"index" json:"category"`
	Hidden       bool          `gorm:"index" json:"hidden"`
	Attempt      int64         `json:"attempt"`
	Hints        StringList    `gorm:"default:null;type:jsonb" json:"hints"`
	Tags         StringList    `gorm:"default:null;type:jsonb" json:"tags"`
	BaseModel
}
