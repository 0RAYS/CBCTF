package router

import (
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	db "CBCTF/internal/repo"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetContest(ctx *gin.Context) {
	contest := middleware.GetContest(ctx)
	champion, _, _, _ := service.GetTeamRanking(db.DB.WithContext(ctx), contest.ID, 1, 0)
	data := resp.GetContestResp(contest, middleware.IsAdmin(ctx))
	data["highest"] = 0
	if len(champion) > 0 {
		data["highest"] = champion[0].Score
	}
	data["solved"], _, _ = db.InitSubmissionRepo(db.DB.WithContext(ctx)).Count(db.CountOptions{
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
	if _, exists := ctx.GetQuery("limit"); !exists {
		form.Limit = 5
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		form.Offset = 0
	}
	options := db.GetOptions{
		Preloads: map[string]db.GetOptions{
			"Teams":   {},
			"Users":   {},
			"Notices": {},
		},
	}
	if !middleware.IsAdmin(ctx) {
		options.Conditions = map[string]any{"hidden": false}
	}
	contests, count, ok, msg := db.InitContestRepo(db.DB.WithContext(ctx)).List(form.Limit, form.Offset, options)
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
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
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
