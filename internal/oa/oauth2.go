package oa

import (
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
	"CBCTF/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type oauth2Protocol struct{}

func NewOAuth2Protocol() Protocol {
	return &oauth2Protocol{}
}

func (o *oauth2Protocol) ID() string {
	return model.OauthProtocolOAuth2
}

func (o *oauth2Protocol) config(provider model.Oauth) *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     provider.ClientID,
		ClientSecret: provider.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  provider.AuthURL,
			TokenURL: provider.TokenURL,
		},
		RedirectURL: provider.CallbackURL,
	}
	if len(provider.Scopes) > 0 {
		config.Scopes = append([]string(nil), provider.Scopes...)
	}
	return config
}

func (o *oauth2Protocol) LoginURL(provider model.Oauth) (string, model.RetVal) {
	state := utils.UUID()
	verifier := oauth2.GenerateVerifier()
	if ret := redis.SetOauthState(provider.Provider, state, verifier); !ret.OK {
		return "", ret
	}
	url := o.config(provider).AuthCodeURL(state, oauth2.AccessTypeOnline, oauth2.S256ChallengeOption(verifier))
	return url, model.SuccessRetVal()
}

func (o *oauth2Protocol) Exchange(ctx *gin.Context, provider model.Oauth) (map[string]any, model.RetVal) {
	var form dto.OauthCallbackForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		return nil, ret
	}
	defer redis.DelOauthState(provider.Provider, form.State)
	verifier, ret := redis.GetOauthVerifier(provider.Provider, form.State)
	if !ret.OK {
		return nil, ret
	}
	oauthConfig := o.config(provider)
	tok, err := oauthConfig.Exchange(ctx, form.Code, oauth2.VerifierOption(verifier))
	if err != nil {
		log.Logger.Warningf("Failed to get token for provider %s: %s", provider.Provider, err)
		return nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	client := oauthConfig.Client(ctx, tok)
	client.Timeout = time.Second * 10
	response, err := client.Get(provider.UserInfoURL)
	if err != nil {
		log.Logger.Warningf("Failed to get User info by provider %s: %s", provider.Provider, err)
		return nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Logger.Warningf("Failed to close response body for provider %s: %s", provider.Provider, err)
		}
	}(response.Body)
	var result map[string]any
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		err = fmt.Errorf("unexpected status %d: %s", response.StatusCode, strings.TrimSpace(string(body)))
		log.Logger.Warningf("Failed to get User info by provider %s: %s", provider.Provider, err)
		return nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Logger.Warningf("Failed to decode response body for provider %s: %s", provider.Provider, err)
		return nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	if err = ApplyUserInfoCallback(provider, client, result); err != nil {
		log.Logger.Warningf("Failed to apply oauth callback for provider %s: %s", provider.Provider, err)
	}
	return result, model.SuccessRetVal()
}
