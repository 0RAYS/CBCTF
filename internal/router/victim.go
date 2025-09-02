package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"CBCTF/internal/websocket"
	wm "CBCTF/internal/websocket/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func StartVictim(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.StartVictimEventType)
	user := middleware.GetSelf(ctx).(model.User)
	team := middleware.GetTeam(ctx)
	challenge := middleware.GetChallenge(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	go func() {
		_, ok, _ := service.StartVictim(db.DB, user.ID, team.ID, contestChallenge.ID, challenge.ID)
		if !ok {
			websocket.Send(false, user.ID, wm.ErrorLevel, wm.StartVictimWSType, "Start Victim", "Failed")
			go func() {
				victim, ok, _ := db.InitVictimRepo(db.DB).HasAliveVictim(team.ID, challenge.ID)
				if !ok {
					return
				}
				service.StopVictim(db.DB, victim)
			}()
			return
		}
		websocket.Send(false, user.ID, wm.SuccessLevel, wm.StartVictimWSType, "Start Victim", "Done")
		return
	}()
	status := service.GetVictimStatus(db.DB, team.ID, challenge)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": status})
}

func IncreaseVictimDuration(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.IncreaseVictimEventType)
	team := middleware.GetTeam(ctx)
	challenge := middleware.GetChallenge(ctx)
	DB := db.DB
	repo := db.InitVictimRepo(DB)
	victim, ok, msg := repo.HasAliveVictim(team.ID, challenge.ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	if !victim.Start.Add(victim.Duration).Before(time.Now().Add(20 * time.Minute)) {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.HasMuchTime, "data": nil})
		return
	}
	tx := DB.Begin()
	repo = db.InitVictimRepo(tx)
	duration := victim.Duration + time.Hour
	if ok, msg = repo.Update(victim.ID, db.UpdateVictimOptions{
		Duration: &duration,
	}); !ok {
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
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data[0]})
}

func StopVictim(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.StopVictimEventType)
	team := middleware.GetTeam(ctx)
	challenge := middleware.GetChallenge(ctx)
	victim, ok, msg := db.InitVictimRepo(db.DB).HasAliveVictim(team.ID, challenge.ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ok, msg = service.StopVictim(db.DB, victim)
	if ok {
		ctx.Set(middleware.CTXEventSuccessKey, true)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func GetVictims(ctx *gin.Context) {
	var form f.GetModelsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	team := middleware.GetTeam(ctx)
	repo := db.InitVictimRepo(db.DB)
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
