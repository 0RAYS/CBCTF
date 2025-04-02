package resp

import (
	"CBCTF/internel/config"
	"CBCTF/internel/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

var LoginResp = GetUserResp
var RegisterResp = GetUserResp

func GetUserResp(user model.User, admin bool) gin.H {
	data := gin.H{
		"id":       user.ID,
		"name":     user.Name,
		"email":    user.Email,
		"country":  user.Country,
		"avatar":   fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(user.Avatar, "/")),
		"desc":     user.Desc,
		"verified": user.Verified,
	}
	if admin {
		data["hidden"] = user.Hidden
		data["verified"] = user.Verified
		data["banned"] = user.Banned
		data["teams"] = len(user.Teams)
		data["contests"] = len(user.Contests)
	}
	return data
}
