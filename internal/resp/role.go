package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
)

func GetRoleResp(role model.Role) gin.H {
	return gin.H{
		"id":          role.ID,
		"name":        role.Name,
		"description": role.Description,
		"default":     role.Default,
	}
}
