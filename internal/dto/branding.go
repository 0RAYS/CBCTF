package dto

import "CBCTF/internal/model"

type UpdateBrandingForm struct {
	SiteName           *model.LocalizedText       `form:"site_name" json:"site_name"`
	AdminName          *model.LocalizedText       `form:"admin_name" json:"admin_name"`
	BrowserTitle       *model.LocalizedText       `form:"browser_title" json:"browser_title"`
	BrowserDescription *model.LocalizedText       `form:"browser_description" json:"browser_description"`
	FooterCopyright    *model.LocalizedText       `form:"footer_copyright" json:"footer_copyright"`
	FooterICPNumber    *string                    `form:"footer_icp_number" json:"footer_icp_number"`
	FooterICPLink      *string                    `form:"footer_icp_link" json:"footer_icp_link"`
	FooterContactEmail *string                    `form:"footer_contact_email" json:"footer_contact_email"`
	FooterGithubURL    *string                    `form:"footer_github_url" json:"footer_github_url"`
	HomeLogo           *model.FileURL             `form:"home_logo" json:"home_logo"`
	HomeLogoAlt        *model.LocalizedText       `form:"home_logo_alt" json:"home_logo_alt"`
	Home               *model.BrandingHomeContent `form:"home" json:"home"`
}
