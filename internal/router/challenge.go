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
	if err := os.MkdirAll(challenge.Path, 0755); err != nil {
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
	challenges, count, ok, msg := db.GetChallenges(ctx, form.Limit, form.Offset, form.Type, form.Category)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "challenges": challenges}})
}

func GetChallengeFiles(ctx *gin.Context) {
	challenge, ok, msg := db.GetChallengeByID(ctx, middleware.GetChallengeID(ctx))
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	dir, err := os.ReadDir(challenge.Path)
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
	_, msg := db.DeleteChallenge(ctx, middleware.GetChallengeID(ctx))
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
	case model.Static:
		if file.Filename != model.StaticFile {
			ctx.JSON(http.StatusOK, gin.H{"msg": "InvalidFileName", "data": nil})
			return
		}
		path = fmt.Sprintf("%s/%s", challenge.Path, model.StaticFile)
	case model.Dynamic:
		if file.Filename != model.DynamicFile {
			ctx.JSON(http.StatusOK, gin.H{"msg": "InvalidFileName", "data": nil})
			return
		}
		path = fmt.Sprintf("%s/%s", challenge.Path, model.DynamicFile)
	case model.Container:
		if file.Filename != model.ContainerFile {
			ctx.JSON(http.StatusOK, gin.H{"msg": "InvalidFileName", "data": nil})
			return
		}
		path = fmt.Sprintf("%s/%s", challenge.Path, model.ContainerFile)
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
	path := fmt.Sprintf("%s/%s", challenge.Path, form.File)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"msg": "FileNotFound", "data": nil})
		return
	}
	ctx.File(path)
}
