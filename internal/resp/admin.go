package resp

import (
	"CBCTF/internal/model"
	"github.com/gin-gonic/gin"
)

func GetAdminResp(admin model.Admin) gin.H {
	data := gin.H{
		"id":       admin.ID,
		"name":     admin.Name,
		"email":    admin.Email,
		"avatar":   admin.Avatar,
		"verified": admin.Verified,
	}
	return data
}
