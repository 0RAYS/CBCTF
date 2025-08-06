package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type OauthCallbackForm struct {
	Code  string `form:"code" json:"code" binding:"required"`
	State string `form:"state" json:"state" binding:"required"`
}

func (f *OauthCallbackForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}

type CreateOauthProviderForm struct {
	AuthURL         string `form:"auth_url" json:"auth_url" binding:"required"`
	TokenURL        string `form:"token_url" json:"token_url" binding:"required"`
	UserInfoURL     string `form:"user_info_url" json:"user_info_url" binding:"required"`
	CallbackURL     string `form:"callback_url" json:"callback_url" binding:"required"`
	ClientID        string `form:"client_id" json:"client_id" binding:"required"`
	ClientSecret    string `form:"client_secret" json:"client_secret" binding:"required"`
	Provider        string `form:"provider" json:"provider" binding:"required"`
	URI             string `form:"uri" json:"uri" binding:"required"`
	RespIDField     string `form:"id_field" json:"id_field" binding:"required"`
	RespNameField   string `form:"name_field" json:"name_field" binding:"required"`
	RespEmailField  string `form:"email_field" json:"email_field" binding:"required"`
	RespAvatarField string `form:"avatar_field" json:"avatar_field"`
	RespDescField   string `form:"desc_field" json:"desc_field"`
}

func (f *CreateOauthProviderForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return false, i18n.BadRequest
	}
	f.AuthURL = strings.TrimSpace(f.AuthURL)
	if !strings.HasPrefix(f.AuthURL, "http://") && !strings.HasPrefix(f.AuthURL, "https://") {
		return false, i18n.BadRequest
	}
	f.TokenURL = strings.TrimSpace(f.TokenURL)
	if !strings.HasPrefix(f.TokenURL, "http://") && !strings.HasPrefix(f.TokenURL, "https://") {
		return false, i18n.BadRequest
	}
	f.UserInfoURL = strings.TrimSpace(f.UserInfoURL)
	if !strings.HasPrefix(f.UserInfoURL, "http://") && !strings.HasPrefix(f.UserInfoURL, "https://") {
		return false, i18n.BadRequest
	}
	f.CallbackURL = strings.TrimSpace(f.CallbackURL)
	if !strings.HasPrefix(f.CallbackURL, "http://") && !strings.HasPrefix(f.CallbackURL, "https://") {
		return false, i18n.BadRequest
	}
	f.ClientID = strings.TrimSpace(f.ClientID)
	if f.ClientID == "" {
		return false, i18n.BadRequest
	}
	f.ClientSecret = strings.TrimSpace(f.ClientSecret)
	if f.ClientSecret == "" {
		return false, i18n.BadRequest
	}
	f.Provider = strings.TrimSpace(f.Provider)
	if f.Provider == "" {
		f.Provider = utils.UUID()
	}
	f.URI = strings.TrimSpace(f.URI)
	if f.URI == "" {
		f.URI = utils.RandStr(10)
	}
	f.RespIDField = strings.TrimSpace(f.RespIDField)
	if f.RespIDField == "" {
		f.RespIDField = "{id}"
	}
	f.RespNameField = strings.TrimSpace(f.RespNameField)
	if f.RespNameField == "" {
		f.RespNameField = "{name}"
	}
	f.RespEmailField = strings.TrimSpace(f.RespEmailField)
	if f.RespEmailField == "" {
		f.RespEmailField = "{email}"
	}
	f.RespAvatarField = strings.TrimSpace(f.RespAvatarField)
	f.RespDescField = strings.TrimSpace(f.RespDescField)
	return true, i18n.Success
}

type UpdateOauthProviderForm struct {
	AuthURL         *string          `form:"auth_url" json:"auth_url"`
	TokenURL        *string          `form:"token_url" json:"token_url"`
	UserInfoURL     *string          `form:"user_info_url" json:"user_info_url"`
	CallbackURL     *string          `form:"callback_url" json:"callback_url"`
	ClientID        *string          `form:"client_id" json:"client_id"`
	ClientSecret    *string          `form:"client_secret" json:"client_secret"`
	Provider        *string          `form:"provider" json:"provider"`
	URI             *string          `form:"uri" json:"uri"`
	RespIDField     *string          `form:"id_field" json:"id_field"`
	RespNameField   *string          `form:"name_field" json:"name_field"`
	RespEmailField  *string          `form:"email_field" json:"email_field"`
	RespAvatarField *string          `form:"avatar_field" json:"avatar_field"`
	RespDescField   *string          `form:"desc_field" json:"desc_field"`
	On              *bool            `form:"on" json:"on"`
	Avatar          *model.AvatarURL `form:"avatar" json:"avatar"`
}

func (f *UpdateOauthProviderForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return false, i18n.BadRequest
	}
	if f.AuthURL != nil {
		*f.AuthURL = strings.TrimSpace(*f.AuthURL)
		if !strings.HasPrefix(*f.AuthURL, "http://") && !strings.HasPrefix(*f.AuthURL, "https://") {
			return false, i18n.BadRequest
		}
	}
	if f.TokenURL != nil {
		*f.TokenURL = strings.TrimSpace(*f.TokenURL)
		if !strings.HasPrefix(*f.TokenURL, "http://") && !strings.HasPrefix(*f.TokenURL, "https://") {
			return false, i18n.BadRequest
		}
	}
	if f.UserInfoURL != nil {
		*f.UserInfoURL = strings.TrimSpace(*f.UserInfoURL)
		if !strings.HasPrefix(*f.UserInfoURL, "http://") && !strings.HasPrefix(*f.UserInfoURL, "https://") {
			return false, i18n.BadRequest
		}
	}
	if f.CallbackURL != nil {
		*f.CallbackURL = strings.TrimSpace(*f.CallbackURL)
		if !strings.HasPrefix(*f.CallbackURL, "http://") && !strings.HasPrefix(*f.CallbackURL, "https://") {
			return false, i18n.BadRequest
		}
	}
	if f.ClientID != nil {
		*f.ClientID = strings.TrimSpace(*f.ClientID)
		if *f.ClientID == "" {
			return false, i18n.BadRequest
		}
	}
	if f.ClientSecret != nil {
		*f.ClientSecret = strings.TrimSpace(*f.ClientSecret)
		if *f.ClientSecret == "" {
			return false, i18n.BadRequest
		}
	}
	if f.Provider != nil {
		*f.Provider = strings.TrimSpace(*f.Provider)
		if *f.Provider == "" {
			*f.Provider = utils.UUID()
		}
	}
	if f.URI != nil {
		*f.URI = strings.TrimSpace(*f.URI)
		if *f.URI == "" {
			*f.URI = utils.RandStr(10)
		}
	}
	if f.RespIDField != nil {
		*f.RespIDField = strings.TrimSpace(*f.RespIDField)
		if *f.RespIDField == "" {
			*f.RespIDField = "id"
		}
	}
	if f.RespNameField != nil {
		*f.RespNameField = strings.TrimSpace(*f.RespNameField)
		if *f.RespNameField == "" {
			*f.RespNameField = "name"
		}
	}
	if f.RespEmailField != nil {
		*f.RespEmailField = strings.TrimSpace(*f.RespEmailField)
		if *f.RespEmailField == "" {
			*f.RespEmailField = "email"
		}
	}
	if f.RespAvatarField != nil {
		*f.RespAvatarField = strings.TrimSpace(*f.RespAvatarField)
	}
	if f.RespDescField != nil {
		*f.RespDescField = strings.TrimSpace(*f.RespDescField)
	}
	if f.Avatar != nil {
		*f.Avatar = model.AvatarURL(strings.TrimSpace(string(*f.Avatar)))
	}
	return true, i18n.Success
}
