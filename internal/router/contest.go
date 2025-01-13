package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetContests(ctx *gin.Context) {
	var getContestsForm GetContestsForm
	var (
		contests []model.Contest
		total    int64
		ok       bool
		msg      string
	)
	trace := middleware.GetTraceID(ctx)
	if err := ctx.ShouldBind(&getContestsForm); err == nil {
		if middleware.GetSelf(ctx).Type == "admin" {
			contests, total, ok, msg = db.GetContests(getContestsForm.Limit, getContestsForm.Offset, true)
		} else {
			contests, total, ok, msg = db.GetContests(getContestsForm.Limit, getContestsForm.Offset, false)
		}
		if !ok {
			log.Logger.Infof("| %s | %s", trace, msg)
			ctx.JSONP(http.StatusInternalServerError, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": nil})
			return
		}
		ctx.JSONP(http.StatusOK, gin.H{"trace": trace, "msg": utils.M(ctx, msg), "data": gin.H{
			"contests": utils.TidyRetData(contests, "captcha"),
			"total":    total,
		}})
		return
	} else {
		log.Logger.Infof("| %s | %s", trace, err)
		ctx.JSONP(http.StatusBadRequest, gin.H{"trace": trace, "msg": utils.M(ctx, "BadRequest"), "data": nil})
		return
	}
}
