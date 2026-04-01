package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"

	"github.com/gin-gonic/gin"
)

func SubmitFlag(ctx *gin.Context) {
	var form dto.SubmitFlagForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.SubmitFlagEventType)
	user := middleware.GetSelf(ctx)
	team := middleware.GetTeam(ctx)
	contest := middleware.GetContest(ctx)
	challenge := middleware.GetChallenge(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	contestFlags, _, ret := db.InitContestFlagRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ret = db.WithTransaction(func(tx *db.Tx) model.RetVal {
		_, ret := service.Submit(tx, user, team, contest, contestChallenge, form, ctx.ClientIP())
		return ret
	})
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	if contestChallenge.Type == model.PodsChallengeType && service.CheckIfSolved(db.DB, team, contestFlags) {
		go func() {
			victim, ret := db.InitVictimRepo(db.DB).HasAliveVictim(team.ID, challenge.ID)
			if !ret.OK {
				return
			}
			service.ForceStopVictim(db.DB, victim)
		}()
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, ret)
}

func GetContestFlags(ctx *gin.Context) {
	contestChallenge := middleware.GetContestChallenge(ctx)
	contestFlags, _, ret := db.InitContestFlagRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, contestFlag := range contestFlags {
		data = append(data, resp.GetContestFlagResp(contestFlag))
	}
	resp.JSON(ctx, model.SuccessRetVal(data))
}

func UpdateContestFlag(ctx *gin.Context) {
	var form dto.UpdateContestFlagForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
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
	ret := db.InitContestFlagRepo(db.DB).Update(contestFlag.ID, db.UpdateContestFlagOptions{
		Value:        form.Value,
		Score:        form.Score,
		CurrentScore: &currentScore,
		Decay:        form.Decay,
		MinScore:     form.MinScore,
		ScoreType:    form.ScoreType,
	})
	if ret.OK {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	resp.JSON(ctx, ret)
}
