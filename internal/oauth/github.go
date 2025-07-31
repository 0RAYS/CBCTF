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
	//go:embed avatar/github-mark-white.svg
	GithubMarkWhiteFile embed.FS
	GithubMarkWhite, _  = GithubMarkWhiteFile.ReadFile("avatar/github-mark-white.svg")
	//go:embed avatar/github-mark.svg
	GithubMarkFile embed.FS
	GithubMark, _  = GithubMarkFile.ReadFile("avatar/github-mark.svg")
)

func GetDefaultGithubOauth() model.Oauth {
	return model.Oauth{
		AuthURL:         github.Endpoint.AuthURL,
		TokenURL:        github.Endpoint.TokenURL,
		UserInfoURL:     "https://api.github.com/user",
		ClientID:        "",
		ClientSecret:    "",
		Provider:        "Github",
		RedirectURI:     "/oauth/callback/github",
		RespNameField:   "{name}",
		RespEmailField:  "{email}",
		RespAvatarField: "{avatar_url}",
		RespDescField:   "{blog}",
		Avatar:          model.AvatarURL(fmt.Sprintf("%s/assets?filename=github", config.Env.Backend)),
	}
}
