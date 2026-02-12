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
		"id_field":          oauth.IDField,
		"name_field":        oauth.NameField,
		"email_field":       oauth.EmailField,
		"picture_field":     oauth.PictureField,
		"description_field": oauth.DescriptionField,
		"on":                oauth.On,
		"picture":           oauth.Picture,
	}
}
