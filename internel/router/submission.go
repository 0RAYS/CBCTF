package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSubmissions(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	if _, exists := ctx.GetQuery("limit"); !exists {
		form.Limit = 10
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		form.Offset = 0
	}
	var (
		submissions = make([]model.Submission, 0)
		count       int64
		ok          bool
		msg         string
	)
	DB := db.DB.WithContext(ctx)
	team := middleware.GetTeam(ctx)
	submissions, count, ok, msg = db.InitSubmissionRepo(DB).GetAllByKeyID("team_id", team.ID, form.Limit, form.Offset, false)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, submission := range submissions {
		data = append(data, resp.GetSubmissionResp(submission))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": gin.H{"submissions": data, "count": count}})
}
