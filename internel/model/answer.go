package model

type Answer struct {
	TeamID uint   `json:"team_id"`
	Team   Team   `json:"-"`
	FlagID uint   `json:"flag_id"`
	Flag   Flag   `json:"-"`
	Value  string `gorm:"not null" json:"value"`
	Solved bool   `json:"solved"`
	BaseModel
}
