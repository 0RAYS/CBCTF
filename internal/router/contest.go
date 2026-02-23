package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
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
	champion, _, _ := service.GetTeamRanking(db.DB, contest, 1, 0)
	data := resp.GetContestResp(contest, middleware.IsFullAccess(ctx))
	data["highest"] = 0
	if len(champion) > 0 {
		data["highest"] = champion[0].Score
	}
	data["solved"], _ = db.InitSubmissionRepo(db.DB).Count(db.CountOptions{
		Conditions: map[string]any{"solved": true, "contest_id": contest.ID},
	})
	ctx.JSON(http.StatusOK, model.SuccessRetVal(data))
}

func GetContests(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	if _, ok := ctx.GetQuery("limit"); !ok {
		form.Limit = 5
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		form.Offset = 0
	}
	options := db.GetOptions{}
	if !middleware.IsFullAccess(ctx) {
		options.Conditions = map[string]any{"hidden": false}
	}
	contests, count, ret := db.InitContestRepo(db.DB).List(form.Limit, form.Offset, options)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, contest := range contests {
		data = append(data, resp.GetContestResp(contest, middleware.IsFullAccess(ctx)))
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"contests": data, "count": count}))
}

func CreateContest(ctx *gin.Context) {
	var form dto.CreateContestForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateContestEventType)
	contest, ret := service.CreateContest(db.DB, form)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	contest.Duration = time.Duration(contest.Duration.Seconds())
	ctx.JSON(http.StatusOK, model.SuccessRetVal(contest))
}

func UpdateContest(ctx *gin.Context) {
	var form dto.UpdateContestForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateContestEventType)
	contest := middleware.GetContest(ctx)
	ret := service.UpdateContest(db.DB, contest, form)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}

func DeleteContest(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteContestEventType)
	contest := middleware.GetContest(ctx)
	tx := db.DB.Begin()
	ret := db.InitContestRepo(tx).Delete(contest.ID)
	if !ret.OK {
		tx.Rollback()
	} else {
		tx.Commit()
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}
