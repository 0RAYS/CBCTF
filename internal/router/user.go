package router

import (
	"RayWar/internal/db"
	"RayWar/internal/log"
	"RayWar/internal/middleware"
	"RayWar/internal/model"
	"RayWar/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUsers(ctx *gin.Context) {
	var getUsersForm GetUsersForm
	var (
		users []model.User
		total int64
		ok    bool
		msg   string
	)
	trace := middleware.GetTraceID(ctx)
	if err := ctx.ShouldBind(&getUsersForm); err == nil {
		if middleware.GetSelf(ctx).Type == "admin" {
			users, total, ok, msg = db.GetUsers(getUsersForm.Limit, getUsersForm.Offset, true)
		} else {
			users, total, ok, msg = db.GetUsers(getUsersForm.Limit, getUsersForm.Offset, false)
		}
		if !ok {
			log.Logger.Infof("| %s | %s", trace, msg)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
		ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": gin.H{
			"users": utils.TidyRetData(users, "password", "email"),
			"total": total,
		}})
		return
	} else {
		log.Logger.Infof("| %s | %s", trace, err)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		return
	}
}
