package resp

import (
	"CBCTF/internel/config"
	"CBCTF/internel/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func GetAdminResp(admin model.Admin) gin.H {
	data := gin.H{
		"id":       admin.ID,
		"name":     admin.Name,
		"email":    admin.Email,
		"avatar":   fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(admin.Avatar, "/")),
		"verified": admin.Verified,
	}
	return data
}
