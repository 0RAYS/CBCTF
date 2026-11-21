package model

// ScoreAdjustment 是一条管理员对队伍分数的人工变动记录。
// Score 是带符号的变动分数，队伍显示分数等于解题分数加上全部变动分数之和。
type ScoreAdjustment struct {
	BaseModel
	ContestID uint    `gorm:"index" json:"contest_id"`
	TeamID    uint    `gorm:"index" json:"team_id"`
	Score     float64 `json:"score"`
	Reason    string  `json:"reason"`
	CreatorID uint    `gorm:"index" json:"creator_id"`
}
