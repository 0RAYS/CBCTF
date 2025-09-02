package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetContest(ctx *gin.Context) {
	contest := middleware.GetContest(ctx)
	champion, _, _, _ := service.GetTeamRanking(db.DB, contest.ID, 1, 0)
	data := resp.GetContestResp(contest, middleware.IsAdmin(ctx))
	data["highest"] = 0
	if len(champion) > 0 {
		data["highest"] = champion[0].Score
	}
	data["solved"], _, _ = db.InitSubmissionRepo(db.DB).Count(db.CountOptions{
		Conditions: map[string]any{"solved": true, "contest_id": contest.ID},
	})
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data})
}

func GetContests(ctx *gin.Context) {
	var form f.GetModelsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if _, ok := ctx.GetQuery("limit"); !ok {
		form.Limit = 5
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		form.Offset = 0
	}
	options := db.GetOptions{Preloads: map[string]db.GetOptions{"Teams": {}, "Users": {}, "Notices": {}}}
	if !middleware.IsAdmin(ctx) {
		options.Conditions = map[string]any{"hidden": false}
	}
	contests, count, ok, msg := db.InitContestRepo(db.DB).List(form.Limit, form.Offset, options)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, contest := range contests {
		data = append(data, resp.GetContestResp(contest, middleware.IsAdmin(ctx)))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"contests": data, "count": count}})
}

func CreateContest(ctx *gin.Context) {
	var form f.CreateContestForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateContestEventType)
	contest, ok, msg := service.CreateContest(db.DB, form)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	contest.Duration = time.Duration(contest.Duration.Seconds())
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": &contest})
}

func UpdateContest(ctx *gin.Context) {
	var form f.UpdateContestForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateContestEventType)
	contest := middleware.GetContest(ctx)
	ok, msg := service.UpdateContest(db.DB, contest, form)
	if ok {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteContest(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteContestEventType)
	tx := db.DB.Begin()
	contest := middleware.GetContest(ctx)
	ok, msg := db.InitContestRepo(tx).Delete(contest.ID)
	if !ok {
		tx.Rollback()
	} else {
		ctx.Set(middleware.CTXEventSuccessKey, true)
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
