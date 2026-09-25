package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
)

func GetPermissionResp(perm model.Permission) gin.H {
	return gin.H{
		"id":          perm.ID,
		"name":        perm.Name,
		"resource":    perm.Resource,
		"operation":   perm.Operation,
		"description": perm.Description,
	}
}
