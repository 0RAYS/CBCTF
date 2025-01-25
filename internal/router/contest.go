package router

import (
	"CBCTF/internal/constants"
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetContest(ctx *gin.Context) {
	contest, ok, msg := db.GetContestByID(ctx, middleware.GetContestID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": contest})
}

func GetContestCaptcha(ctx *gin.Context) {
	contest, ok, msg := db.GetContestByID(ctx, middleware.GetContestID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": contest.Captcha})
}

func GetContests(ctx *gin.Context) {
	var form constants.GetModelsForm
	all := false
	if middleware.GetRole(ctx) == "admin" {
		all = true
	}
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contests, count, ok, msg := db.GetContests(ctx, form.Limit, form.Offset, all)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "contests": contests}})
}

func CreateContest(ctx *gin.Context) {
	var form constants.CreateContestForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contest, ok, msg := db.CreateContest(ctx, form)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": contest})
}

func UpdateContest(ctx *gin.Context) {
	var form constants.UpdateContestForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contest, ok, msg := db.GetContestByID(ctx, middleware.GetContestID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := utils.Form2Map(form)
	if name, ok := data["name"]; ok && name.(string) != contest.Name && !db.IsUniqueTeamName(name.(string), contest.ID) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "ContestNameExists", "data": nil})
		return
	}
	_, msg = db.UpdateContest(ctx, contest.ID, data)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteContest(ctx *gin.Context) {
	_, msg := db.DeleteContest(ctx, middleware.GetContestID(ctx))
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
