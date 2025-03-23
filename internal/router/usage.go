package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func AddUsage(ctx *gin.Context) {
	var form f.CreateUsageForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	usages, ok, msg := db.CreateUsage(tx, form, middleware.GetContest(ctx).ID)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": &usages})
}

func GetUsages(ctx *gin.Context) {
	var (
		usages  []model.Usage
		ok      bool
		msg     string
		all     = middleware.GetRole(ctx) == "admin"
		DB      = db.DB.WithContext(ctx)
		contest = middleware.GetContest(ctx)
		team    = middleware.GetTeam(ctx)
	)
	usages, ok, msg = db.GetUsageByContestID(DB, contest.ID, all)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	var challenges []map[string]interface{}
	for _, usage := range usages {
		tmp := map[string]interface{}{}
		if !all {
			usage.Flag = ""
			usage.Flags = make([]string, 0)
			usage.DockerImage = ""
			usage.GeneratorImage = ""
			usage.NetworkPolicy.From = make([]utils.IPBlock, 0)
			usage.NetworkPolicy.To = make([]utils.IPBlock, 0)
			usage.NetworkPolicies = make([]utils.NetworkPolicy, 0)
		}
		tmp["usage"] = usage
		if !all {
			tmp["status"] = gin.H{
				"solved":   db.IsSolved(DB, contest.ID, team.ID, usage.ChallengeID),
				"attempts": db.CountAttempts(DB, contest.ID, team.ID, usage.ChallengeID),
				"init": func() bool {
					_, ok, _ = db.GetFlagBy3ID(db.DB.WithContext(ctx), contest.ID, team.ID, usage.ChallengeID)
					return ok
				}(),
				"files": func() string {
					if _, err := os.Stat(usage.AttachmentPath(team.ID)); err != nil {
						if !os.IsNotExist(err) {
							log.Logger.Warningf("Failed to check attachment: %s", err)
						}
						return ""
					}
					return model.AttachmentFile
				}(),
				"remote": func() gin.H {
					if usage.Type == model.Docker {
						if docker, ok, _ := db.GetContainerBy3ID(db.DB.WithContext(ctx), contest.ID, team.ID, usage.ChallengeID); ok {
							return gin.H{
								"target":    docker.RemoteAddr(),
								"remaining": docker.Remaining().Seconds(),
							}
						}
					}
					return gin.H{
						"target":    "",
						"remaining": "",
					}
				}(),
			}
		}
		challenges = append(challenges, tmp)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": &challenges})
}

func RemoveUsage(ctx *gin.Context) {
	var DB = db.DB.WithContext(ctx)
	usage := middleware.GetUsage(ctx)
	tx := DB.Begin()
	ok, msg := db.DeleteUsage(tx, usage.ID)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func UpdateUsage(ctx *gin.Context) {
	var DB = db.DB.WithContext(ctx)
	usage := middleware.GetUsage(ctx)
	var form f.UpdateUsageForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	data := utils.Form2Map(form)
	tx := DB.Begin()
	ok, msg := db.UpdateUsage(tx, usage.ID, data)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
