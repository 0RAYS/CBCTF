package dto

type GetTrafficForm struct {
	TimeShift int64 `form:"time_shift" json:"time_shift" binding:"min=0,max=9223372036854"`
	Duration  int64 `form:"duration" json:"duration" binding:"min=1,max=9223372036854"`
}
