package router

import (
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetChallenges(ctx *gin.Context) {
	var form f.GetChallengesForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
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
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
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
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateChallengeEventType)
	tx := db.DB.WithContext(ctx).Begin()
	challenge, ok, msg := service.CreateChallenge(tx, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	if err := os.MkdirAll(challenge.BasicDir(), 0755); err != nil {
		log.Logger.Warningf("create challenge dir err: %s", err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.CreateDirError, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": resp.GetChallengeResp(challenge)})
}

func UpdateChallenge(ctx *gin.Context) {
	var form f.UpdateChallengeForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateChallengeEventType)
	challenge := middleware.GetChallenge(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.UpdateChallenge(tx, challenge, challenge.ChallengeFlags, form)
	if !ok {
		tx.Rollback()
		return
	}
	tx.Commit()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteChallenge(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteChallengeEventType)
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
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
}
