package resp

import (
	"CBCTF/internal/db"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetUserResp(user model.User, admin bool) gin.H {
	userRepo := db.InitUserRepo(db.DB)
	data := gin.H{
		"id":               user.ID,
		"name":             user.Name,
		"email":            user.Email,
		"picture":          user.Picture,
		"description":      user.Description,
		"verified":         user.Verified,
		"score":            user.Score,
		"solved":           user.Solved,
		"provider":         user.Provider,
		"has_no_pwd":       user.Password == model.NeverLoginPWD,
		"hidden":           user.Hidden,
		"banned":           user.Banned,
		"has_admin_access": db.InitPermissionRepo(db.DB).HasAdminAccess(user.ID),
	}
	if admin {
		data["teams"], _ = userRepo.CountTeams(user.ID)
		data["contests"], _ = userRepo.CountContests(user.ID)
	}
	return data
}
