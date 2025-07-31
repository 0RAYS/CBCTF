package router

import (
	"CBCTF/internal/config"
	f "CBCTF/internal/form"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	db "CBCTF/internal/repo"
	"CBCTF/internal/resp"
	"CBCTF/internal/service"
	"CBCTF/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"strings"
	"sync"
)

var (
	AvailableOauthProviders = make(map[string]gin.H)
	OauthProvidersMutex     = sync.RWMutex{}
)

func GetOauthProviders(ctx *gin.Context) {
	var form f.GetModelsForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	oauthProviders, count, ok, msg := db.InitOauthRepo(db.DB.WithContext(ctx)).List(form.Limit, form.Offset)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"providers": oauthProviders, "count": count}})
}

func CreateOauthProvider(ctx *gin.Context) {
	var form f.CreateOauthProviderForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx := db.DB.WithContext(ctx).Begin()
	provider, ok, msg := service.CreateOauthProvider(tx, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": resp.GetOauthResp(provider)})
}

func UpdateOauthProvider(ctx *gin.Context) {
	var form f.UpdateOauthProviderForm
	if ok, msg := form.Bind(ctx); !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	oauth := middleware.GetOauth(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	ok, msg := service.UpdateOauthProvider(tx, oauth, form)
	if !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	if form.On != nil && !oauth.On && *form.On {
		oauth, ok, msg = db.InitOauthRepo(db.DB.WithContext(ctx)).Get(db.GetOptions{Conditions: map[string]any{"id": oauth.ID}})
		if !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		registerOauthRoutes(router, oauth)
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
}

func DeleteOauthProvider(ctx *gin.Context) {
	oauth := middleware.GetOauth(ctx)
	tx := db.DB.WithContext(ctx).Begin()
	if ok, msg := db.InitOauthRepo(tx).Delete(oauth.ID); !ok {
		tx.Rollback()
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
		return
	}
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": nil})
}

func registerOauthRoutes(router *gin.Engine, provider model.Oauth) {
	oauthConfig := &oauth2.Config{
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  provider.AuthURL,
			TokenURL: provider.TokenURL,
		},
		RedirectURL: fmt.Sprintf("%s/oauth/%s/callback", config.Env.Backend, provider.URI),
	}

	for _, r := range router.Routes() {
		if (r.Path == fmt.Sprintf("/oauth/%s", provider.URI) && r.Method == http.MethodPost) || (r.Path == fmt.Sprintf("/oauth/%s/callback", provider.URI) && r.Method == http.MethodGet) {
			log.Logger.Infof("Route for provider %s already exists: %s %s", provider.Provider, r.Method, r.Path)
			return
		}
	}

	router.POST(fmt.Sprintf("/oauth/%s", provider.URI), func(ctx *gin.Context) {
		state := utils.UUID()
		verifier := oauth2.GenerateVerifier()
		if err := redis.SetOauthState(provider.Provider, state, verifier); err != nil {
			log.Logger.Warningf("Failed to set oauth state for provider %s: %v", provider.Provider, err)
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.RedisError, "data": nil})
			return
		}
		url := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline, oauth2.S256ChallengeOption(verifier))
		ctx.Redirect(http.StatusTemporaryRedirect, url)
	})

	router.GET(fmt.Sprintf("/oauth/%s/callback", provider.URI), func(ctx *gin.Context) {
		var form f.OauthCallbackForm
		if ok, msg := form.Bind(ctx); !ok {
			ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
			return
		}
		verifier, err := redis.GetOauthVerifier(provider.Provider, form.State)
		if err != nil {
			log.Logger.Warningf("Failed to get oauth verifier for provider %s: %v", provider.Provider, err)
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.RedisError, "data": nil})
			return
		}
		tok, err := oauthConfig.Exchange(ctx, form.Code, oauth2.VerifierOption(verifier))
		if err != nil {
			log.Logger.Warningf("Failed to get token for provider %s: %v", provider.Provider, err)
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
		}
		client := oauthConfig.Client(ctx, tok)
		response, err := client.Get(provider.UserInfoURL)
		if err != nil {
			log.Logger.Warningf("Failed to get user info for provider %s: %v", provider.Provider, err)
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
			return
		}
		defer func(Body io.ReadCloser) {
			if err = Body.Close(); err != nil {
				log.Logger.Warningf("Failed to close response body for provider %s: %v", provider.Provider, err)
			}
		}(response.Body)
		var result map[string]any
		if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
			log.Logger.Warningf("Failed to decode response body for provider %s: %v", provider.Provider, err)
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
		}
		id, ok := result[provider.RespIDField].(string)
		if !ok {
			log.Logger.Warningf("Failed to provider user_id for provider %s: %v", provider.Provider, result)
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
			return
		}
		userRepo := db.InitUserRepo(db.DB.WithContext(ctx))
		user, ok, msg := userRepo.Get(db.GetOptions{Conditions: map[string]any{"provider": provider.Provider, "provider_id": id}})
		if ok {
			if msg != i18n.UserNotFound {
				ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
				return
			}
			name, ok := utils.GetFiledValue(result, provider.RespNameField)
			if !ok {
				name = fmt.Sprintf("%s_%s", provider.Provider, utils.RandStr(10))
			} else {
				name = fmt.Sprintf("%s_%s", name, utils.RandStr(5))
			}
			email, ok := utils.GetFiledValue(result, provider.RespEmailField)
			if !ok {
				email = fmt.Sprintf("%s_%s@example.com", provider.Provider, utils.RandStr(10))
			} else {
				email = fmt.Sprintf("%s_%s", utils.RandStr(10), email)
			}
			avatar, _ := utils.GetFiledValue(result, provider.RespAvatarField)
			desc, _ := utils.GetFiledValue(result, provider.RespDescField)
			raw, _ := json.Marshal(result)
			user, ok, msg = userRepo.Create(db.CreateUserOptions{
				Name:           name,
				Password:       "never_login_pwd",
				Email:          email,
				Avatar:         model.AvatarURL(avatar),
				Desc:           desc,
				Verified:       true,
				Provider:       provider.Provider,
				ProviderUserID: id,
				OauthRaw:       string(raw),
			})
			if !ok {
				log.Logger.Warningf("Failed to create user for provider %s: %s", provider.Provider, msg)
				ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": nil})
				return
			}
		}
		token, err := utils.GenerateToken(user.ID, user.Name, false, middleware.GetMagic(ctx))
		if err != nil {
			log.Logger.Warningf("Failed to generate token: %s", err)
			ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
			return
		}
		ctx.Set("IsAdmin", false)
		ctx.Set("Self", user)
		log.Logger.Infof("%s:%d login", user.Name, user.ID)
		ctx.Writer.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
		ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": resp.LoginResp(user, false)})
	})
	OauthProvidersMutex.Lock()
	AvailableOauthProviders[provider.Provider] = gin.H{
		"url":    fmt.Sprintf("%s/oauth/%s", config.Env.Backend, provider.URI),
		"name":   provider.Provider,
		"avatar": provider.Avatar,
	}
	OauthProvidersMutex.Unlock()
}

func RegisterOauthRouter(router *gin.Engine) {
	oauthProviders, _, ok, _ := db.InitOauthRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"on": true},
	})
	if !ok || len(oauthProviders) == 0 {
		return
	}
	for _, provider := range oauthProviders {
		if strings.HasPrefix(provider.Provider, "http://") || strings.HasPrefix(provider.Provider, "https://") {
			continue
		}
		registerOauthRoutes(router, provider)
	}
	router.GET("/oauth", func(ctx *gin.Context) {
		OauthProvidersMutex.RLock()
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": AvailableOauthProviders})
		OauthProvidersMutex.RUnlock()
	})
}
