package form

// CreateUsageForm 将题目添加至比赛
type CreateUsageForm struct {
	ChallengeID []string `form:"challenge_id" json:"challenge_id" binding:"required"`
}

// UpdateUsageForm for admin update usage info
type UpdateUsageForm struct {
	Hidden       *bool    `form:"hidden" json:"hidden"`
	Score        *float64 `form:"score" json:"score"`
	ScoreType    *uint    `form:"score_type" json:"score_type"`
	CurrentScore *float64 `form:"current_score" json:"current_score"`
	MinScore     *float64 `form:"min_score" json:"min_score"`
	Decay        *float64 `form:"decay" json:"decay"`
	Attempt      *int64   `form:"attempt" json:"attempt"`
	Hints        *string  `form:"hints" json:"hints"`
	Tags         *string  `form:"tags" json:"tags"`
}
