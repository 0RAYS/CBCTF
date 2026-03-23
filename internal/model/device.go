package model

const OauthLoginDeviceMagic = "oauth_login"

// Device
// BelongsTo User
// HasMany Request
type Device struct {
	UserID uint   `json:"user_id"`
	Magic  string `json:"magic"`
	Count  int64  `gorm:"default:0" json:"count"`
	BaseModel
}
