package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSubmissions(all bool) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var form f.GetModelsForm
		if err := ctx.ShouldBind(&form); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
			return
		}
		var (
			submissions = make([]model.Submission, 0)
			count       int64
			ok          bool
			msg         string
		)
		DB := db.DB.WithContext(ctx)
		if all {
			submissions, count, ok, msg = db.InitSubmissionRepo(DB).GetAll(form.Limit, form.Offset, false, 0)
		} else {
			team := middleware.GetTeam(ctx)
			submissions, count, ok, msg = db.InitSubmissionRepo(DB).GetAllByKeyID("team_id", team.ID, form.Limit, form.Offset, false, 0, false)
		}
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"submissions": &submissions, "count": count}})
	}
}
