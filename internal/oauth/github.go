package oauth

import (
	"CBCTF/internal/config"
	"CBCTF/internal/model"
	"embed"
	"fmt"

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
		IDClaim:          "{id}",
		NameClaim:        "{login}",
		EmailClaim:       "{email}",
		PictureClaim:     "{picture_url}",
		DescriptionClaim: "{html_url}",
		On:               false,
		Picture:          model.FileURL(fmt.Sprintf("%s/assets?filename=github", config.Env.Host)),
	}
}
