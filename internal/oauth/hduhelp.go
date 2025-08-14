package oauth

import (
	"CBCTF/internal/config"
	"CBCTF/internal/model"
	"embed"
	"fmt"
)

var (
	//go:embed avatar/hduhelp.png
	hduhelpFile      embed.FS
	HDUHelpAvatar, _ = hduhelpFile.ReadFile("avatar/hduhelp.png")
)

func GetDefaultHDUHelpOauth() model.Oauth {
	return model.Oauth{
		AuthURL:         "https://api.hduhelp.com/oauth/authorize",
		TokenURL:        "https://api.hduhelp.com/oauth/token",
		UserInfoURL:     "https://api.hduhelp.com/salmon_base/student/info",
		CallbackURL:     fmt.Sprintf("%s/oauth/hduhelp/callback", config.Env.Backend),
		ClientID:        "",
		ClientSecret:    "",
		Provider:        "HDUHelp",
		Uri:             "hduhelp",
		RespIDField:     "{data.staffId}",
		RespNameField:   "HDU_{data.staffId}",
		RespEmailField:  "{data.staffId}@hdu.edu.cn",
		RespAvatarField: "",
		RespDescField:   "{data.unitName} {data.majorName} {data.staffName}",
		On:              false,
		Avatar:          model.AvatarURL(fmt.Sprintf("%s/assets?filename=hduhelp", config.Env.Backend)),
	}
}
