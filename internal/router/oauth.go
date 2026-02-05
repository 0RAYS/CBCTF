package router

import (
	"CBCTF/internal/config"
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
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
	oauthProviders, _, ret := db.InitOauthRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"on": true},
	})
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
	ctx.JSON(http.StatusOK, model.SuccessRetVal(data))
}

func Oauth(ctx *gin.Context) {
	uri := middleware.GetOauthUri(ctx)
	oauthProviderMapLock.RLock()
	provider, ok := oauthProviderMap[uri]
	oauthProviderMapLock.RUnlock()
	if !ok {
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Request.BadRequest})
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
		ctx.JSON(http.StatusOK, ret)
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
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Request.BadRequest})
		return
	}
	var form dto.OauthCallbackForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
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
		ctx.JSON(http.StatusOK, ret)
		return
	}
	tok, err := oauthConfig.Exchange(ctx, form.Code, oauth2.VerifierOption(verifier))
	if err != nil {
		log.Logger.Warningf("Failed to get token for provider %s: %s", provider.Provider, err)
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err}})
	}
	client := oauthConfig.Client(ctx, tok)
	response, err := client.Get(provider.UserInfoURL)
	if err != nil {
		log.Logger.Warningf("Failed to get User info by provider %s: %s", provider.Provider, err)
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err}})
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
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err}})
	}
	id, ok := utils.GetFiledValue(result, provider.RespIDField)
	if !ok {
		log.Logger.Warningf("Failed to get user_id by provider %s: %s", provider.Provider, result)
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": "Get value failed"}})
		return
	}
	userRepo := db.InitUserRepo(db.DB)
	user, ret := userRepo.Get(db.GetOptions{Conditions: map[string]any{"provider": provider.Provider, "provider_user_id": id}})
	if !ret.OK {
		if ret.Msg != i18n.Model.NotFound {
			ctx.JSON(http.StatusOK, ret)
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
		picture, _ := utils.GetFiledValue(result, provider.RespPictureField)
		description, _ := utils.GetFiledValue(result, provider.RespDescriptionField)
		raw, _ := json.Marshal(result)
		user, ret = userRepo.Create(db.CreateUserOptions{
			Name:           name,
			Password:       model.NeverLoginPWD,
			Email:          email,
			Picture:        model.FileURL(picture),
			Description:    description,
			Verified:       true,
			Provider:       provider.Provider,
			ProviderUserID: id,
			OauthRaw:       string(raw),
		})
		if !ret.OK {
			ctx.JSON(http.StatusOK, ret)
			return
		}
		prometheus.UpdateUserRegisterMetrics(provider.Provider)
	}
	ctx.Set("Self", user)
	token, err := utils.GenerateToken(user.ID, user.Name, false, model.OauthLoginType)
	if err != nil {
		log.Logger.Warningf("Failed to generate token: %s", err)
		ctx.JSON(http.StatusOK, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err}})
		return
	}
	prometheus.UpdateUserLoginMetrics(provider.Provider)
	ctx.Set(middleware.CTXEventSuccessKey, true)
	url := fmt.Sprintf("%s/platform/#/oauth/callback?token=%s", config.Env.Host, token)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func GetOauthProviders(ctx *gin.Context) {
	var form dto.ListModelsForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	oauthProviders, count, ret := db.InitOauthRepo(db.DB).List(form.Limit, form.Offset)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	data := make([]gin.H, 0)
	for _, provider := range oauthProviders {
		data = append(data, resp.GetOauthResp(provider))
	}
	ctx.JSON(http.StatusOK, model.SuccessRetVal(gin.H{"providers": data, "count": count}))
}

func GetOauthProvider(ctx *gin.Context) {
	provider := middleware.GetOauth(ctx)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(resp.GetOauthResp(provider)))
}

func CreateOauthProvider(ctx *gin.Context) {
	var form dto.CreateOauthProviderForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.CreateOauthEventType)
	provider, ret := db.InitOauthRepo(db.DB).Create(db.CreateOauthOptions{
		AuthURL:              form.AuthURL,
		TokenURL:             form.TokenURL,
		UserInfoURL:          form.UserInfoURL,
		ClientID:             form.ClientID,
		ClientSecret:         form.ClientSecret,
		Provider:             form.Provider,
		Uri:                  form.Uri,
		RespIDField:          form.RespIDField,
		RespNameField:        form.RespNameField,
		RespEmailField:       form.RespEmailField,
		RespPictureField:     form.RespPictureField,
		RespDescriptionField: form.RespDescriptionField,
		On:                   false,
	})
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal(resp.GetOauthResp(provider)))
}

func UpdateOauthProvider(ctx *gin.Context) {
	var form dto.UpdateOauthProviderForm
	if ret := form.Bind(ctx); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateOauthEventType)
	oauth := middleware.GetOauth(ctx)
	if ret := db.InitOauthRepo(db.DB).Update(oauth.ID, db.UpdateOauthOptions{
		AuthURL:              form.AuthURL,
		TokenURL:             form.TokenURL,
		UserInfoURL:          form.UserInfoURL,
		CallbackURL:          form.CallbackURL,
		ClientID:             form.ClientID,
		ClientSecret:         form.ClientSecret,
		Provider:             form.Provider,
		Uri:                  form.Uri,
		RespIDField:          form.RespIDField,
		RespNameField:        form.RespNameField,
		RespEmailField:       form.RespEmailField,
		RespPictureField:     form.RespPictureField,
		RespDescriptionField: form.RespDescriptionField,
		On:                   form.On,
	}); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	newOauth, ret := db.InitOauthRepo(db.DB).GetByID(oauth.ID)
	if !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	oauthProviderMapLock.Lock()
	if _, ok := oauthProviderMap[newOauth.Uri]; ok {
		delete(oauthProviderMap, newOauth.Uri)
	}
	if newOauth.On {
		oauthProviderMap[newOauth.Uri] = newOauth
	}
	oauthProviderMapLock.Unlock()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, ret)
}

func DeleteOauthProvider(ctx *gin.Context) {
	ctx.Set(middleware.CTXEventTypeKey, model.DeleteOauthEventType)
	oauth := middleware.GetOauth(ctx)
	if ret := db.InitOauthRepo(db.DB).Delete(oauth.ID); !ret.OK {
		ctx.JSON(http.StatusOK, ret)
		return
	}
	oauthProviderMapLock.Lock()
	if _, ok := oauthProviderMap[oauth.Uri]; ok {
		delete(oauthProviderMap, oauth.Uri)
	}
	oauthProviderMapLock.Unlock()
	ctx.Set(middleware.CTXEventSuccessKey, true)
	ctx.JSON(http.StatusOK, model.SuccessRetVal())
}
