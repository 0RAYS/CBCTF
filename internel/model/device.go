package model

type Device struct {
	UserID uint   `json:"user_id"`
	User   User   `json:"-"`
	Magic  string `json:"magic"`
	Count  int    `json:"count"`
	BaseModel
}
