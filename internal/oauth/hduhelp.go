package oauth

import "CBCTF/internal/model"

func GetDefaultHDUHelpOauth() model.Oauth {
	return model.Oauth{
		AuthURL:         "https://api.hduhelp.com/oauth/authorize",
		TokenURL:        "https://api.hduhelp.com/oauth/token",
		UserInfoURL:     "https://api.hduhelp.com/salmon_base/student/info",
		ClientID:        "",
		ClientSecret:    "",
		Provider:        "HDUHelp",
		RedirectURI:     "/oauth/callback/hduhelp",
		RespNameField:   "staff_name",
		RespEmailField:  "",
		RespAvatarField: "",
		RespDescField:   "",
		Avatar:          "",
	}
}
