package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
	"CBCTF/internal/view"
)

func GetUserResp(userView view.UserView, admin bool) gin.H {
	user := userView.User
	data := gin.H{
		"id":               user.ID,
		"name":             user.Name,
		"email":            user.Email,
		"picture":          user.Picture,
		"description":      user.Description,
		"verified":         user.Verified,
		"score":            user.Score,
		"solved":           user.Solved,
		"provider":         user.Provider,
		"has_no_pwd":       user.Password == model.NeverLoginPWD,
		"hidden":           user.Hidden,
		"banned":           user.Banned,
		"has_admin_access": userView.HasAdminAccess,
	}
	if admin {
		data["teams"] = userView.TeamCount
		data["contests"] = userView.ContestCount
	}
	return data
}
