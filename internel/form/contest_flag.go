package form

type UpdateContestFlagForm struct {
	Value     *string  `form:"value" json:"value"`
	Score     *float64 `form:"score" json:"score"`
	Decay     *float64 `form:"decay" json:"decay"`
	MinScore  *float64 `form:"min_score" json:"min_score"`
	ScoreType *uint    `form:"score_type" json:"score_type"`
}
