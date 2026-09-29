package db

import (
	"gorm.io/gorm"

	"CBCTF/internal/model"
	"CBCTF/internal/oa"
)

type OauthRepo struct {
	BaseRepo[model.Oauth]
}

type UpdateOauthOptions struct {
	Protocol         *string
	AuthURL          *string
	TokenURL         *string
	UserInfoURL      *string
	CallbackURL      *string
	ClientID         *string
	ClientSecret     *string
	Provider         *string
	Uri              *string
	Scopes           *[]string
	IDClaim          *string
	NameClaim        *string
	EmailClaim       *string
	PictureClaim     *string
	DescriptionClaim *string
	GroupsClaim      *string
	AdminGroup       *string
	DefaultGroup     *uint
	On               *bool
	Picture          *model.FileURL
}

func (u UpdateOauthOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Protocol != nil {
		options["protocol"] = *u.Protocol
	}
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
	if u.Scopes != nil {
		options["scopes"] = *u.Scopes
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
	if u.GroupsClaim != nil {
		options["groups_claim"] = *u.GroupsClaim
	}
	if u.AdminGroup != nil {
		options["admin_group"] = *u.AdminGroup
	}
	if u.DefaultGroup != nil {
		options["default_group"] = *u.DefaultGroup
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
		DB: tx,
	}
}

func (o *OauthRepo) RegisterDefault() model.RetVal {
	count, ret := o.Count()
	if !ret.OK {
		return ret
	}
	if count > 0 {
		return model.SuccessRetVal()
	}
	return WithTransactionDB(o.DB, func(tx *gorm.DB) model.RetVal {
		repo := InitOauthRepo(tx)
		for _, provider := range []model.Oauth{oa.GetDefaultGithubOauth(), oa.GetDefaultHDUHelpOauth(), oa.GetDefaultHDUCASOauth()} {
			if _, ret := repo.Create(provider); !ret.OK {
				return ret
			}
		}
		return model.SuccessRetVal()
	})
}
