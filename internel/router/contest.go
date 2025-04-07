package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/middleware"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetContest(ctx *gin.Context) {
	contest := middleware.GetContest(ctx)
	champion, _, _, _ := service.GetTeamRanking(db.DB.WithContext(ctx), contest.ID, 1, 0)
	data := resp.GetContestResp(contest, middleware.GetRole(ctx) == "admin")
	data["highest"] = 0
	if len(champion) > 0 {
		data["highest"] = champion[0].Score
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": data})
}

func GetContests(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	all := middleware.GetRole(ctx) == "admin"
	contests, count, ok, msg := db.InitContestRepo(db.DB.WithContext(ctx)).GetAll(form.Limit, form.Offset, all, true)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, contest := range contests {
		data = append(data, resp.GetContestResp(contest, all))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"contests": data, "count": count}})
}

func CreateContest(ctx *gin.Context) {
	var form f.CreateContestForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	contest, ok, msg := service.CreateContest(tx, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	contest.Duration = time.Duration(contest.Duration.Seconds())
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": &contest})
}

func UpdateContest(ctx *gin.Context) {
	var form f.UpdateContestForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	contest := middleware.GetContest(ctx)
	ok, msg := service.UpdateContest(tx, contest, form)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteContest(ctx *gin.Context) {
	tx := db.DB.WithContext(ctx).Begin()
	contest := middleware.GetContest(ctx)
	ok, msg := db.InitContestRepo(tx).Delete(contest.ID)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
