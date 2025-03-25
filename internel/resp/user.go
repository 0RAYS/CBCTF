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

func GetUserResp(user model.User, all bool) gin.H {
	data := gin.H{
		"name":     user.Name,
		"email":    user.Email,
		"country":  user.Country,
		"avatar":   fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(user.Avatar, "/")),
		"desc":     user.Desc,
		"verified": user.Verified,
		"hidden":   user.Hidden,
	}
	if all {
		data["id"] = user.ID
		data["banned"] = user.Banned
		data["score"] = user.Score
		data["solved"] = user.Solved
		data["teams"] = len(user.Teams)
		data["contests"] = len(user.Contests)
		data["devices"] = func() []string {
			var devices []string
			for _, device := range user.Devices {
				devices = append(devices, device.Magic)
			}
			return devices
		}
		data["cheats"] = func() []string {
			var cheats []string
			for _, cheat := range user.Cheats {
				cheats = append(cheats, cheat.ID)
			}
			return cheats
		}
	}
	return data
}
