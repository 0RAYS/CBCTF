package oauth

import (
	"CBCTF/internal/config"
	"CBCTF/internal/model"
	"embed"
	"fmt"
)

var (
	//go:embed logo/hduhelp.png
	hduhelpFile       embed.FS
	HDUHelpPicture, _ = hduhelpFile.ReadFile("logo/hduhelp.png")
)

func GetDefaultHDUHelpOauth() model.Oauth {
	return model.Oauth{
		AuthURL:              "https://api.hduhelp.com/oauth/authorize",
		TokenURL:             "https://api.hduhelp.com/oauth/token",
		UserInfoURL:          "https://api.hduhelp.com/salmon_base/student/info",
		CallbackURL:          fmt.Sprintf("%s/oauth/hduhelp/callback", config.Env.Host),
		ClientID:             "",
		ClientSecret:         "",
		Provider:             "HDUHelp",
		Uri:                  "hduhelp",
		RespIDField:          "{data.staffId}",
		RespNameField:        "HDU_{data.staffId}",
		RespEmailField:       "{data.staffId}@hdu.edu.cn",
		RespPictureField:     "",
		RespDescriptionField: "{data.unitName} {data.majorName} {data.staffName}",
		On:                   false,
		Picture:              model.FileURL(fmt.Sprintf("%s/assets?filename=hduhelp", config.Env.Host)),
	}
}
