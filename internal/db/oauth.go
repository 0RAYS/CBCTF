package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/oauth"

	"gorm.io/gorm"
)

type OauthRepo struct {
	BaseRepo[model.Oauth]
}

type CreateOauthOptions struct {
	AuthURL              string
	TokenURL             string
	UserInfoURL          string
	CallbackURL          string
	ClientID             string
	ClientSecret         string
	Provider             string
	Uri                  string
	RespIDField          string
	RespNameField        string
	RespEmailField       string
	RespPictureField     string
	RespDescriptionField string
	On                   bool
	Picture              model.FileURL
}

func (c CreateOauthOptions) Convert2Model() model.Model {
	return model.Oauth{
		AuthURL:              c.AuthURL,
		TokenURL:             c.TokenURL,
		UserInfoURL:          c.UserInfoURL,
		CallbackURL:          c.CallbackURL,
		ClientID:             c.ClientID,
		ClientSecret:         c.ClientSecret,
		Provider:             c.Provider,
		Uri:                  c.Uri,
		RespIDField:          c.RespIDField,
		RespNameField:        c.RespNameField,
		RespEmailField:       c.RespEmailField,
		RespPictureField:     c.RespPictureField,
		RespDescriptionField: c.RespDescriptionField,
		On:                   c.On,
		Picture:              c.Picture,
	}
}

type UpdateOauthOptions struct {
	AuthURL              *string
	TokenURL             *string
	UserInfoURL          *string
	CallbackURL          *string
	ClientID             *string
	ClientSecret         *string
	Provider             *string
	Uri                  *string
	RespIDField          *string
	RespNameField        *string
	RespEmailField       *string
	RespPictureField     *string
	RespDescriptionField *string
	On                   *bool
	Picture              *model.FileURL
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
		options["user_info_url"] = *u.UserInfoURL
	}
	if u.CallbackURL != nil {
		options["callback_url"] = *u.CallbackURL
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
	if u.Uri != nil {
		options["uri"] = *u.Uri
	}
	if u.RespIDField != nil {
		options["resp_id_field"] = *u.RespIDField
	}
	if u.RespNameField != nil {
		options["resp_name_field"] = *u.RespNameField
	}
	if u.RespEmailField != nil {
		options["resp_email_field"] = *u.RespEmailField
	}
	if u.RespPictureField != nil {
		options["resp_picture_field"] = *u.RespPictureField
	}
	if u.RespDescriptionField != nil {
		options["resp_description_field"] = *u.RespDescriptionField
	}
	if u.On != nil {
		options["on"] = *u.On
	}
	if u.Picture != nil {
		options["picture"] = *u.Picture
	}
	return options
}

func InitOauthRepo(tx *gorm.DB) *OauthRepo {
	return &OauthRepo{
		BaseRepo: BaseRepo[model.Oauth]{
			DB: tx,
		},
	}
}

func (o *OauthRepo) RegisterDefault() {
	github := oauth.GetDefaultGithubOauth()
	_, ret := o.GetByUniqueKey("provider", github.Provider, GetOptions{Selects: []string{"id"}})
	if !ret.OK {
		if err := o.DB.Model(&model.Oauth{}).Create(&github).Error; err != nil {
			log.Logger.Warningf("Failed to register default github oauth provider: %s", err)
		}
	}
	hduhelp := oauth.GetDefaultHDUHelpOauth()
	_, ret = o.GetByUniqueKey("provider", hduhelp.Provider, GetOptions{Selects: []string{"id"}})
	if !ret.OK {
		if err := o.DB.Model(&model.Oauth{}).Create(&hduhelp).Error; err != nil {
			log.Logger.Warningf("Failed to register default hduhelp oauth provider: %s", err)
		}
	}
}
