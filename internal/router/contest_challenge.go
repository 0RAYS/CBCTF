package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"os"

	"github.com/gin-gonic/gin"
)

func GetContestChallenges(ctx *gin.Context) {
	var form dto.GetContestChallengesForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	options := db.GetOptions{
		Conditions: map[string]any{"contest_id": middleware.GetContest(ctx).ID, "hidden": false},
		Preloads:   map[string]db.GetOptions{"Challenge": {}, "ContestFlags": {}},
	}
	if form.Category != "" {
		options.Conditions["category"] = form.Category
	}
	contestChallengeL, count, ret := db.InitContestChallengeRepo(db.DB).List(form.Limit, form.Offset, options)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	team := middleware.GetTeam(ctx)
	for _, contestChallenge := range contestChallengeL {
		tmp := resp.GetContestChallengeResp(contestChallenge)
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
		data = append(data, tmp)
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"challenges": data, "count": count}))
}

func GetAllContestChallenges(ctx *gin.Context) {
	var form dto.GetAllContestChallengesForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	options := db.GetOptions{
		Conditions: map[string]any{"contest_id": middleware.GetContest(ctx).ID},
		Preloads:   map[string]db.GetOptions{"Challenge": {}, "ContestFlags": {}},
	}
	if form.Category != "" {
		options.Conditions["category"] = form.Category
	}
	if middleware.IsFullAccess(ctx) && form.Type != "" {
		options.Conditions["type"] = form.Type
	}
	contestChallengeL, count, ret := db.InitContestChallengeRepo(db.DB).List(form.Limit, form.Offset, options)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, contestChallenge := range contestChallengeL {
		tmp := resp.GetContestChallengeResp(contestChallenge)
		data = append(data, tmp)
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"challenges": data, "count": count}))
}

func GetContestChallengeCategories(ctx *gin.Context) {
	var form dto.GetCategoriesForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	contest := middleware.GetContest(ctx)
	categories, ret := db.InitContestChallengeRepo(db.DB).ListCategories(contest.ID, form.Type)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(categories))
}

func GetContestChallengeStatus(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	challenge := middleware.GetChallenge(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	contestFlags, _, ret := db.InitContestFlagRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ret.OK {
		resp.JSON(ctx, ret)
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
			if string(record.Path) == path {
				filename = record.Filename
			}
			if _, err := os.Stat(path); err != nil {
				return ""
			}
			return filename
		}(),
	}
	resp.JSON(ctx, model.SuccessRetVal(data))
}

func AddContestChallenge(ctx *gin.Context) {
	var form dto.CreateContestChallengeForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateContestChallengeEventType)
	contestChallengeL, failedL, _ := service.CreateContestChallenge(db.DB, middleware.GetContest(ctx), form)
	data := make([]gin.H, 0)
	for _, contestChallenge := range contestChallengeL {
		data = append(data, resp.GetContestChallengeResp(contestChallenge))
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"contest_challenge": data, "failed": failedL}))
}

func UpdateContestChallenge(ctx *gin.Context) {
	var form dto.UpdateContestChallengeForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
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
	resp.JSON(ctx, ret)
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
	resp.JSON(ctx, ret)
}

func GetContestFlagSolvers(ctx *gin.Context) {
	contestFlag := middleware.GetContestFlag(ctx)
	rows, ret := db.InitSubmissionRepo(db.DB).ListFlagSolvers(contestFlag.ID)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		data = append(data, gin.H{
			"user_id":   row.UserID,
			"user_name": row.UserName,
			"team_id":   row.TeamID,
			"team_name": row.TeamName,
			"score":     row.Score,
			"solved_at": row.SolvedAt,
		})
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"solvers": data, "count": int64(len(data))}))
}
