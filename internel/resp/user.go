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

func GetUserResp(user model.User) gin.H {
	data := gin.H{
		"name":     user.Name,
		"email":    user.Email,
		"country":  user.Country,
		"avatar":   fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(user.Avatar, "/")),
		"desc":     user.Desc,
		"verified": user.Verified,
		"hidden":   user.Hidden,
	}
	return data
}
