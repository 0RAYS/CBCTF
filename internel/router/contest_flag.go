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

func GetContestFlags(ctx *gin.Context) {
	contestChallenge := middleware.GetContestChallenge(ctx)
	data := make([]gin.H, 0)
	for _, contestFlag := range contestChallenge.ContestFlags {
		data = append(data, resp.GetContestFlagResp(contestFlag))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data})
}

func GetContestFlag(ctx *gin.Context) {
	contestFlag := middleware.GetContestFlag(ctx)
	data := resp.GetContestFlagResp(contestFlag)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data})
}

func UpdateContestFlag(ctx *gin.Context) {
	var form f.UpdateContestFlagForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	contestFlag := middleware.GetContestFlag(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := db.InitContestFlagRepo(tx).Update(contestFlag.ID, db.UpdateContestFlagOptions{
		Value:     form.Value,
		Score:     form.Score,
		Decay:     form.Decay,
		MinScore:  form.MinScore,
		ScoreType: form.ScoreType,
	})
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
