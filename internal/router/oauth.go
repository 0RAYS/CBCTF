package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"CBCTF/internal/redis"
	"CBCTF/internal/resp"
	"CBCTF/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

var (
	// oauthProviderMap[model.Oauth.Uri] = model.Oauth
	oauthProviderMap     = make(map[string]model.Oauth)
	oauthProviderMapLock sync.RWMutex
)

func RegisterOauthRouter() {
	oauthProviders, _, ok, _ := db.InitOauthRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"on": true},
	})
	if !ok {
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
			"url":    fmt.Sprintf("%s/oauth/%s", config.Env.Backend, provider.Uri),
			"name":   provider.Provider,
			"avatar": provider.Avatar,
		})
	}
	oauthProviderMapLock.RUnlock()
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": data})
}

func Oauth(ctx *gin.Context) {
	uri := middleware.GetOauthUri(ctx)
	oauthProviderMapLock.RLock()
	provider, ok := oauthProviderMap[uri]
	oauthProviderMapLock.RUnlock()
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
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
	if err := redis.SetOauthState(provider.Provider, state, verifier); err != nil {
		log.Logger.Warningf("Failed to set oauth state for provider %s: %s", provider.Provider, err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.RedisError, "data": nil})
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
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.BadRequest, "data": nil})
		return
	}
	var form f.OauthCallbackForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
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
	defer func(provider string, state string) {
		if err := redis.DelOauthState(provider, state); err != nil {
			log.Logger.Warningf("Failed to delete oauth state for provider %s: %s", provider, err)
		}
	}(provider.Provider, form.State)
	verifier, err := redis.GetOauthVerifier(provider.Provider, form.State)
	if err != nil {
		log.Logger.Warningf("Failed to get oauth verifier for provider %s: %s", provider.Provider, err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.RedisError, "data": nil})
		return
	}
	tok, err := oauthConfig.Exchange(ctx, form.Code, oauth2.VerifierOption(verifier))
	if err != nil {
		log.Logger.Warningf("Failed to get token for provider %s: %s", provider.Provider, err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
	}
	client := oauthConfig.Client(ctx, tok)
	response, err := client.Get(provider.UserInfoURL)
	if err != nil {
		log.Logger.Warningf("Failed to get User info by provider %s: %s", provider.Provider, err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
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
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
	}
	id, ok := utils.GetFiledValue(result, provider.RespIDField)
	if !ok {
		log.Logger.Warningf("Failed to get user_id by provider %s: %s", provider.Provider, result)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
		return
	}
	userRepo := db.InitUserRepo(db.DB)
	user, ok, msg := userRepo.Get(db.GetOptions{Conditions: map[string]any{"provider": provider.Provider, "provider_user_id": id}})
	if !ok {
		if msg != i18n.UserNotFound {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		name, ok := utils.GetFiledValue(result, provider.RespNameField)
		if !ok {
			name = fmt.Sprintf("%s_%s", provider.Provider, utils.RandStr(10))
		}
		email, ok := utils.GetFiledValue(result, provider.RespEmailField)
		if !ok {
			email = fmt.Sprintf("%s_%s@example.com", provider.Provider, utils.RandStr(10))
		}
		avatar, _ := utils.GetFiledValue(result, provider.RespAvatarField)
		desc, _ := utils.GetFiledValue(result, provider.RespDescField)
		raw, _ := json.Marshal(result)
		user, ok, msg = userRepo.Create(db.CreateUserOptions{
			Name:           name,
			Password:       model.NeverLoginPWD,
			Email:          email,
			Avatar:         model.AvatarURL(avatar),
			Desc:           desc,
			Verified:       true,
			Provider:       provider.Provider,
			ProviderUserID: id,
			OauthRaw:       string(raw),
		})
		if !ok {
			log.Logger.Warningf("Failed to create User by provider %s: %s", provider.Provider, msg)
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		prometheus.UpdateUserRegisterMetrics(provider.Provider)
	}
	ctx.Set("Self", user)
	token, err := utils.GenerateToken(user.ID, user.Name, false, model.OauthLoginType)
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
		return
	}
	prometheus.UpdateUserLoginMetrics(provider.Provider)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/platform/#/oauth/callback?token=%s", config.Env.Frontend, token))
}

func GetOauthProviders(ctx *gin.Context) {
	var form f.GetModelsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	oauthProviders, count, ok, msg := db.InitOauthRepo(db.DB).List(form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	data := make([]gin.H, 0)
	for _, provider := range oauthProviders {
		data = append(data, resp.GetOauthResp(provider))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"providers": data, "count": count}})
}

func GetOauthProvider(ctx *gin.Context) {
	provider := middleware.GetOauth(ctx)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": resp.GetOauthResp(provider)})
}

func CreateOauthProvider(ctx *gin.Context) {
	var form f.CreateOauthProviderForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateOauthEventType)
	provider, ok, msg := db.InitOauthRepo(db.DB).Create(db.CreateOauthOptions{
		AuthURL:         form.AuthURL,
		TokenURL:        form.TokenURL,
		UserInfoURL:     form.UserInfoURL,
		ClientID:        form.ClientID,
		ClientSecret:    form.ClientSecret,
		Provider:        form.Provider,
		Uri:             form.Uri,
		RespIDField:     form.RespIDField,
		RespNameField:   form.RespNameField,
		RespEmailField:  form.RespEmailField,
		RespAvatarField: form.RespAvatarField,
		RespDescField:   form.RespDescField,
		On:              false,
	})
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": resp.GetOauthResp(provider)})
}

func UpdateOauthProvider(ctx *gin.Context) {
	var form f.UpdateOauthProviderForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateOauthEventType)
	oauth := middleware.GetOauth(ctx)
	if ok, msg := db.InitOauthRepo(db.DB).Update(oauth.ID, db.UpdateOauthOptions{
		AuthURL:         form.AuthURL,
		TokenURL:        form.TokenURL,
		UserInfoURL:     form.UserInfoURL,
		CallbackURL:     form.CallbackURL,
		ClientID:        form.ClientID,
		ClientSecret:    form.ClientSecret,
		Provider:        form.Provider,
		Uri:             form.Uri,
		RespIDField:     form.RespIDField,
		RespNameField:   form.RespNameField,
		RespEmailField:  form.RespEmailField,
		RespAvatarField: form.RespAvatarField,
		RespDescField:   form.RespDescField,
		On:              form.On,
	}); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	newOauth, ok, msg := db.InitOauthRepo(db.DB).GetByID(oauth.ID)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	oauthProviderMapLock.Lock()
	if _, ok = oauthProviderMap[newOauth.Uri]; ok {
		delete(oauthProviderMap, newOauth.Uri)
	}
	if newOauth.On {
		oauthProviderMap[newOauth.Uri] = newOauth
	}
	oauthProviderMapLock.Unlock()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteOauthProvider(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteOauthEventType)
	oauth := middleware.GetOauth(ctx)
	if ok, msg := db.InitOauthRepo(db.DB).Delete(oauth.ID); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	oauthProviderMapLock.Lock()
	if _, ok := oauthProviderMap[oauth.Uri]; ok {
		delete(oauthProviderMap, oauth.Uri)
	}
	oauthProviderMapLock.Unlock()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
}
