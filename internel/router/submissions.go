package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/middleware"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSubmissions(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	if _, exists := ctx.GetQuery("limit"); !exists {
		form.Limit = 10
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		form.Offset = 0
	}
	team := middleware.GetTeam(ctx)
	submissions, count, ok, msg := db.InitSubmissionRepo(db.DB.WithContext(ctx)).
		ListWithConditions(form.Limit, form.Offset, db.GetOptions{
			{Key: "team_id", Value: team.ID, Op: "and"},
		}, false)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, submission := range submissions {
		data = append(data, resp.GetSubmissionResp(submission))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": gin.H{"submissions": data, "count": count}})
}
