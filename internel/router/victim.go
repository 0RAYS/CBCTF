package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/middleware"
	"CBCTF/internel/model"
	"CBCTF/internel/redis"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func StartVictim(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	user := middleware.GetSelf(ctx).(model.User)
	if ok, err := redis.CheckVictimCreate(team.ID, contestChallenge.ChallengeID); ok || err != nil {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"msg": i18n.TooManyRequests, "data": nil})
		return
	}
	if err := redis.RecordVictimCreate(team.ID, contestChallenge.ChallengeID); err != nil {
		log.Logger.Warningf("Failed to record container create: %v", err)
	}
	tx := db.DB.WithContext(ctx).Begin()
	_, ok, msg := service.StartVictim(tx, user, team, contestChallenge)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	status := service.GetVictimStatus(db.DB.WithContext(ctx), team, contestChallenge)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": status})
}

func IncreaseVictimDuration(ctx *gin.Context) {
	team := middleware.GetTeam(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	DB := db.DB.WithContext(ctx)
	repo := db.InitVictimRepo(DB)
	victims, _, ok, msg := repo.ListWithConditions(-1, -1, db.GetOptions{
		{Key: "team_id", Value: team.ID, Op: "and"},
		{Key: "contest_challenge_id", Value: contestChallenge.ID, Op: "and"},
	}, false, "Pods")
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
	_, msg := service.StopVictim(db.DB.WithContext(ctx), team, contestChallenge)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func GetVictim(ctx *gin.Context) {
	victim := middleware.GetVictim(ctx)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": resp.GetVictimResp(victim)})
}

func GetVictims(ctx *gin.Context) {
	var form f.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	if _, exists := ctx.GetQuery("limit"); !exists {
		form.Limit = 5
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		form.Offset = 0
	}
	team := middleware.GetTeam(ctx)
	repo := db.InitVictimRepo(db.DB.WithContext(ctx))
	victims, count, ok, msg := repo.ListWithConditions(form.Limit, form.Offset, db.GetOptions{
		{Key: "team_id", Value: team.ID, Op: "and"},
	}, true)
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
