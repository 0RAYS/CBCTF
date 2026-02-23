package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/oauth"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"CBCTF/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	if !config.Env.Registration.Enabled {
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Request.Forbidden})
		return
	}
	var form dto.RegisterForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.RegisterEventType)
	tx := db.DB.Begin()
	user, ret := service.CreateUser(tx, form)
	if !ret.OK {
		tx.Rollback()
		ctx.JSON(http.StatusOK, ret)
		return
	}
	if config.Env.Registration.DefaultGroup != 0 {
		if defaultGroup, ret := db.InitGroupRepo(tx).GetByID(config.Env.Registration.DefaultGroup); ret.OK {
			if ret = db.AppendUserToGroup(tx, user, defaultGroup); !ret.OK {
				tx.Rollback()
				ctx.JSON(http.StatusOK, ret)
				return
			}
		}
	}
	tx.Commit()
	if ret = service.SendEmail(user); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	token, err := utils.GenerateToken(user.ID, user.Name, middleware.GetMagic(ctx))
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	ctx.Set("Self", user)
	log.Logger.Infof("%s:%d register", user.Name, user.ID)
	ctx.Writer.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	prometheus.UpdateUserRegisterMetrics(oauth.LocalProvider)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(resp.GetUserResp(user, false)))
}

func Login(ctx *gin.Context) {
	var form dto.LoginForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.LoginEventType)
	user, ret := service.VerifyUser(db.DB, form)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	token, err := utils.GenerateToken(user.ID, user.Name, middleware.GetMagic(ctx))
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	ctx.Set("Self", user)
	log.Logger.Infof("%s:%d login", user.Name, user.ID)
	ctx.Writer.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	prometheus.UpdateUserLoginMetrics(oauth.LocalProvider)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(resp.GetUserResp(user, false)))
}
