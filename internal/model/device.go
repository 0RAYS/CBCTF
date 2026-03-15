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

func (d Device) TableName() string {
	return "devices"
}

func (d Device) ModelName() string {
	return "Device"
}

func (d Device) GetBaseModel() BaseModel {
	return d.BaseModel
}

func (d Device) UniqueFields() []string {
	return []string{"id"}
}

func (d Device) QueryFields() []string {
	return []string{"id", "user_id", "magic"}
}
