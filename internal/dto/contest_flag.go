package dto

type UpdateContestFlagForm struct {
	Value     *string  `form:"value" json:"value" binding:"omitempty,min=1"`
	Score     *float64 `form:"score" json:"score" binding:"omitempty,gte=0"`
	Decay     *float64 `form:"decay" json:"decay" binding:"omitempty,gte=0"`
	MinScore  *float64 `form:"min_score" json:"min_score" binding:"omitempty,gte=0"`
	ScoreType *uint    `form:"score_type" json:"score_type" binding:"omitempty,oneof=0 1 2"`
}
