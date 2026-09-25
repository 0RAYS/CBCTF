package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/view"
)

func GetGroupResp(groupView view.GroupView) gin.H {
	group := groupView.Group
	data := gin.H{
		"id":          group.ID,
		"name":        group.Name,
		"description": group.Description,
		"default":     group.Default,
		"role_id":     group.RoleID,
		"users":       groupView.UserCount,
	}
	return data
}
