package db

import (
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
	IDClaim          string
	NameClaim        string
	EmailClaim       string
	PictureClaim     string
	DescriptionClaim string
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
		IDClaim:          c.IDClaim,
		NameClaim:        c.NameClaim,
		EmailClaim:       c.EmailClaim,
		PictureClaim:     c.PictureClaim,
		DescriptionClaim: c.DescriptionClaim,
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
	IDClaim          *string
	NameClaim        *string
	EmailClaim       *string
	PictureClaim     *string
	DescriptionClaim *string
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
	if u.IDClaim != nil {
		options["id_claim"] = *u.IDClaim
	}
	if u.NameClaim != nil {
		options["name_claim"] = *u.NameClaim
	}
	if u.EmailClaim != nil {
		options["email_claim"] = *u.EmailClaim
	}
	if u.PictureClaim != nil {
		options["picture_claim"] = *u.PictureClaim
	}
	if u.DescriptionClaim != nil {
		options["description_claim"] = *u.DescriptionClaim
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
	_, ret := o.GetByUniqueKey("provider", github.Provider)
	if !ret.OK {
		o.Insert(github)
	}
	hduhelp := oauth.GetDefaultHDUHelpOauth()
	_, ret = o.GetByUniqueKey("provider", hduhelp.Provider)
	if !ret.OK {
		o.Insert(hduhelp)
	}
}
