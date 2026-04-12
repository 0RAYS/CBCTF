package oauth

import (
	"CBCTF/internal/config"
	"CBCTF/internal/model"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/oauth2/github"
)

// Download from: https://github.com/logos
var (
	//go:embed logo/github-mark-white.png
	githubMarkWhiteFile embed.FS
	GithubMarkWhite, _  = githubMarkWhiteFile.ReadFile("logo/github-mark-white.png")
	//go:embed logo/github-mark.png
	githubMarkFile embed.FS
	GithubMark, _  = githubMarkFile.ReadFile("logo/github-mark.png")
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
		Scopes:           model.StringList{"read:user", "user:email"},
		IDClaim:          "{id}",
		NameClaim:        "{login}",
		EmailClaim:       "{email}",
		PictureClaim:     "{picture_url}",
		DescriptionClaim: "{html_url}",
		On:               false,
		Picture:          model.FileURL(fmt.Sprintf("%s/assets?filename=github", config.Env.Host)),
	}
}

func IsGithubProvider(provider model.Oauth) bool {
	return strings.HasPrefix(strings.ToLower(provider.UserInfoURL), "https://api.github.com/")
}

func SetGithubEmail(_ model.Oauth, client *http.Client, data map[string]any) error {
	response, err := client.Get("https://api.github.com/user/emails")
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
