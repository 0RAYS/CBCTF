package router

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"CBCTF/internel/middleware"
	db "CBCTF/internel/repo"
	"CBCTF/internel/resp"
	"CBCTF/internel/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func GetChallenges(ctx *gin.Context) {
	var form f.GetChallengesForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	if _, exists := ctx.GetQuery("limit"); !exists {
		form.Limit = 5
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		form.Offset = 0
	}
	challenges, count, ok, msg := service.GetChallenges(db.DB.WithContext(ctx), form)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, challenge := range challenges {
		data = append(data, resp.GetChallengeResp(challenge))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"count": count, "challenges": data}})
}

func GetChallenge(ctx *gin.Context) {
	challenge := middleware.GetChallenge(ctx)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": resp.GetChallengeResp(challenge)})
}

func GetCategories(ctx *gin.Context) {
	var form f.GetCategoriesForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	repo := db.InitChallengeRepo(db.DB.WithContext(ctx))
	categories, ok, msg := repo.ListCategories(form.Type)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": categories})
}

func GetChallengeFiles(ctx *gin.Context) {
	challenge := middleware.GetChallenge(ctx)
	dir, err := os.ReadDir(challenge.BasicDir())
	if err != nil {
		log.Logger.Warningf("Failed to read challenge directory: %s", err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.ReadDirError, "data": nil})
		return
	}
	files := make([]string, 0)
	for _, file := range dir {
		files = append(files, file.Name())
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": files})
}

func CreateChallenge(ctx *gin.Context) {
	var form f.CreateChallengeForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	challenge, ok, msg := service.CreateChallenge(tx, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	if err := os.MkdirAll(challenge.BasicDir(), 0755); err != nil {
		log.Logger.Warningf("create challenge dir err: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.CreateDirError, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": resp.GetChallengeResp(challenge)})
}

func UpdateChallenge(ctx *gin.Context) {
	var form f.UpdateChallengeForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.UpdateChallenge(tx, middleware.GetChallenge(ctx), form)
	if !ok {
		tx.Rollback()
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteChallenge(ctx *gin.Context) {
	challenge := middleware.GetChallenge(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := db.InitChallengeRepo(tx).Delete(challenge.RandID)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	if err := os.RemoveAll(challenge.BasicDir()); err != nil {
		log.Logger.Warningf("Failed to remove challenge basic dir: %s", err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
}
