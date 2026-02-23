package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
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
	user := middleware.GetSelf(ctx)
	team := middleware.GetTeam(ctx)
	contest := middleware.GetContest(ctx)
	challenge := middleware.GetChallenge(ctx)
	contestChallenge := middleware.GetContestChallenge(ctx)
	go func() {
		_, ret := service.StartVictim(db.DB, user.ID, team.ID, contest.ID, contestChallenge.ID, challenge.ID)
		if !ret.OK {
			websocket.Send(user.ID, wm.ErrorLevel, wm.StartVictimWSType, "Start Victim", "Failed")
			victim, ret := db.InitVictimRepo(db.DB).HasAliveVictim(team.ID, challenge.ID)
			if !ret.OK {
				return
			}
			service.StopVictim(db.DB, victim)
			return
		}
		websocket.Send(user.ID, wm.SuccessLevel, wm.StartVictimWSType, "Start Victim", "Done")
		return
	}()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal())
}

func IncreaseVictimDuration(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.IncreaseVictimEventType)
	team := middleware.GetTeam(ctx)
	challenge := middleware.GetChallenge(ctx)
	repo := db.InitVictimRepo(db.DB)
	victim, ret := repo.HasAliveVictim(team.ID, challenge.ID)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	if !victim.Start.Add(victim.Duration).Before(time.Now().Add(20 * time.Minute)) {
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Model.Victim.HasMuchTime})
		return
	}
	duration := victim.Duration + time.Hour
	if ret = db.InitVictimRepo(db.DB).Update(victim.ID, db.UpdateVictimOptions{
		Duration: &duration,
	}); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{
		"target":    victim.RemoteAddr(),
		"remaining": victim.Remaining().Seconds(),
		"status":    "Running",
	}))
}

func StopVictim(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.StopVictimEventType)
	team := middleware.GetTeam(ctx)
	challenge := middleware.GetChallenge(ctx)
	victim, ret := db.InitVictimRepo(db.DB).HasAliveVictim(team.ID, challenge.ID)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	if ret = service.StopVictim(db.DB, victim); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, ret)
}

func GetVictims(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	team := middleware.GetTeam(ctx)
	victims, count, ret := db.InitVictimRepo(db.DB).List(form.Limit, form.Offset, db.GetOptions{
		Conditions: map[string]any{"team_id": team.ID},
		Preloads:   map[string]db.GetOptions{"ContestChallenge": {}},
		Deleted:    true,
	})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, victim := range victims {
		data = append(data, resp.GetVictimResp(victim))
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"victims": data, "count": count}))
}
