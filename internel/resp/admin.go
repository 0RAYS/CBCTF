package resp

import (
	"CBCTF/internel/config"
	"CBCTF/internel/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func AdminLoginResp(admin model.Admin) gin.H {
	return gin.H{
		"name":   admin.Name,
		"email":  admin.Email,
		"avatar": fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(admin.Avatar, "/")),
	}
}
