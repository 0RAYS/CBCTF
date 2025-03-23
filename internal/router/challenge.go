package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func CreateChallenge(ctx *gin.Context) {
	var form f.CreateChallengeForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	form.Category = utils.ToTitle(strings.TrimSpace(form.Category))
	tx := db.DB.WithContext(ctx).Begin()
	challenge, ok, msg := db.CreateChallenge(tx, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	if err := os.MkdirAll(challenge.BasicDir(), 0755); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": "CreateDirError", "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": &challenge})
}

func GetChallenge(ctx *gin.Context) {
	challenge := middleware.GetChallenge(ctx)
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": &challenge})
}

func GetChallenges(ctx *gin.Context) {
	var form f.GetChallengesForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	if ctx.Query("type") == "" && ctx.Query("category") == "" {
		form.Type = ""
		form.Category = ""
	}
	challenges, count, ok, msg := db.GetChallenges(db.DB.WithContext(ctx), form.Limit, form.Offset, form.Type, form.Category)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "challenges": &challenges}})
}

func GetAttachment(ctx *gin.Context) {
	usage := middleware.GetUsage(ctx)
	team := middleware.GetTeam(ctx)
	path := usage.AttachmentPath(team.ID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": "FileNotFound", "data": nil})
		return
	}
	ctx.File(path)
}

func GetChallengeFiles(ctx *gin.Context) {
	challenge := middleware.GetChallenge(ctx)
	var files []string
	dir, err := os.ReadDir(challenge.BasicDir())
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": "ReadDirError", "data": nil})
		return
	}
	for _, file := range dir {
		files = append(files, file.Name())
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": &files})
}

func GetCategories(ctx *gin.Context) {
	var form f.GetCategoriesForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	categories, ok, msg := db.GetCategories(db.DB.WithContext(ctx), form.Type)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": &categories})
}

func UpdateChallenge(ctx *gin.Context) {
	var (
		challenge model.Challenge
		msg       string
		data      map[string]interface{}
	)
	var form f.UpdateChallengeForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	challenge = middleware.GetChallenge(ctx)
	tmp := utils.ToTitle(strings.TrimSpace(*form.Category))
	form.Category = &tmp
	data = utils.Form2Map(form)
	if t, ok := data["type"]; ok && !db.IsValidChallengeType(t.(string)) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "InvalidChallengeType", "data": nil})
		return
	}
	if category, ok := data["category"]; ok && category.(string) != challenge.Category {
		data["category"] = category.(string)
	}
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := db.UpdateChallenge(tx, challenge.ID, data)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteChallenge(ctx *gin.Context) {
	var form f.DeleteChallengeForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	challenge := middleware.GetChallenge(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	usages, ok, msg := db.GetUsageByChallengeID(tx, challenge.ID)
	if ok {
		for _, usage := range usages {
			if ok, msg := db.DeleteUsage(tx, usage.ID); !ok {
				tx.Rollback()
				ctx.JSON(http.StatusInternalServerError, gin.H{"msg": msg, "data": nil})
				return
			}
		}
	}
	ok, msg = db.DeleteChallenge(tx, challenge.ID)
	if !ok {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	if form.Force && os.RemoveAll(challenge.BasicDir()) != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "UnknownError", "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func ChallengeStatus(ctx *gin.Context) {
	data := gin.H{
		"init":  false,
		"files": "",
		"remote": gin.H{
			"target":    "",
			"remaining": "",
		},
		"solved": false,
	}
	team := middleware.GetTeam(ctx)
	usage := middleware.GetUsage(ctx)
	contest := middleware.GetContest(ctx)
	if _, ok, msg := db.GetFlagBy3ID(db.DB.WithContext(ctx), contest.ID, team.ID, usage.ChallengeID); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": data})
		return
	}
	data["init"] = true
	if db.IsSolved(db.DB.WithContext(ctx), contest.ID, team.ID, usage.ChallengeID) {
		data["solved"] = true
	}
	if _, err := os.Stat(usage.AttachmentPath(team.ID)); err != nil {
		if !os.IsNotExist(err) {
			log.Logger.Warningf("Failed to check attachment: %s", err)
		}
	} else {
		data["files"] = model.AttachmentFile
	}
	if usage.Type == model.Docker {
		if container, ok, _ := db.GetContainerBy3ID(db.DB.WithContext(ctx), contest.ID, team.ID, usage.ChallengeID); ok {
			data["remote"] = gin.H{
				"target":    container.RemoteAddress(),
				"remaining": container.Remaining().Seconds(),
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": data})
}

func InitChallenge(reset bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			team    = middleware.GetTeam(ctx)
			contest = middleware.GetContest(ctx)
			usage   = middleware.GetUsage(ctx)
			DB      = db.DB.WithContext(ctx)
			ok      bool
			msg     string
			err     error
		)
		if ok, err = redis.CheckChallengeInit(team.ID, usage.ChallengeID); ok || err != nil {
			ctx.JSON(http.StatusTooManyRequests, gin.H{"msg": "TooQuick", "data": nil})
			return
		}
		_ = redis.RecordChallengeInit(team.ID, usage.ChallengeID)
		if !reset {
			if _, ok, msg = db.GetFlagBy3ID(DB, contest.ID, team.ID, usage.ChallengeID); ok {
				ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
				return
			}
		}
		container, ok, _ := db.GetContainerBy3ID(DB, contest.ID, team.ID, usage.ChallengeID)
		if ok {
			tx := DB.Begin()
			ok, msg = db.DeleteContainer(tx, container)
			if !ok {
				tx.Rollback()
				ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
				return
			}
			tx.Commit()
		}
		tx := DB.Begin()
		_, ok, msg = db.InitFlag(tx, contest, team, usage)
		if !ok {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
	}
}
