package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetEmails(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	options := db.GetOptions{}
	smtp := middleware.GetSmtp(ctx)
	if smtp.ID > 0 {
		options.Conditions = map[string]any{"smtp_id": smtp.ID}
	}
	emails, count, ret := db.InitEmailRepo(db.DB).List(form.Limit, form.Offset, options)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, email := range emails {
		data = append(data, resp.GetEmailResp(email))
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"emails": data, "count": count}))
}

func VerifyEmail(ctx *gin.Context) {
	var form dto.VerifyEmail
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.VerifyEmailEventType)
	ret := service.VerifyEmail(db.DB, form)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.Redirect(http.StatusTemporaryRedirect, config.Env.Host)
}

func ActivateEmail(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.ActivateEmailEventType)
	user := middleware.GetSelf(ctx).(model.User)
	if user.Verified {
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Model.User.AlreadyVerified})
		return
	}
	ret := service.SendEmail(user)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, ret)
}
