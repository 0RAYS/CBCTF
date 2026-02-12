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
	AuthURL          string
	TokenURL         string
	UserInfoURL      string
	CallbackURL      string
	ClientID         string
	ClientSecret     string
	Provider         string
	Uri              string
	IDField          string
	NameField        string
	EmailField       string
	PictureField     string
	DescriptionField string
	On               bool
	Picture          model.FileURL
}

func (c CreateOauthOptions) Convert2Model() model.Model {
	return model.Oauth{
		AuthURL:          c.AuthURL,
		TokenURL:         c.TokenURL,
		UserInfoURL:      c.UserInfoURL,
		CallbackURL:      c.CallbackURL,
		ClientID:         c.ClientID,
		ClientSecret:     c.ClientSecret,
		Provider:         c.Provider,
		Uri:              c.Uri,
		IDField:          c.IDField,
		NameField:        c.NameField,
		EmailField:       c.EmailField,
		PictureField:     c.PictureField,
		DescriptionField: c.DescriptionField,
		On:               c.On,
		Picture:          c.Picture,
	}
}

type UpdateOauthOptions struct {
	AuthURL          *string
	TokenURL         *string
	UserInfoURL      *string
	CallbackURL      *string
	ClientID         *string
	ClientSecret     *string
	Provider         *string
	Uri              *string
	IDField          *string
	NameField        *string
	EmailField       *string
	PictureField     *string
	DescriptionField *string
	On               *bool
	Picture          *model.FileURL
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
	if u.IDField != nil {
		options["id_field"] = *u.IDField
	}
	if u.NameField != nil {
		options["name_field"] = *u.NameField
	}
	if u.EmailField != nil {
		options["email_field"] = *u.EmailField
	}
	if u.PictureField != nil {
		options["picture_field"] = *u.PictureField
	}
	if u.DescriptionField != nil {
		options["description_field"] = *u.DescriptionField
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
