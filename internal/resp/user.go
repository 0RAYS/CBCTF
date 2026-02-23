package resp

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetUserResp(user model.User, admin bool) gin.H {
	data := gin.H{
		"id":          user.ID,
		"name":        user.Name,
		"email":       user.Email,
		"picture":     user.Picture,
		"description": user.Description,
		"verified":    user.Verified,
		"score":       user.Score,
		"solved":      user.Solved,
		"provider":    user.Provider,
		"hidden":      user.Hidden,
		"banned":      user.Banned,
		"is_admin":    db.InitUserRepo(db.DB).IsInGroup(user.ID, model.AdminGroupName),
	}
	if admin {
		data["teams"] = db.InitUserRepo(db.DB).CountAssociation(user, "Teams")
		data["contests"] = db.InitUserRepo(db.DB).CountAssociation(user, "Contests")
	}
	return data
}
