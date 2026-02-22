package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetContestChallenges(ctx *gin.Context) {
	var form dto.GetChallengesForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
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
	contestChallengeL, count, ret := db.InitContestChallengeRepo(db.DB).List(form.Limit, form.Offset, options)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, contestChallenge := range contestChallengeL {
		tmp := resp.GetContestChallengeResp(contestChallenge)
		if !middleware.IsAdmin(ctx) {
			team := middleware.GetTeam(ctx)
			tmp["hidden"] = false
			tmp["attempts"] = service.CountAttempts(db.DB, team, contestChallenge)
			tmp["init"] = service.CheckIfGenerated(db.DB, team, contestChallenge.ContestFlags)
			tmp["solved"] = service.CheckIfSolved(db.DB, team, contestChallenge.ContestFlags)
			tmp["remote"] = service.GetVictimStatus(db.DB, team.ID, contestChallenge.Challenge)
			tmp["file"] = func() string {
				if _, err := os.Stat(contestChallenge.Challenge.AttachmentPath(team.ID)); err != nil {
					return ""
				}
				return contestChallenge.Challenge.AttachmentPath(team.ID)
			}()
		}
		data = append(data, tmp)
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"challenges": data, "count": count}))
}

func GetContestChallengeCategories(ctx *gin.Context) {
	var form dto.GetCategoriesForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	contest := middleware.GetContest(ctx)
	categories, ret := db.InitContestChallengeRepo(db.DB).ListCategories(contest.ID, form.Type)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(categories))
}

func GetContestChallengeStatus(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	challenge := middleware.GetChallenge(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	contestFlags, _, ret := db.InitContestFlagRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := gin.H{
		"attempts": service.CountAttempts(db.DB, team, contestChallenge),
		"init":     service.CheckIfGenerated(db.DB, team, contestFlags),
		"solved":   service.CheckIfSolved(db.DB, team, contestFlags),
		"remote":   service.GetVictimStatus(db.DB, team.ID, challenge),
		"file": func() string {
			path := challenge.AttachmentPath(team.ID)
			record, _ := db.InitFileRepo(db.DB).Get(db.GetOptions{
				Conditions: map[string]any{"model": challenge.ModelName(), "model_id": challenge.ID, "type": model.ChallengeFileType}},
			)
			filename := "attachment.zip"
			if record.Path == path {
				filename = record.Filename
			}
			if _, err := os.Stat(path); err != nil {
				return ""
			}
			return filename
		}(),
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(data))
}

func AddContestChallenge(ctx *gin.Context) {
	var form dto.CreateContestChallengeForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateContestChallengeEventType)
	contestChallengeL, failedL, _ := service.CreateContestChallenge(db.DB, middleware.GetContest(ctx), form)
	data := make([]gin.H, 0)
	for _, contestChallenge := range contestChallengeL {
		data = append(data, resp.GetContestChallengeResp(contestChallenge))
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"contest_challenge": data, "failed": failedL}))
}

func UpdateContestChallenge(ctx *gin.Context) {
	var form dto.UpdateContestChallengeForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateContestChallengeEventType)
	contestChallenge := middleware.GetContestChallenge(ctx)
	ret := db.InitContestChallengeRepo(db.DB).Update(contestChallenge.ID, db.UpdateContestChallengeOptions{
		Name:        form.Name,
		Description: form.Description,
		Hidden:      form.Hidden,
		Attempt:     form.Attempt,
		Hints:       form.Hints,
		Tags:        form.Tags,
	})
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}

func DeleteContestChallenge(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteContestChallengeEventType)
	contestChallenge := middleware.GetContestChallenge(ctx)
	tx := db.DB.Begin()
	ret := db.InitContestChallengeRepo(tx).Delete(contestChallenge.ID)
	if !ret.OK {
		tx.Rollback()
	} else {
		tx.Commit()
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, ret)
}
