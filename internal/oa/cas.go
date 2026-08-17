package oa

import (
	"CBCTF/internal/dto"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type casProtocol struct{}

func NewCASProtocol() Protocol {
	return &casProtocol{}
}

func (c *casProtocol) ID() string {
	return model.OauthProtocolCAS
}

func (c *casProtocol) LoginURL(provider model.Oauth) (string, model.RetVal) {
	loginURL, err := buildCASLoginURL(provider)
	if err != nil {
		return "", model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	return loginURL, model.SuccessRetVal()
}

func (c *casProtocol) Exchange(ctx *gin.Context, provider model.Oauth) (map[string]any, model.RetVal) {
	var form dto.CASCallbackForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		return nil, ret
	}
	validateURL, err := buildCASValidateURL(provider, form.Ticket)
	if err != nil {
		return nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	log.Logger.Debugf("CAS callback for provider %s: %s", provider.Provider, validateURL)
	client := http.Client{Timeout: time.Second * 10}
	response, err := client.Get(validateURL)
	if err != nil {
		log.Logger.Warningf("Failed to validate CAS ticket for provider %s: %s", provider.Provider, err)
		return nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Logger.Warningf("Failed to close CAS response body for provider %s: %s", provider.Provider, err)
		}
	}(response.Body)
	result, err := parseCASServiceResponse(response.Body)
	if err != nil {
		log.Logger.Warningf("Failed to parse CAS response for provider %s: %s", provider.Provider, err)
		return nil, model.RetVal{Msg: i18n.Common.UnknownError, Attr: map[string]any{"Error": err.Error()}}
	}
	return result, model.SuccessRetVal()
}

func buildCASLoginURL(provider model.Oauth) (string, error) {
	loginURL, err := url.Parse(provider.AuthURL)
	if err != nil {
		return "", err
	}
	query := loginURL.Query()
	query.Set("service", provider.CallbackURL)
	loginURL.RawQuery = query.Encode()
	return loginURL.String(), nil
}

func buildCASValidateURL(provider model.Oauth, ticket string) (string, error) {
	validateURL, err := url.Parse(provider.UserInfoURL)
	if err != nil {
		return "", err
	}
	query := validateURL.Query()
	query.Set("ticket", ticket)
	query.Set("service", provider.CallbackURL)
	validateURL.RawQuery = query.Encode()
	return validateURL.String(), nil
}

// parseCASServiceResponse parses a CAS 3.0 serviceValidate response into the
// same claims shape consumed by service.OauthLogin: a top-level "user" key plus
// an "attributes" map of additional claims.
func parseCASServiceResponse(body io.Reader) (map[string]any, error) {
	decoder := xml.NewDecoder(body)
	result := map[string]any{}
	attributes := map[string]any{}
	var stack []string
	var failure string
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			stack = append(stack, t.Name.Local)
		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		case xml.CharData:
			value := strings.TrimSpace(string(t))
			if value == "" || len(stack) == 0 {
				continue
			}
			current := stack[len(stack)-1]
			if current == "user" {
				result["user"] = value
				continue
			}
			if slices.Contains(stack, "attributes") && current != "attributes" {
				attributes[current] = value
				continue
			}
			if current == "authenticationFailure" {
				failure = value
			}
		}
	}
	if failure != "" {
		return nil, fmt.Errorf("CAS authentication failed: %s", failure)
	}
	if _, ok := result["user"]; !ok {
		return nil, errors.New("CAS response missing user")
	}
	result["attributes"] = attributes
	return result, nil
}
