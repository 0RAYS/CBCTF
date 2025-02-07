package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/middleware"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetContest(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": middleware.GetContest(ctx)})
}

func GetContestCaptcha(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": middleware.GetContest(ctx).Captcha})
}

func GetContests(ctx *gin.Context) {
	var form f.GetModelsForm
	all := false
	if middleware.GetRole(ctx) == "admin" {
		all = true
	}
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contests, count, ok, msg := db.GetContests(db.DB.WithContext(ctx), form.Limit, form.Offset, all)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "contests": contests}})
}

func CreateContest(ctx *gin.Context) {
	var form f.CreateContestForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	contest, ok, msg := db.CreateContest(tx, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": contest})
}

func UpdateContest(ctx *gin.Context) {
	var form f.UpdateContestForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	contest := middleware.GetContest(ctx)
	data := utils.Form2Map(form)
	if name, ok := data["name"]; ok && name.(string) != contest.Name && !db.IsUniqueTeamName(db.DB.WithContext(ctx), name.(string), contest.ID) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "ContestNameExists", "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := db.UpdateContest(tx, contest.ID, data)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteContest(ctx *gin.Context) {
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := db.DeleteContest(tx, middleware.GetContest(ctx))
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
