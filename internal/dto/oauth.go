package dto

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

type OauthCallbackForm struct {
	Code  string `form:"code" json:"code" binding:"required"`
	State string `form:"state" json:"state" binding:"required"`
}

func (f *OauthCallbackForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

type CreateOauthProviderForm struct {
	AuthURL          string `form:"auth_url" json:"auth_url" binding:"required"`
	TokenURL         string `form:"token_url" json:"token_url" binding:"required"`
	UserInfoURL      string `form:"user_info_url" json:"user_info_url" binding:"required"`
	CallbackURL      string `form:"callback_url" json:"callback_url" binding:"required"`
	ClientID         string `form:"client_id" json:"client_id" binding:"required"`
	ClientSecret     string `form:"client_secret" json:"client_secret" binding:"required"`
	Provider         string `form:"provider" json:"provider" binding:"required"`
	Uri              string `form:"uri" json:"uri" binding:"required"`
	RespIDField      string `form:"id_field" json:"id_field" binding:"required"`
	RespNameField    string `form:"name_field" json:"name_field" binding:"required"`
	RespEmailField   string `form:"email_field" json:"email_field" binding:"required"`
	RespPictureField string `form:"picture_field" json:"picture_field"`
	RespDescField    string `form:"desc_field" json:"desc_field"`
}

func (f *CreateOauthProviderForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

type UpdateOauthProviderForm struct {
	AuthURL          *string        `form:"auth_url" json:"auth_url"`
	TokenURL         *string        `form:"token_url" json:"token_url"`
	UserInfoURL      *string        `form:"user_info_url" json:"user_info_url"`
	CallbackURL      *string        `form:"callback_url" json:"callback_url"`
	ClientID         *string        `form:"client_id" json:"client_id"`
	ClientSecret     *string        `form:"client_secret" json:"client_secret"`
	Provider         *string        `form:"provider" json:"provider"`
	Uri              *string        `form:"uri" json:"uri"`
	RespIDField      *string        `form:"resp_id_field" json:"resp_id_field"`
	RespNameField    *string        `form:"resp_name_field" json:"resp_name_field"`
	RespEmailField   *string        `form:"resp_email_field" json:"resp_email_field"`
	RespPictureField *string        `form:"resp_picture_field" json:"resp_picture_field"`
	RespDescField    *string        `form:"resp_desc_field" json:"resp_desc_field"`
	On               *bool          `form:"on" json:"on"`
	Picture          *model.FileURL `form:"picture" json:"picture"`
}

func (f *UpdateOauthProviderForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
