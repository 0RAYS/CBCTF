package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetRoleResp(role model.Role) gin.H {
	data := gin.H{
		"id":          role.ID,
		"name":        role.Name,
		"description": role.Description,
		"default":     role.Default,
	}
	return data
}
