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
		"id_claim":          oauth.IDClaim,
		"name_claim":        oauth.NameClaim,
		"email_claim":       oauth.EmailClaim,
		"picture_claim":     oauth.PictureClaim,
		"description_claim": oauth.DescriptionClaim,
		"groups_claim":      oauth.GroupsClaim,
		"admin_group":       oauth.AdminGroup,
		"on":                oauth.On,
		"picture":           oauth.Picture,
	}
}
