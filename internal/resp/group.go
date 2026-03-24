package resp

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetGroupResp(group model.Group) gin.H {
	users, _ := db.InitGroupRepo(db.DB).CountUsers(group.ID)
	data := gin.H{
		"id":          group.ID,
		"name":        group.Name,
		"description": group.Description,
		"default":     group.Default,
		"role_id":     group.RoleID,
		"users":       users,
	}
	return data
}
