package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/middleware"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetContest(ctx *gin.Context) {
	contest := middleware.GetContest(ctx)
	champion, _, _, _ := service.GetTeamRanking(db.DB.WithContext(ctx), contest.ID, 1, 0)
	data := resp.GetContestResp(contest)
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
	contests, count, ok, msg := db.InitContestRepo(db.DB.WithContext(ctx)).GetAll(form.Limit, form.Offset, false, 0, all)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if !all {
		data := make([]gin.H, 0)
		for _, contest := range contests {
			data = append(data, resp.GetContestResp(contest))
		}
		ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"contests": contests, "count": count}})
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"contests": contests, "count": count}})
	return
}
