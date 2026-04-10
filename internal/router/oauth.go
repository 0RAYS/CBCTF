package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"CBCTF/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
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
			"url":     fmt.Sprintf("%s/oauth/%s", config.Env.Host, provider.Uri),
			"name":    provider.Provider,
			"picture": provider.Picture,
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
	oauthConfig := &oauth2.Config{
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  provider.AuthURL,
			TokenURL: provider.TokenURL,
		},
		RedirectURL: provider.CallbackURL,
	}
	state := utils.UUID()
	verifier := oauth2.GenerateVerifier()
	if ret := redis.SetOauthState(provider.Provider, state, verifier); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	url := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline, oauth2.S256ChallengeOption(verifier))
	ctx.Redirect(http.StatusTemporaryRedirect, url)
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
	var form dto.OauthCallbackForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	oauthConfig := &oauth2.Config{
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  provider.AuthURL,
			TokenURL: provider.TokenURL,
		},
		RedirectURL: provider.CallbackURL,
	}
	ctx.Set(middleware.CTXEventTypeKey, model.OauthLoginEventType)
	defer redis.DelOauthState(provider.Provider, form.State)
	verifier, ret := redis.GetOauthVerifier(provider.Provider, form.State)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	tok, err := oauthConfig.Exchange(ctx, form.Code, oauth2.VerifierOption(verifier))
	if err != nil {
		log.Logger.Warningf("Failed to get token for provider %s: %s", provider.Provider, err)
		resp.JSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	client := oauthConfig.Client(ctx, tok)
	client.Timeout = time.Second * 10
	response, err := client.Get(provider.UserInfoURL)
	if err != nil {
		log.Logger.Warningf("Failed to get User info by provider %s: %s", provider.Provider, err)
		resp.JSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Logger.Warningf("Failed to close response body for provider %s: %s", provider.Provider, err)
		}
	}(response.Body)
	var result map[string]any
	if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Logger.Warningf("Failed to decode response body for provider %s: %s", provider.Provider, err)
		resp.JSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	user, ret := service.OauthLoginWithTransaction(db.DB, provider, result)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	token, err := utils.GenerateToken(user.ID, user.Name, model.OauthLoginDeviceMagic, config.Env.Gin.JWT.Secret)
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		resp.JSON(ctx, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}})
		return
	}
	ctx.Set("Self", user)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	code := utils.UUID()
	if ret := redis.SetOauthCode(code, token); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	url := fmt.Sprintf("%s/platform/#/oauth/callback?code=%s", config.Env.Host, code)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func ExchangeOauthCode(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		resp.JSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest})
		return
	}
	token, ret := redis.GetAndDelOauthToken(code)
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(gin.H{"token": token}))
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
	oauth := middleware.GetOauth(ctx)
	if ret := service.DeleteOauthProvider(db.DB, oauth); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	oauthProviderMapLock.Lock()
	if _, ok := oauthProviderMap[oauth.Uri]; ok {
		delete(oauthProviderMap, oauth.Uri)
	}
	oauthProviderMapLock.Unlock()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal())
}
