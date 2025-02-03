package router

import (
	"CBCTF/internal/constants"
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

func CreateChallenge(ctx *gin.Context) {
	var form constants.CreateChallengeForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	form.Category = utils.ToTitle(strings.TrimSpace(form.Category))
	challenge, ok, msg := db.CreateChallenge(ctx, form)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if err := os.MkdirAll(challenge.BasicDir(), 0755); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": "CreateDirError", "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": challenge})
}

func GetChallenge(ctx *gin.Context) {
	contest, ok, msg := db.GetChallengeByID(ctx, middleware.GetChallengeID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": contest})
}

func GetChallenges(ctx *gin.Context) {
	var form constants.GetChallengesForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	if ctx.Query("type") == "" && ctx.Query("category") == "" {
		form.Type = -1
		form.Category = ""
	}
	challenges, count, ok, msg := db.GetChallenges(ctx, form.Limit, form.Offset, form.Type, form.Category)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "challenges": challenges}})
}

func GetAttachment(ctx *gin.Context) {
	challenge, ok, msg := db.GetChallengeByID(ctx, middleware.GetChallengeID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	team, ok, msg := db.GetTeamByUserID(ctx, middleware.GetSelfID(ctx), middleware.GetContestID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	path := challenge.AttachmentPath(team.ID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": "FileNotFound", "data": nil})
		return
	}
	ctx.File(path)
}

func GetChallengeFiles(ctx *gin.Context) {
	challenge, ok, msg := db.GetChallengeByID(ctx, middleware.GetChallengeID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	dir, err := os.ReadDir(challenge.BasicDir())
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": "ReadDirError", "data": nil})
		return
	}
	var files []string
	for _, file := range dir {
		files = append(files, file.Name())
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": files})
}

func GetCategories(ctx *gin.Context) {
	var form constants.GetCategoriesForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	categories, ok, msg := db.GetCategories(ctx, form.Type)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "Success", "data": categories})
}

func UpdateChallenge(ctx *gin.Context) {
	var (
		challenge model.Challenge
		ok        bool
		msg       string
		data      map[string]interface{}
	)
	var form constants.UpdateChallengeForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	challenge, ok, msg = db.GetChallengeByID(ctx, middleware.GetChallengeID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
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
	_, msg = db.UpdateChallenge(ctx, challenge.ID, data)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteChallenge(ctx *gin.Context) {
	var form constants.DeleteChallengeForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	challenge, ok, msg := db.GetChallengeByID(ctx, middleware.GetChallengeID(ctx))
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": msg, "data": nil})
		return
	}
	usages, ok, msg := db.GetUsageByChallengeID(ctx, challenge.ID)
	if ok {
		for _, usage := range usages {
			db.DeleteUsage(ctx, usage.ID)
		}
	}
	_, msg = db.DeleteChallenge(ctx, middleware.GetChallengeID(ctx))
	if form.Force && os.RemoveAll(challenge.BasicDir()) != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"msg": "UnknownError", "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func UploadChallenge(ctx *gin.Context) {
	challenge, ok, msg := db.GetChallengeByID(ctx, middleware.GetChallengeID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
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
	var form constants.DownloadChallengeForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	challenge, ok, msg := db.GetChallengeByID(ctx, middleware.GetChallengeID(ctx))
	if !ok {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": msg, "data": nil})
		return
	}
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
	team, ok, msg = db.GetTeamByUserID(ctx, middleware.GetSelfID(ctx), middleware.GetContestID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	_, ok, msg = db.GetFlagBy3ID(ctx, middleware.GetContestID(ctx), team.ID, middleware.GetChallengeID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": false})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": true})

}

func InitChallenge(ctx *gin.Context) {
	var (
		team    model.Team
		contest model.Contest
		usage   model.Usage
		ok      bool
		msg     string
	)
	team, ok, msg = db.GetTeamByUserID(ctx, middleware.GetSelfID(ctx), middleware.GetContestID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	contest, ok, msg = db.GetContestByID(ctx, middleware.GetContestID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if !contest.IsRunning() {
		ctx.JSON(http.StatusOK, gin.H{"msg": contest.Status(), "data": nil})
		return
	}
	usage, ok, msg = db.GetUsageBy2ID(ctx, contest.ID, middleware.GetChallengeID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	if _, ok, msg = db.GetFlagBy3ID(ctx, contest.ID, team.ID, usage.ChallengeID); ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	_, ok, msg = db.InitFlag(ctx, contest, team, usage)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
