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
	"os"

	"github.com/gin-gonic/gin"
)

func GetContestChallenges(ctx *gin.Context) {
	var form f.GetChallengesForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if _, ok := ctx.GetQuery("limit"); !ok {
		form.Limit = 10
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		form.Offset = 0
	}
	options := db.GetOptions{
		Conditions: map[string]any{"contest_id": middleware.GetContest(ctx).ID},
		Preloads:   map[string]db.GetOptions{"Challenge": {}, "ContestFlags": {}},
	}
	if form.Category != "" {
		options.Conditions["category"] = form.Category
	}
	if middleware.IsAdmin(ctx) && form.Type != "" {
		options.Conditions["type"] = form.Type
	}
	if !middleware.IsAdmin(ctx) {
		options.Conditions["hidden"] = false
	}
	contestChallengeL, count, ok, msg := db.InitContestChallengeRepo(db.DB.WithContext(ctx)).
		List(form.Limit, form.Offset, options)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, contestChallenge := range contestChallengeL {
		tmp := resp.GetContestChallengeResp(contestChallenge.Challenge, contestChallenge)
		if !middleware.IsAdmin(ctx) {
			team := middleware.GetTeam(ctx)
			tmp["hidden"] = false
			tmp["attempts"] = service.CountAttempts(db.DB.WithContext(ctx), team, contestChallenge)
			tmp["init"] = service.CheckIfGenerated(db.DB.WithContext(ctx), team, contestChallenge.ContestFlags)
			tmp["solved"] = service.CheckIfSolved(db.DB.WithContext(ctx), team, contestChallenge.ContestFlags)
			tmp["remote"] = service.GetTeamVictimStatus(db.DB.WithContext(ctx), team, contestChallenge)
			tmp["file"] = func() string {
				if _, err := os.Stat(contestChallenge.Challenge.AttachmentPath(team.ID)); err != nil {
					return ""
				}
				return contestChallenge.Challenge.AttachmentPath(team.ID)
			}()
		}
		data = append(data, tmp)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": gin.H{"challenges": data, "count": count}})
}

func GetContestChallenge(ctx *gin.Context) {
	challenge := middleware.GetChallenge(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	data := resp.GetContestChallengeResp(challenge, contestChallenge)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data})
}

func GetContestChallengeStatus(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	challenge := middleware.GetChallenge(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	data := gin.H{
		"attempts": service.CountAttempts(db.DB.WithContext(ctx), team, contestChallenge),
		"init":     service.CheckIfGenerated(db.DB.WithContext(ctx), team, contestChallenge.ContestFlags),
		"solved":   service.CheckIfSolved(db.DB.WithContext(ctx), team, contestChallenge.ContestFlags),
		"remote":   service.GetTeamVictimStatus(db.DB.WithContext(ctx), team, contestChallenge),
		"file": func() string {
			if _, err := os.Stat(challenge.AttachmentPath(team.ID)); err != nil {
				return ""
			}
			return "attachment.zip"
		}(),
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data})
}

func AddContestChallenge(ctx *gin.Context) {
	var form f.CreateContestChallengeForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateContestChallengeEventType)
	contestChallengeL, failedL, _, _ := service.CreateContestChallenge(db.DB.WithContext(ctx), middleware.GetContest(ctx), form)
	data := make([]gin.H, 0)
	for _, contestChallenge := range contestChallengeL {
		data = append(data, resp.GetContestChallengeResp(contestChallenge.Challenge, contestChallenge))
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": gin.H{"contest_challenge": data, "failed": failedL}})
}

func UpdateContestChallenge(ctx *gin.Context) {
	var form f.UpdateContestChallengeForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateContestChallengeEventType)
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
		ctx.Set(middleware.CTXEventSuccessKey, true)
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteContestChallenge(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteContestChallengeEventType)
	contestChallenge := middleware.GetContestChallenge(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := db.InitContestChallengeRepo(tx).Delete(contestChallenge.ID)
	if !ok {
		tx.Rollback()
	} else {
		ctx.Set(middleware.CTXEventSuccessKey, true)
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
