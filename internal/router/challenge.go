package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"fmt"
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
		form.Type = -1
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
	challenge := middleware.GetChallenge(ctx)
	team := middleware.GetTeam(ctx)
	path := challenge.AttachmentPath(team.ID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": "FileNotFound", "data": nil})
		return
	}
	ctx.File(path)
}

func GetChallengeFiles(ctx *gin.Context) {
	challenge := middleware.GetChallenge(ctx)
	var files []string
	if middleware.GetRole(ctx) == "admin" {
		dir, err := os.ReadDir(challenge.BasicDir())
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"msg": "ReadDirError", "data": nil})
			return
		}
		for _, file := range dir {
			files = append(files, file.Name())
		}
	} else {
		team := middleware.GetTeam(ctx)
		if _, err := os.Stat(challenge.AttachmentPath(team.ID)); err == nil {
			files = append(files, model.DynamicFile)
		}
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
	if t, ok := data["type"]; ok && !db.IsValidChallengeType(t.(int)) {
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

func UploadChallenge(ctx *gin.Context) {
	challenge := middleware.GetChallenge(ctx)
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	var path string
	switch challenge.Type {
	case model.Static, model.Container:
		if file.Filename != model.StaticFile {
			ctx.JSON(http.StatusOK, gin.H{"msg": "InvalidFileName", "data": nil})
			return
		}
		path = fmt.Sprintf("%s/%s", challenge.BasicDir(), model.StaticFile)
	case model.Dynamic:
		if file.Filename != model.DynamicFile {
			ctx.JSON(http.StatusOK, gin.H{"msg": "InvalidFileName", "data": nil})
			return
		}
		path = fmt.Sprintf("%s/%s", challenge.BasicDir(), model.DynamicFile)
	default:
		ctx.JSON(http.StatusOK, gin.H{"msg": "InvalidChallengeType", "data": nil})
		return
	}
	if err := ctx.SaveUploadedFile(file, path); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "UnknownError", "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": nil})
}

func DownloadChallenge(ctx *gin.Context) {
	var form f.DownloadChallengeForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	challenge := middleware.GetChallenge(ctx)
	var path string
	switch form.File {
	case model.StaticFile, model.DynamicFile:
		path = fmt.Sprintf("%s/%s", challenge.BasicDir(), form.File)
	default:
		ctx.JSON(http.StatusOK, gin.H{"msg": "InvalidFileName", "data": nil})
		return

	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": "FileNotFound", "data": nil})
		return
	}
	ctx.File(path)
}

func ChallengeStatus(ctx *gin.Context) {
	var (
		team model.Team
		ok   bool
		msg  string
	)
	team = middleware.GetTeam(ctx)
	_, ok, msg = db.GetFlagBy3ID(db.DB.WithContext(ctx), middleware.GetContest(ctx).ID, team.ID, middleware.GetChallenge(ctx).ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": false})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": true})

}

func InitChallenge(reset bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			team    model.Team
			contest model.Contest
			usage   model.Usage
			ok      bool
			msg     string
			err     error
			DB      = db.DB.WithContext(ctx)
		)
		team = middleware.GetTeam(ctx)
		if ok, err = redis.CheckChallengeInit(team.ID, middleware.GetChallenge(ctx).ID); ok || err != nil {
			ctx.JSON(http.StatusTooManyRequests, gin.H{"msg": "TooQuick", "data": nil})
			return
		}
		_ = redis.RecordChallengeInit(team.ID, middleware.GetChallenge(ctx).ID)
		contest = middleware.GetContest(ctx)
		usage, ok, msg = db.GetUsageBy2ID(DB, contest.ID, middleware.GetChallenge(ctx).ID)
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		if !reset {
			if _, ok, msg = db.GetFlagBy3ID(DB, contest.ID, team.ID, usage.ChallengeID); ok {
				ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
				return
			}
		}
		docker, ok, _ := db.GetDockerBy3ID(DB, contest.ID, team.ID, usage.ChallengeID)
		if ok {
			tx := DB.Begin()
			ok, msg = db.DeleteDocker(tx, docker)
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
