package oa

import (
	"context"

	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
)

// Protocol covers only what differs between authentication protocols: building
// the login redirect and turning the callback request into a normalized claims
// map. The OAuth router owns the shared login/token/redirect logic.
type Protocol interface {
	// ID returns the protocol identifier persisted on model.Oauth.Protocol.
	ID() string
	// LoginURL builds the URL the end-user's browser is redirected to in
	// order to start authentication.
	LoginURL(ctx context.Context, provider model.Oauth) (string, model.RetVal)
	// Exchange consumes the authentication callback request and returns the
	// normalized claims map consumed by service.OauthLogin.
	Exchange(ctx *gin.Context, provider model.Oauth) (map[string]any, model.RetVal)
}

var protocols = make(map[string]Protocol)

func RegisterProtocol(p Protocol) {
	protocols[p.ID()] = p
}

func GetProtocol(id string) (Protocol, bool) {
	p, ok := protocols[id]
	return p, ok
}
