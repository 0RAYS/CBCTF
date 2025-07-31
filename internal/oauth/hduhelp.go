package oauth

import (
	"CBCTF/internal/config"
	"CBCTF/internal/model"
	"embed"
	"fmt"
)

var (
	//go:embed avatar/hduhelp.png
	HDUHelpFile      embed.FS
	HDUHelpAvatar, _ = HDUHelpFile.ReadFile("avatar/hduhelp.png")
)

func GetDefaultHDUHelpOauth() model.Oauth {
	return model.Oauth{
		AuthURL:         "https://api.hduhelp.com/oauth/authorize",
		TokenURL:        "https://api.hduhelp.com/oauth/token",
		UserInfoURL:     "https://api.hduhelp.com/salmon_base/student/info",
		ClientID:        "",
		ClientSecret:    "",
		Provider:        "HDUHelp",
		RedirectURI:     "/oauth/callback/hduhelp",
		RespNameField:   "{data.staffName} {data.staffId}",
		RespEmailField:  "{data.staffId}@hdu.edu.cn",
		RespAvatarField: "",
		RespDescField:   "{data.unitName} {data.majorName}",
		On:              false,
		Avatar:          model.AvatarURL(fmt.Sprintf("%s/assets?filename=hduhelp", config.Env.Backend)),
	}
}
