package router

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/resp"
)

// respondRBACChange answers a permission change and drops the cached
// permission sets when it succeeded.
func respondRBACChange(ctx *gin.Context, ret model.RetVal) {
	if ret.OK {
		invalidateRBAC()
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}

// invalidateRBAC drops the cached permission sets after a role or group
// change, because a stale cache would keep revoked permissions alive until it
// expires.
func invalidateRBAC() {
	if ret := redis.DeleteRBAC(); !ret.OK {
		log.Logger.Warningf("Failed to invalidate permission cache: %v", ret)
	}
}

// invalidateUserRBAC drops the cached permission set of one user after their
// group membership changed.
func invalidateUserRBAC(userID uint) {
	if ret := redis.DeleteUserRBAC(userID); !ret.OK {
		log.Logger.Warningf("Failed to invalidate permission cache: user_id=%d reason=%v", userID, ret)
	}
}
