package dto

type ListLogForm struct {
	ListModelsForm
	Level string `form:"level" json:"level" binding:"omitempty,oneof=TRACE DEBUG INFO WARNING ERROR FATAL PANIC trace debug info warning error fatal panic"`
}
