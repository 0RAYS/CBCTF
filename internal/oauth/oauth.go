package oauth

import (
	"CBCTF/internal/model"
	"net/http"
)

type ProviderCallback func(provider model.Oauth, client *http.Client, response map[string]any) error

type ProviderMatcher func(provider model.Oauth) bool

type providerHandler struct {
	match    ProviderMatcher
	scopes   []string
	callback ProviderCallback
}

var providerHandlers []providerHandler

func Init() {
	providerHandlers = make([]providerHandler, 0)
	providerHandlers = append(providerHandlers, providerHandler{
		match:    IsGithubProvider,
		scopes:   []string{"read:user", "user:email"},
		callback: SetGithubEmail,
	})
}

func GetScopes(provider model.Oauth) []string {
	for _, handler := range providerHandlers {
		if handler.match != nil && handler.match(provider) {
			if len(handler.scopes) == 0 {
				return nil
			}
			return append([]string(nil), handler.scopes...)
		}
	}
	return nil
}

func ApplyUserInfoCallback(provider model.Oauth, client *http.Client, response map[string]any) error {
	for _, handler := range providerHandlers {
		if handler.match != nil && handler.match(provider) {
			if handler.callback == nil {
				return nil
			}
			return handler.callback(provider, client, response)
		}
	}
	return nil
}
