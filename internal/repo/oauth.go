package repo

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/oauth"
	"gorm.io/gorm"
)

type OauthRepo struct {
	BasicRepo[model.Oauth]
}

type CreateOauthOptions struct {
	AuthURL         string
	TokenURL        string
	UserInfoURL     string
	ClientID        string
	ClientSecret    string
	Provider        string
	RedirectURI     string
	RespNameField   string
	RespEmailField  string
	RespAvatarField string
	RespDescField   string
	On              bool
	Avatar          model.AvatarURL
}

func (c CreateOauthOptions) Convert2Model() model.Oauth {
	return model.Oauth{
		AuthURL:         c.AuthURL,
		TokenURL:        c.TokenURL,
		UserInfoURL:     c.UserInfoURL,
		ClientID:        c.ClientID,
		ClientSecret:    c.ClientSecret,
		Provider:        c.Provider,
		RedirectURI:     c.RedirectURI,
		RespNameField:   c.RespNameField,
		RespEmailField:  c.RespEmailField,
		RespAvatarField: c.RespAvatarField,
		RespDescField:   c.RespDescField,
		On:              c.On,
		Avatar:          c.Avatar,
	}
}

type UpdateOauthOptions struct {
	AuthURL         *string
	TokenURL        *string
	UserInfoURL     *string
	ClientID        *string
	ClientSecret    *string
	Provider        *string
	RedirectURI     *string
	RespNameField   *string
	RespEmailField  *string
	RespAvatarField *string
	RespDescField   *string
	On              *bool
	Avatar          *model.AvatarURL
}

func (u UpdateOauthOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.AuthURL != nil {
		options["auth_url"] = *u.AuthURL
	}
	if u.TokenURL != nil {
		options["token_url"] = *u.TokenURL
	}
	if u.UserInfoURL != nil {
		options["userinfo_url"] = *u.UserInfoURL
	}
	if u.ClientID != nil {
		options["client_id"] = *u.ClientID
	}
	if u.ClientSecret != nil {
		options["client_secret"] = *u.ClientSecret
	}
	if u.Provider != nil {
		options["provider"] = *u.Provider
	}
	if u.RedirectURI != nil {
		options["redirect_uri"] = *u.RedirectURI
	}
	if u.RespNameField != nil {
		options["name_field"] = *u.RespNameField
	}
	if u.RespEmailField != nil {
		options["email_field"] = *u.RespEmailField
	}
	if u.RespAvatarField != nil {
		options["avatar_field"] = *u.RespAvatarField
	}
	if u.RespDescField != nil {
		options["desc_field"] = *u.RespDescField
	}
	if u.On != nil {
		options["on"] = *u.On
	}
	if u.Avatar != nil {
		options["avatar"] = *u.Avatar
	}
	return options
}

func InitOauthRepo(tx *gorm.DB) *OauthRepo {
	return &OauthRepo{
		BasicRepo: BasicRepo[model.Oauth]{
			DB: tx,
		},
	}
}

func (o *OauthRepo) RegisterDefault() {
	github := oauth.GetDefaultGithubOauth()
	_, ok, _ := o.GetByUniqueKey("provider", github.Provider, GetOptions{Selects: []string{"id"}})
	if !ok {
		if err := o.DB.Model(&model.Oauth{}).Create(&github).Error; err != nil {
			log.Logger.Warningf("Failed to register default github oauth provider: %v", err)
		}
	}
	hduhelp := oauth.GetDefaultHDUHelpOauth()
	_, ok, _ = o.GetByUniqueKey("provider", hduhelp.Provider, GetOptions{Selects: []string{"id"}})
	if !ok {
		if err := o.DB.Model(&model.Oauth{}).Create(&hduhelp).Error; err != nil {
			log.Logger.Warningf("Failed to register default hduhelp oauth provider: %v", err)
		}
	}
}
