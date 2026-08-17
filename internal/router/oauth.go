package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/oa"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/redis"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"CBCTF/internal/utils"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	// oauthProviderMap[model.Oauth.Uri] = model.Oauth
	oauthProviderMap     = make(map[string]model.Oauth)
	oauthProviderMapLock sync.RWMutex
)

func RegisterOauthRouter() {
	oauthProviders, ret := service.ListEnabledOauthProviders(db.DB)
	if !ret.OK {
		return
	}
	oauthProviderMapLock.Lock()
	for _, provider := range oauthProviders {
		oauthProviderMap[provider.Uri] = provider
	}
	oauthProviderMapLock.Unlock()
}

func ListOauth(ctx *gin.Context) {
	data := make([]gin.H, 0)
	oauthProviderMapLock.RLock()
	for _, provider := range oauthProviderMap {
		data = append(data, gin.H{
			"url":      fmt.Sprintf("%s/oauth/%s", config.Env.Host, provider.Uri),
			"name":     provider.Provider,
			"protocol": provider.Protocol,
			"picture":  provider.Picture,
		})
	}
	oauthProviderMapLock.RUnlock()
	resp.JSON(ctx, model.SuccessRetVal(data))
}

func Oauth(ctx *gin.Context) {
	uri := middleware.GetOauthUri(ctx)
	oauthProviderMapLock.RLock()
	provider, ok := oauthProviderMap[uri]
	oauthProviderMapLock.RUnlock()
	if !ok {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	protocol, ok := oa.GetProtocol(provider.Protocol)
	if !ok {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	loginURL, ret := protocol.LoginURL(provider)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Redirect(http.StatusTemporaryRedirect, loginURL)
}

func OauthCallback(ctx *gin.Context) {
	uri := middleware.GetOauthUri(ctx)
	oauthProviderMapLock.RLock()
	provider, ok := oauthProviderMap[uri]
	oauthProviderMapLock.RUnlock()
	if !ok {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	protocol, ok := oa.GetProtocol(provider.Protocol)
	if !ok {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.OauthLoginEventType)
	result, ret := protocol.Exchange(ctx, provider)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	user, ret := service.OauthLoginWithTransaction(db.DB, provider, result)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	token, err := utils.GenerateToken(user.ID, user.Name, config.Env.Gin.JWT.Secret)
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		resp.JSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	ctx.Set("Self", user)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	code := utils.UUID()
	if ret = redis.SetOauthCode(code, token); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	prometheus.RecordUserLogin(provider.Provider)
	url := fmt.Sprintf("%s/platform/#/oauth/callback?code=%s", config.Env.Host, code)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func ExchangeOauthCode(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.OauthLoginEventType)
	tempToken, ret := redis.GetAndDelOauthToken(code)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	claims, err := utils.ParseToken(tempToken, config.Env.Gin.JWT.Secret)
	if err != nil {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Response.Unauthorized})
		return
	}
	token, err := utils.GenerateToken(claims.UserID, claims.Name, config.Env.Gin.JWT.Secret)
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		resp.JSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	setAuthCookie(ctx, token)
	ctx.Set(middleware.CTXEventModelsKey, model.UintMap{"Self": claims.UserID})
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}

func GetOauthProviders(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	oauthProviders, count, ret := service.ListOauthProviders(db.DB, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, provider := range oauthProviders {
		data = append(data, resp.GetOauthResp(provider))
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"providers": data, "count": count}))
}

func GetOauthProvider(ctx *gin.Context) {
	resp.JSON(ctx, model.SuccessRetVal(resp.GetOauthResp(middleware.GetOauth(ctx))))
}

func CreateOauthProvider(ctx *gin.Context) {
	var form dto.CreateOauthProviderForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateOauthEventType)
	provider, ret := service.CreateOauthProvider(db.DB, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetOauthResp(provider)))
}

func UpdateOauthProvider(ctx *gin.Context) {
	var form dto.UpdateOauthProviderForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateOauthEventType)
	oldOauth := middleware.GetOauth(ctx)
	newOauth, ret := service.UpdateOauthProvider(db.DB, oldOauth, form)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	oauthProviderMapLock.Lock()
	if _, ok := oauthProviderMap[oldOauth.Uri]; ok {
		delete(oauthProviderMap, oldOauth.Uri)
	}
	if newOauth.On {
		oauthProviderMap[newOauth.Uri] = newOauth
	}
	oauthProviderMapLock.Unlock()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, ret)
}

func DeleteOauthProvider(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteOauthEventType)
	provider := middleware.GetOauth(ctx)
	if ret := service.DeleteOauthProvider(db.DB, provider); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	oauthProviderMapLock.Lock()
	if _, ok := oauthProviderMap[provider.Uri]; ok {
		delete(oauthProviderMap, provider.Uri)
	}
	oauthProviderMapLock.Unlock()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}
