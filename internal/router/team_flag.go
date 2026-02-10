package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
	"CBCTF/internal/task"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTeamFlags(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.ReadFlagEventType)
	team := middleware.GetTeam(ctx)
	teamFlags, _, ret := db.InitTeamFlagRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"team_id": team.ID},
		Preloads: map[string]db.GetOptions{"ContestFlag": {
			Preloads: map[string]db.GetOptions{"ContestChallenge": {}},
		}},
	})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	challengeFlagMap := make(map[uint]gin.H)
	for _, flag := range teamFlags {
		if _, ok := challengeFlagMap[flag.ContestFlag.ContestChallengeID]; !ok {
			challengeFlagMap[flag.ContestFlag.ContestChallengeID] = gin.H{
				"name":     flag.ContestFlag.ContestChallenge.Name,
				"type":     flag.ContestFlag.ContestChallenge.Type,
				"category": flag.ContestFlag.ContestChallenge.Category,
				"hidden":   flag.ContestFlag.ContestChallenge.Hidden,
				"flags":    make([]gin.H, 0),
			}
		}
		challengeFlagMap[flag.ContestFlag.ContestChallengeID]["flags"] = append(challengeFlagMap[flag.ContestFlag.ContestChallengeID]["flags"].([]gin.H), gin.H{
			"value":         flag.Value,
			"solved":        flag.Solved,
			"template":      flag.ContestFlag.Value,
			"init_score":    flag.ContestFlag.Score,
			"current_score": flag.ContestFlag.CurrentScore,
			"decay":         flag.ContestFlag.Decay,
			"min_score":     flag.ContestFlag.MinScore,
			"solvers":       flag.ContestFlag.Solvers,
		})
	}
	data := make([]gin.H, 0)
	for _, flags := range challengeFlagMap {
		data = append(data, flags)
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(data))
}

func InitTeamFlag(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.InitChallengeEventType)
	user := middleware.GetSelf(ctx).(model.User)
	team := middleware.GetTeam(ctx)
	contest := middleware.GetContest(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	contestFlags, _, ret := db.InitContestFlagRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	contestChallenge.ContestFlags = contestFlags
	challenge := middleware.GetChallenge(ctx)
	tx := db.DB.Begin()
	teamFlags, ret := service.CreateTeamFlag(tx, team, contest, contestChallenge)
	if !ret.OK {
		tx.Rollback()
		ctx.JSON(http.StatusOK, ret)
		return
	}
	if challenge.Type == model.DynamicChallengeType {
		if _, err := task.EnqueueGenAttachmentTask(user.ID, challenge, team, teamFlags); err != nil {
			log.Logger.Warningf("Failed to enqueue gen attachment task: %s", err)
			tx.Rollback()
			ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Task.EnqueueError, Attr: map[string]any{"Error": err.Error()}})
			return
		}
	}
	tx.Commit()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal())
}

func ResetTeamFlag(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.ResetChallengeEventType)
	user := middleware.GetSelf(ctx).(model.User)
	team := middleware.GetTeam(ctx)
	contest := middleware.GetContest(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	contestFlags, _, ret := db.InitContestFlagRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"contest_challenge_id": contestChallenge.ID},
	})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	contestChallenge.ContestFlags = contestFlags
	challenge := middleware.GetChallenge(ctx)
	tx := db.DB.Begin()
	teamFlags, ret := service.UpdateTeamFlag(tx, team, contest, contestChallenge)
	if !ret.OK {
		tx.Rollback()
		ctx.JSON(http.StatusOK, ret)
		return
	}
	switch challenge.Type {
	case model.DynamicChallengeType:
		if _, err := task.EnqueueGenAttachmentTask(user.ID, challenge, team, teamFlags); err != nil {
			log.Logger.Warningf("Failed to enqueue gen attachment task: %s", err)
			tx.Rollback()
			ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Task.EnqueueError, Attr: map[string]any{"Error": err.Error()}})
			return
		}
		tx.Commit()
	case model.PodsChallengeType:
		tx.Commit()
		// 不考虑失败
		go func() {
			victim, ret := db.InitVictimRepo(db.DB).HasAliveVictim(team.ID, challenge.ID)
			if !ret.OK {
				return
			}
			service.StopVictim(db.DB, victim)
		}()
	default:
		tx.Commit()
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, ret)
}
