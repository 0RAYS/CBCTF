package model

import "CBCTF/internal/i18n"

type Oauth struct {
	AuthURL         string    `json:"auth_url"`
	TokenURL        string    `json:"token_url"`
	UserInfoURL     string    `json:"userinfo_url"`
	ClientID        string    `json:"client_id"`
	ClientSecret    string    `json:"client_secret"`
	Provider        string    `json:"provider"`
	RedirectURI     string    `json:"redirect_uri"`
	RespNameField   string    `json:"name_field"`
	RespEmailField  string    `json:"email_field"`
	RespAvatarField string    `json:"avatar_field"`
	RespDescField   string    `json:"desc_field"`
	On              bool      `json:"on"`
	Avatar          AvatarURL `gorm:"type:json" json:"avatar"`
	BasicModel
}

func (o Oauth) GetModelName() string {
	return "Oauth"
}

func (o Oauth) GetVersion() uint {
	return o.Version
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
	return []string{"id"}
}
