package oa

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/oauth2/github"

	"CBCTF/internal/config"
	"CBCTF/internal/model"
)

// Download from: https://github.com/logos
var (
	//go:embed logo/github-mark.png
	githubMarkFile embed.FS
	GithubLogo, _  = githubMarkFile.ReadFile("logo/github-mark.png")
)

func GetDefaultGithubOauth() model.Oauth {
	return model.Oauth{
		AuthURL:          github.Endpoint.AuthURL,
		TokenURL:         github.Endpoint.TokenURL,
		UserInfoURL:      "https://api.github.com/user",
		CallbackURL:      fmt.Sprintf("%s/oauth/github/callback", config.Env.Host),
		ClientID:         "",
		ClientSecret:     "",
		Provider:         "Github",
		Uri:              "github",
		Scopes:           []string{"user:email"},
		IDClaim:          "{id}",
		NameClaim:        "{login}",
		EmailClaim:       "{email}",
		PictureClaim:     "{avatar_url}",
		DescriptionClaim: "{html_url}",
		On:               false,
		Picture:          model.FileURL(fmt.Sprintf("%s/assets?filename=github", config.Env.Host)),
	}
}

func IsGithubProvider(provider model.Oauth) bool {
	return strings.HasPrefix(strings.ToLower(provider.UserInfoURL), "https://api.github.com/")
}

func SetGithubEmail(ctx context.Context, _ model.Oauth, client *http.Client, data map[string]any) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user/emails", nil)
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 4096))
		return fmt.Errorf("unexpected status %d: %s", response.StatusCode, strings.TrimSpace(string(body)))
	}
	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err = json.NewDecoder(response.Body).Decode(&emails); err != nil {
		return err
	}
	for _, email := range emails {
		if email.Primary && email.Verified && email.Email != "" {
			data["email"] = email.Email
			return nil
		}
	}
	return nil
}
