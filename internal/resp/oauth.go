package resp

import (
	"CBCTF/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetOauthResp(oauth model.Oauth) gin.H {
	runes := []rune(oauth.ClientSecret)
	var masked string
	if len(runes) >= 2 {
		masked = fmt.Sprintf("%s******%s", string(runes[0]), string(runes[len(runes)-1]))
	} else {
		masked = fmt.Sprintf("error: too short")
	}
	return gin.H{
		"id":                oauth.ID,
		"auth_url":          oauth.AuthURL,
		"token_url":         oauth.TokenURL,
		"user_info_url":     oauth.UserInfoURL,
		"callback_url":      oauth.CallbackURL,
		"client_id":         oauth.ClientID,
		"client_secret":     masked,
		"provider":          oauth.Provider,
		"uri":               oauth.URI,
		"resp_id_field":     oauth.RespIDField,
		"resp_name_field":   oauth.RespNameField,
		"resp_email_field":  oauth.RespEmailField,
		"resp_avatar_field": oauth.RespAvatarField,
		"resp_desc_field":   oauth.RespDescField,
		"on":                oauth.On,
		"avatar":            oauth.Avatar,
	}
}
