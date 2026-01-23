package model

import "CBCTF/internal/i18n"

type Oauth struct {
	AuthURL         string    `json:"auth_url"`
	TokenURL        string    `json:"token_url"`
	UserInfoURL     string    `json:"user_info_url"`
	CallbackURL     string    `json:"callback_url"`
	ClientID        string    `json:"client_id"`
	ClientSecret    string    `json:"client_secret"`
	Provider        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"provider"`
	Uri             string    `json:"uri"`
	RespIDField     string    `json:"id_field"`
	RespNameField   string    `json:"name_field"`
	RespEmailField  string    `json:"email_field"`
	RespAvatarField string    `json:"avatar_field"`
	RespDescField   string    `json:"desc_field"`
	On              bool      `json:"on"`
	Avatar          AvatarURL `json:"avatar"`
	BaseModel
}

func (o Oauth) GetModelName() string {
	return "Oauth"
}

func (o Oauth) GetBaseModel() BaseModel {
	return o.BaseModel
}

func (o Oauth) CreateErrorString() string {
	return i18n.CreateOauthError
}

func (o Oauth) DeleteErrorString() string {
	return i18n.DeleteOauthError
}

func (o Oauth) GetErrorString() string {
	return i18n.GetOauthError
}

func (o Oauth) NotFoundErrorString() string {
	return i18n.OauthNotFound
}

func (o Oauth) UpdateErrorString() string {
	return i18n.UpdateOauthError
}

func (o Oauth) GetUniqueKey() []string {
	return []string{"id", "provider"}
}

func (o Oauth) GetAllowedQueryFields() []string {
	return []string{"id", "provider", "on"}
}
