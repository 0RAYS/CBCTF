package model

const OauthLoginType = "oauth_login"

// Device
// BelongsTo User
// HasMany Request
type Device struct {
	UserID uint   `json:"user_id"`
	User   User   `json:"-"`
	Magic  string `json:"magic"`
	Count  int    `json:"count"`
	BaseModel
}

func (d Device) GetModelName() string {
	return "Device"
}

func (d Device) GetBaseModel() BaseModel {
	return d.BaseModel
}

func (d Device) GetUniqueKey() []string {
	return []string{"id"}
}

func (d Device) GetAllowedQueryFields() []string {
	return []string{}
}
