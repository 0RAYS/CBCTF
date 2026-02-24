package middleware

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckPermission(ctx *gin.Context) {
	key := fmt.Sprintf("%s %s", ctx.Request.Method, ctx.FullPath())
	permission, ok := model.RoutePermissions[key]
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Response.Forbidden})
		return
	}

	pass, ret := db.InitPermissionRepo(db.DB).CheckUserPermission(GetSelf(ctx).ID, permission)
	if !ret.OK {
		ctx.AbortWithStatusJSON(http.StatusOK, ret)
		return
	}
	if !pass {
		ctx.AbortWithStatusJSON(http.StatusOK, model.RetVal{Msg: i18n.Response.Forbidden})
		return
	}
	ctx.Next()
}
