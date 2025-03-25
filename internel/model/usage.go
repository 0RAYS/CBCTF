package model

type Usage struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	ContestID   uint        `json:"contest_id"`
	Contest     Contest     `json:"-"`
	ChallengeID string      `json:"challenge_id"`
	Challenge   Challenge   `json:"-"`
	Name        string      `gorm:"not null" json:"name"`
	Desc        string      `json:"desc"`
	Flags       []Flag      `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Containers  []Container `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Version     uint        `gorm:"default:1" json:"-"`
}
