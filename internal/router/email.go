package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetEmails(ctx *gin.Context) {
	var form f.GetModelsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	options := db.GetOptions{}
	smtp := middleware.GetSmtp(ctx)
	if smtp.ID > 0 {
		options.Conditions = map[string]any{"smtp_id": smtp.ID}
	}
	emails, count, ok, msg := db.InitEmailRepo(db.DB).List(form.Limit, form.Offset, options)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, email := range emails {
		data = append(data, resp.GetEmailResp(email))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"emails": data, "count": count}})
}

func VerifyEmail(ctx *gin.Context) {
	var form f.VerifyEmail
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.VerifyEmailEventType)
	tx := db.DB.Begin()
	ok, msg := service.VerifyEmail(tx, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.Redirect(http.StatusTemporaryRedirect, config.Env.Backend)
}

func ActivateEmail(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.ActivateEmailEventType)
	user := middleware.GetSelf(ctx).(model.User)
	if user.Verified {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.AlreadyVerifiedEmail, "data": nil})
		return
	}
	_, msg := service.SendEmail(user)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}
