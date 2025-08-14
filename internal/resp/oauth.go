package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetOauthResp(oauth model.Oauth) gin.H {
	return gin.H{
		"id":                oauth.ID,
		"auth_url":          oauth.AuthURL,
		"token_url":         oauth.TokenURL,
		"user_info_url":     oauth.UserInfoURL,
		"callback_url":      oauth.CallbackURL,
		"client_id":         oauth.ClientID,
		"client_secret":     oauth.ClientSecret,
		"provider":          oauth.Provider,
		"uri":               oauth.Uri,
		"resp_id_field":     oauth.RespIDField,
		"resp_name_field":   oauth.RespNameField,
		"resp_email_field":  oauth.RespEmailField,
		"resp_avatar_field": oauth.RespAvatarField,
		"resp_desc_field":   oauth.RespDescField,
		"on":                oauth.On,
		"avatar":            oauth.Avatar,
	}
}
