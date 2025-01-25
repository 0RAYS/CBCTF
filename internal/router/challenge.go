package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/constants"
	"CBCTF/internal/db"
	"CBCTF/internal/log"
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
	form.Category = strings.ToTitle(strings.TrimSpace(form.Category))
	challenge := model.InitChallenge(form)
	base := fmt.Sprintf("%s/challenges/%s", config.Env.Gin.Upload.Path, challenge.Path)
	var failed []string
	for _, p := range []string{"attachment", "mounted", "generator"} {
		path := fmt.Sprintf("%s/%s", base, p)
		if err := os.MkdirAll(path, 0755); err != nil {
			failed = append(failed, path)
			log.Logger.Errorf("Create path %s failed: %s", path, err)
		}
	}
	if len(failed) > 0 {
		ctx.JSON(http.StatusOK, gin.H{"msg": "CreateDirError", "data": failed})
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
	var form constants.GetModelsForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": "BadRequest", "data": nil})
		return
	}
	challenges, count, ok, msg := db.GetChallenges(ctx, form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "challenges": challenges}})
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
	tmp := strings.ToTitle(strings.TrimSpace(*form.Category))
	form.Category = &tmp
	data = utils.Form2Map(form)
	if t, ok := data["type"]; ok && !db.IsValidChallengeType(t.(uint)) {
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

}
