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
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	emails, count, ret := service.ListEmails(db.DB, middleware.GetSmtp(ctx), form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, email := range emails {
		data = append(data, resp.GetEmailResp(email))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"emails": data, "count": count}))
}

func VerifyEmail(ctx *gin.Context) {
	var form dto.VerifyEmail
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.VerifyEmailEventType)
	if ret := service.VerifyEmail(db.DB, form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.Redirect(http.StatusTemporaryRedirect, config.Env.Host)
}

func ActivateEmail(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.ActivateEmailEventType)
	user := middleware.GetSelf(ctx)
	if user.Verified {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Model.User.AlreadyVerified})
		return
	}
	ret := service.SendEmail(user)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, ret)
}
