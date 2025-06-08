package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/middleware"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func GetContestChallenges(ctx *gin.Context) {
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
	conditions := db.GetOptions{
		{Key: "contest_id", Value: middleware.GetContest(ctx).ID, Op: "and"},
	}
	if middleware.GetRole(ctx) != "admin" {
		conditions = append(conditions, db.GetOption{Key: "hidden", Value: false, Op: "and"})
	}
	contestChallengeL, count, ok, msg := db.InitContestChallengeRepo(db.DB.WithContext(ctx)).
		ListWithConditions(form.Limit, form.Offset, conditions, false, "Challenge", "ContestFlags")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, contestChallenge := range contestChallengeL {
		//TODO
		data = append(data, resp.GetContestChallengeResp(contestChallenge))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": gin.H{"challenges": data, "count": count}})
}

func GetContestChallenge(ctx *gin.Context) {
	contestChallenge := middleware.GetContestChallenge(ctx)
	data := resp.GetContestChallengeResp(contestChallenge)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data})
}

func GetContestChallengeStatus(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	data := gin.H{
		"attempts": service.CountAttempts(db.DB.WithContext(ctx), team, contestChallenge),
		"init":     service.CheckIfGenerated(db.DB.WithContext(ctx), team, contestChallenge),
		"solved":   service.CheckIfSolved(db.DB.WithContext(ctx), team, contestChallenge),
		//TODO
		//"remote":   service.GetVictimStatus(db.DB.WithContext(ctx), team, contestChallenge),
		"file": func() string {
			if _, err := os.Stat(contestChallenge.Challenge.AttachmentPath(team.ID)); err != nil {
				return ""
			}
			return "attachment.zip"
		}(),
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data})
}

func AddContestChallenge(ctx *gin.Context) {
	var form f.CreateContestChallengeForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	contestChallengeL, failedL, _, _ := service.CreateContestChallenge(db.DB.WithContext(ctx), middleware.GetContest(ctx), form)
	data := make([]gin.H, 0)
	for _, contestChallenge := range contestChallengeL {
		data = append(data, resp.GetContestChallengeResp(contestChallenge))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": gin.H{"contest_challenge": data, "failed": failedL}})
}

func UpdateContestChallenge(ctx *gin.Context) {
	var form f.UpdateContestChallengeForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	contestChallenge := middleware.GetContestChallenge(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := db.InitContestChallengeRepo(tx).Update(contestChallenge.ID, db.UpdateContestChallengeOptions{
		Name:    form.Name,
		Desc:    form.Desc,
		Hidden:  form.Hidden,
		Attempt: form.Attempt,
		Hints:   form.Hints,
		Tags:    form.Tags,
	})
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteContestChallenge(ctx *gin.Context) {
	contestChallenge := middleware.GetContestChallenge(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := db.InitContestChallengeRepo(tx).Delete(contestChallenge.ID)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
