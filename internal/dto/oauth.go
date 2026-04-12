package dto

import (
	"CBCTF/internal/model"
)

type OauthCallbackForm struct {
	Code  string `form:"code" json:"code" binding:"required"`
	State string `form:"state" json:"state" binding:"required"`
}

type CreateOauthProviderForm struct {
	AuthURL          string `form:"auth_url" json:"auth_url" binding:"required,url"`
	TokenURL         string `form:"token_url" json:"token_url" binding:"required,url"`
	UserInfoURL      string `form:"user_info_url" json:"user_info_url" binding:"required,url"`
	CallbackURL      string `form:"callback_url" json:"callback_url" binding:"required,url"`
	ClientID         string `form:"client_id" json:"client_id" binding:"required"`
	ClientSecret     string `form:"client_secret" json:"client_secret" binding:"required"`
	Provider         string `form:"provider" json:"provider" binding:"required"`
	Uri              string `form:"uri" json:"uri" binding:"required,alphanum"`
	Scopes           model.StringList `form:"scopes" json:"scopes"`
	IDClaim          string `form:"id_claim" json:"id_claim" binding:"required"`
	NameClaim        string `form:"name_claim" json:"name_claim" binding:"required"`
	EmailClaim       string `form:"email_claim" json:"email_claim" binding:"required"`
	PictureClaim     string `form:"picture_claim" json:"picture_claim"`
	DescriptionClaim string `form:"description_claim" json:"description_claim"`
	GroupsClaim      string `form:"groups_claim" json:"groups_claim"`
	AdminGroup       string `form:"admin_group" json:"admin_group"`
	DefaultGroup     uint   `form:"default_group" json:"default_group"`
}

type UpdateOauthProviderForm struct {
	AuthURL          *string        `form:"auth_url" json:"auth_url" binding:"omitempty,url"`
	TokenURL         *string        `form:"token_url" json:"token_url" binding:"omitempty,url"`
	UserInfoURL      *string        `form:"user_info_url" json:"user_info_url" binding:"omitempty,url"`
	CallbackURL      *string        `form:"callback_url" json:"callback_url" binding:"omitempty,url"`
	ClientID         *string        `form:"client_id" json:"client_id" binding:"omitempty,min=1"`
	ClientSecret     *string        `form:"client_secret" json:"client_secret" binding:"omitempty,min=1"`
	Provider         *string        `form:"provider" json:"provider" binding:"omitempty,min=1"`
	Uri              *string        `form:"uri" json:"uri" binding:"omitempty,min=1,alphanum"`
	Scopes           *model.StringList `form:"scopes" json:"scopes"`
	IDClaim          *string        `form:"id_claim" json:"id_claim" binding:"omitempty,min=1"`
	NameClaim        *string        `form:"name_claim" json:"name_claim" binding:"omitempty,min=1"`
	EmailClaim       *string        `form:"email_claim" json:"email_claim" binding:"omitempty,min=1"`
	PictureClaim     *string        `form:"picture_claim" json:"picture_claim"`
	DescriptionClaim *string        `form:"description_claim" json:"description_claim"`
	GroupsClaim      *string        `form:"groups_claim" json:"groups_claim"`
	AdminGroup       *string        `form:"admin_group" json:"admin_group"`
	DefaultGroup     *uint          `form:"default_group" json:"default_group"`
	On               *bool          `form:"on" json:"on"`
	Picture          *model.FileURL `form:"picture" json:"picture"`
}
