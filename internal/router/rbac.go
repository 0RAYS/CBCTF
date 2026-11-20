package router

import (
	"CBCTF/internal/log"
	"CBCTF/internal/redis"
)

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
