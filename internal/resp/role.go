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
	if len(role.Permissions) > 0 {
		permissions := make([]gin.H, 0, len(role.Permissions))
		for _, perm := range role.Permissions {
			permissions = append(permissions, GetPermissionResp(perm))
		}
		data["permissions"] = permissions
	}
	return data
}
