package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"

	"github.com/gin-gonic/gin"
)

func GetSubmissions(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	team := middleware.GetTeam(ctx)
	submissions, count, ret := db.InitSubmissionRepo(db.DB).List(form.Limit, form.Offset, db.GetOptions{
		Conditions: map[string]any{"team_id": team.ID},
	})
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, submission := range submissions {
		data = append(data, resp.GetSubmissionResp(submission))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"submissions": data, "count": count}))
}
