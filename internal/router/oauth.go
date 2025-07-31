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
	data := make([]gin.H, 0)
	for _, provider := range oauthProviders {
		data = append(data, resp.GetOauthResp(provider))
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": msg, "data": gin.H{"providers": data, "count": count}})
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

func RegisterOauthRouter(router *gin.Engine) {
	oauthProviders, _, ok, _ := db.InitOauthRepo(db.DB).List(-1, -1, db.GetOptions{
		Conditions: map[string]any{"on": true},
	})
	if !ok || len(oauthProviders) == 0 {
		return
	}
	availableOauth := make(gin.H)
	for _, provider := range oauthProviders {
		if strings.HasPrefix(provider.Provider, "http://") || strings.HasPrefix(provider.Provider, "https://") {
			continue
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

		router.GET(fmt.Sprintf("/oauth/%s", provider.URI), func(ctx *gin.Context) {
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
			defer func(provider string, state string) {
				if err := redis.DelOauthState(provider, state); err != nil {
					log.Logger.Warningf("Failed to delete oauth state for provider %s: %v", provider, err)
				}
			}(provider.Provider, form.State)
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
			id, ok := utils.GetFiledValue(result, provider.RespIDField)
			if !ok {
				log.Logger.Warningf("Failed to provider user_id for provider %s: %v", provider.Provider, result)
				ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
				return
			}
			userRepo := db.InitUserRepo(db.DB.WithContext(ctx))
			user, ok, msg := userRepo.Get(db.GetOptions{Conditions: map[string]any{"provider": provider.Provider, "provider_user_id": id}})
			if !ok {
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
			token, err := utils.GenerateToken(user.ID, user.Name, false, model.OauthLoginType)
			if err != nil {
				log.Logger.Warningf("Failed to generate token: %s", err)
				ctx.JSON(http.StatusOK, gin.H{"msg": i18n.UnknownError, "data": nil})
				return
			}
			ctx.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("%s?token=%s", config.Env.OauthCallback, token))
		})

		availableOauth[provider.Provider] = gin.H{
			"url":    fmt.Sprintf("%s/oauth/%s", config.Env.Backend, provider.URI),
			"name":   provider.Provider,
			"avatar": provider.Avatar,
		}
	}

	router.GET("/oauth", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"msg": i18n.Success, "data": availableOauth})
	})
}
