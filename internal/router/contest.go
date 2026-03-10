package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
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
	resp.JSON(ctx, model.SuccessRetVal(data))
}

func GetContests(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	options := db.GetOptions{Sort: []string{"id DESC"}}
	if !middleware.IsFullAccess(ctx) {
		options.Conditions = map[string]any{"hidden": false}
	}
	contests, count, ret := db.InitContestRepo(db.DB).List(form.Limit, form.Offset, options)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, contest := range contests {
		data = append(data, resp.GetContestResp(contest, middleware.IsFullAccess(ctx)))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"contests": data, "count": count}))
}

func CreateContest(ctx *gin.Context) {
	var form dto.CreateContestForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateContestEventType)
	contest, ret := service.CreateContest(db.DB, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	contest.Duration = time.Duration(contest.Duration.Seconds())
	resp.JSON(ctx, model.SuccessRetVal(contest))
}

func UpdateContest(ctx *gin.Context) {
	var form dto.UpdateContestForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateContestEventType)
	contest := middleware.GetContest(ctx)
	ret := service.UpdateContest(db.DB, contest, form)
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
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
	resp.JSON(ctx, ret)
}
