package router

import (
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	db "CBCTF/internal/repo"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func StartVictim(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	user := middleware.GetSelf(ctx).(model.User)
	tx := db.DB.WithContext(ctx).Begin()
	_, ok, msg := service.StartTeamVictim(tx, user, team, contestChallenge)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	status := service.GetTeamVictimStatus(db.DB.WithContext(ctx), team, contestChallenge)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": status})
}

func IncreaseVictimDuration(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	DB := db.DB.WithContext(ctx)
	repo := db.InitVictimRepo(DB)
	victims, _, ok, msg := repo.List(-1, -1, db.GetOptions{
		Conditions: map[string]any{
			"team_id":              team.ID,
			"contest_challenge_id": contestChallenge.ID,
		},
		Preloads: map[string]db.GetOptions{
			"Pods": {},
		},
	})
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, victim := range victims {
		if !victim.Start.Add(victim.Duration).Before(time.Now().Add(20 * time.Minute)) {
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.HasMuchTime, "data": nil})
			return
		}
		tx := DB.Begin()
		repo = db.InitVictimRepo(tx)
		duration := victim.Duration + 1*time.Hour
		ok, msg := repo.Update(victim.ID, db.UpdateVictimOptions{
			Duration: &duration,
		})
		if !ok {
			tx.Rollback()
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		tx.Commit()
		data = append(data, gin.H{
			"target":    victim.RemoteAddr(),
			"remaining": victim.Remaining().Seconds(),
			"status":    "Running",
		})
	}
	if len(data) > 0 {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data[0]})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.VictimNotFound, "data": nil})
	}
}

func StopVictim(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	_, msg := service.StopTeamVictim(db.DB.WithContext(ctx), team, contestChallenge)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func GetVictim(ctx *gin.Context) {
	victim := middleware.GetVictim(ctx)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": resp.GetVictimResp(victim)})
}

func GetVictims(ctx *gin.Context) {
	var form f.GetModelsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	team := middleware.GetTeam(ctx)
	repo := db.InitVictimRepo(db.DB.WithContext(ctx))
	victims, count, ok, msg := repo.List(form.Limit, form.Offset, db.GetOptions{
		Conditions: map[string]any{"team_id": team.ID},
		Deleted:    true,
	})
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, victim := range victims {
		data = append(data, resp.GetVictimResp(victim))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": gin.H{"victims": data, "count": count}})
}
