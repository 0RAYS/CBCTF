package model

import (
	"golang.org/x/oauth2"
)

type Oauth struct {
	AuthURL          string  `json:"auth_url"`
	TokenURL         string  `json:"token_url"`
	UserInfoURL      string  `json:"user_info_url"`
	CallbackURL      string  `json:"callback_url"`
	ClientID         string  `json:"client_id"`
	ClientSecret     string  `json:"client_secret"`
	Provider         string  `gorm:"type:varchar(255);uniqueIndex:idx_oauth_provider_active,where:deleted_at IS NULL;not null" json:"provider"`
	Uri              string  `json:"uri"`
	IDClaim          string  `json:"id_claim"`
	NameClaim        string  `json:"name_claim"`
	EmailClaim       string  `json:"email_claim"`
	PictureClaim     string  `json:"picture_claim"`
	DescriptionClaim string  `json:"description_claim"`
	GroupsClaim      string  `json:"groups_claim"`
	AdminGroup       string  `json:"admin_group"`
	DefaultGroup     uint    `json:"default_group"`
	On               bool    `json:"on"`
	Picture          FileURL `json:"picture"`
	BaseModel
}

func (o *Oauth) Config(scopes []string) *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     o.ClientID,
		ClientSecret: o.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  o.AuthURL,
			TokenURL: o.TokenURL,
		},
		RedirectURL: o.CallbackURL,
	}
	if len(scopes) > 0 {
		config.Scopes = scopes
	}
	return config
}
