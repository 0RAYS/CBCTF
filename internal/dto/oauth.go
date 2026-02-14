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
	IDField          string `form:"id_field" json:"id_field" binding:"required"`
	NameField        string `form:"name_field" json:"name_field" binding:"required"`
	EmailField       string `form:"email_field" json:"email_field" binding:"required"`
	PictureField     string `form:"picture_field" json:"picture_field"`
	DescriptionField string `form:"description_field" json:"description_field"`
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
	IDField          *string        `form:"id_field" json:"id_field" binding:"omitempty,min=1"`
	NameField        *string        `form:"name_field" json:"name_field" binding:"omitempty,min=1"`
	EmailField       *string        `form:"email_field" json:"email_field" binding:"omitempty,min=1"`
	PictureField     *string        `form:"picture_field" json:"picture_field"`
	DescriptionField *string        `form:"description_field" json:"description_field"`
	On               *bool          `form:"on" json:"on"`
	Picture          *model.FileURL `form:"picture" json:"picture"`
}
