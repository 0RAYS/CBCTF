package resp

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetGroupResp(group model.Group) gin.H {
	data := gin.H{
		"id":          group.ID,
		"name":        group.Name,
		"description": group.Description,
		"default":     group.Default,
		"role_id":     group.RoleID,
		"users":       db.InitGroupRepo(db.DB).CountAssociation(group, "Users"),
	}
	return data
}
