package resp

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetUserResp(user model.User, admin bool) gin.H {
	data := gin.H{
		"id":       user.ID,
		"name":     user.Name,
		"email":    user.Email,
		"picture":  user.Picture,
		"desc":     user.Desc,
		"verified": user.Verified,
	}
	if admin {
		data["hidden"] = user.Hidden
		data["banned"] = user.Banned
		data["teams"] = db.InitUserRepo(db.DB).CountAssociation(user, "Teams")
		data["contests"] = db.InitUserRepo(db.DB).CountAssociation(user, "Contests")
	}
	return data
}
