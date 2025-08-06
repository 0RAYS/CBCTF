package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

var LoginResp = GetUserResp
var RegisterResp = GetUserResp

func GetUserResp(user model.User, admin bool) gin.H {
	data := gin.H{
		"id":       user.ID,
		"name":     user.Name,
		"email":    user.Email,
		"country":  user.Country,
		"avatar":   user.Avatar,
		"desc":     user.Desc,
		"verified": user.Verified,
	}
	if admin {
		data["hidden"] = user.Hidden
		data["banned"] = user.Banned
		data["teams"] = len(user.Teams)
		data["contests"] = len(user.Contests)
	}
	return data
}
