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
	//go:embed avatar/github-mark-white.png
	githubMarkWhiteFile embed.FS
	GithubMarkWhite, _  = githubMarkWhiteFile.ReadFile("avatar/github-mark-white.png")
	//go:embed avatar/github-mark.png
	githubMarkFile embed.FS
	GithubMark, _  = githubMarkFile.ReadFile("avatar/github-mark.png")
)

func GetDefaultGithubOauth() model.Oauth {
	return model.Oauth{
		AuthURL:         github.Endpoint.AuthURL,
		TokenURL:        github.Endpoint.TokenURL,
		UserInfoURL:     "https://api.github.com/user",
		CallbackURL:     fmt.Sprintf("%s/oauth/github/callback", config.Env.Backend),
		ClientID:        "",
		ClientSecret:    "",
		Provider:        "Github",
		URI:             "github",
		RespIDField:     "{id}",
		RespNameField:   "{name}",
		RespEmailField:  "{email}",
		RespAvatarField: "{avatar_url}",
		RespDescField:   "{blog}",
		On:              false,
		Avatar:          model.AvatarURL(fmt.Sprintf("%s/assets?filename=github", config.Env.Backend)),
	}
}
