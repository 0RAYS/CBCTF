package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetBrandingResp(branding model.Branding) gin.H {
	return gin.H{
		"id":                  branding.ID,
		"code":                branding.Code,
		"site_name":           branding.SiteName,
		"admin_name":          branding.AdminName,
		"browser_title":       branding.BrowserTitle,
		"browser_description": branding.BrowserDescription,
		"footer_copyright":    branding.FooterCopyright,
		"home_logo":           branding.HomeLogo,
		"home_logo_alt":       branding.HomeLogoAlt,
		"home":                branding.Home,
	}
}
