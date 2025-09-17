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

	"github.com/gin-gonic/gin"
)

func SubmitFlag(ctx *gin.Context) {
	var form f.SubmitFlagForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.SubmitFlagEventType)
	user := middleware.GetSelf(ctx).(model.User)
	team := middleware.GetTeam(ctx)
	contest := middleware.GetContest(ctx)
	challenge := middleware.GetChallenge(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	contestFlags, _, ok, msg := db.InitContestFlagRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx := db.DB.Begin()
	result, _, ok, msg := service.Submit(tx, user, team, contest, contestChallenge, form, ctx.ClientIP())
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	if contestChallenge.Type == model.PodsChallengeType && service.CheckIfSolved(db.DB, team, contestFlags) {
		go func() {
			victim, ok, _ := db.InitVictimRepo(db.DB).HasAliveVictim(team.ID, challenge.ID)
			if !ok {
				return
			}
			gtx := db.DB.Begin()
			if ok, _ = service.StopVictim(gtx, victim); !ok {
				gtx.Rollback()
				return
			}
			gtx.Commit()
		}()
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": result, "data": nil})
}

func GetContestFlags(ctx *gin.Context) {
	contestChallenge := middleware.GetContestChallenge(ctx)
	contestFlags, _, ok, msg := db.InitContestFlagRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, contestFlag := range contestFlags {
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
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateContestChallengeFlagEventType)
	contestChallenge := middleware.GetContestChallenge(ctx)
	contestFlag := middleware.GetContestFlag(ctx)
	if contestChallenge.Type == model.QuestionChallengeType && form.Value != nil {
		form.Value = &contestFlag.Value
	}
	currentScore := contestFlag.CurrentScore
	if form.Score != nil && *form.Score < currentScore {
		currentScore = *form.Score
	}
	ok, msg := db.InitContestFlagRepo(db.DB).Update(contestFlag.ID, db.UpdateContestFlagOptions{
		Value:        form.Value,
		Score:        form.Score,
		CurrentScore: &currentScore,
		Decay:        form.Decay,
		MinScore:     form.MinScore,
		ScoreType:    form.ScoreType,
	})
	if ok {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
