package oa

import (
	"CBCTF/internal/config"
	"CBCTF/internal/model"
	"embed"
	"fmt"
)

var (
	//go:embed logo/hdu.png
	hduLogoFile embed.FS
	HDULogo, _  = hduLogoFile.ReadFile("logo/hdu.png")
)

func GetDefaultHDUCASOauth() model.Oauth {
	return model.Oauth{
		Protocol:         model.OauthProtocolCAS,
		AuthURL:          "https://sso.hdu.edu.cn/login",
		UserInfoURL:      "https://sso.hdu.edu.cn/p3/serviceValidate",
		CallbackURL:      fmt.Sprintf("%s/oauth/hducas/callback", config.Env.Host),
		Provider:         "HDU CAS",
		Uri:              "hducas",
		IDClaim:          "{user}",
		NameClaim:        "{attributes.XM}",
		EmailClaim:       "{user}@hdu.edu.cn",
		DescriptionClaim: "杭州电子科技大学 {attributes.DWMC} {attributes.SXZYMC} {user} {attributes.XM}",
		On:               false,
		Picture:          model.FileURL(fmt.Sprintf("%s/assets?filename=hdu", config.Env.Host)),
	}
}
