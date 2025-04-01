package model

type Usage struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	ContestID       uint            `json:"contest_id"`
	Contest         Contest         `json:"-"`
	ChallengeID     string          `json:"challenge_id"`
	Challenge       Challenge       `json:"-"`
	Name            string          `gorm:"not null" json:"name"`
	Desc            string          `json:"desc"`
	Hidden          bool            `json:"hidden"`
	Docker          Docker          `gorm:"type:json" json:"docker"`
	Dockers         Dockers         `gorm:"type:json" json:"dockers"`
	Hints           Strings         `gorm:"type:json" json:"hints"`
	Tags            Strings         `gorm:"type:json" json:"tags"`
	NetworkPolicies NetworkPolicies `gorm:"type:json" json:"network_policies"`
	Flags           []Flag          `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Containers      []Container     `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Submissions     []Submission    `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Version         uint            `gorm:"default:1" json:"-"`
}
