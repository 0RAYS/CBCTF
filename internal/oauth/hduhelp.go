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
		RespNameField:   "{staff_name} {staff_id}",
		RespEmailField:  "{staff_id}@hdu.edu.cn",
		RespAvatarField: "",
		RespDescField:   "{unit_name} {major_name}",
		On:              false,
		Avatar:          "",
	}
}
