package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type BrandingRepo struct {
	BaseRepo[model.Branding]
}

type CreateBrandingOptions struct {
	Code               string
	SiteName           model.LocalizedText
	AdminName          model.LocalizedText
	BrowserTitle       model.LocalizedText
	BrowserDescription model.LocalizedText
	FooterCopyright    model.LocalizedText
	HomeLogo           model.FileURL
	HomeLogoAlt        model.LocalizedText
	Home               model.BrandingHomeContent
}

func (c CreateBrandingOptions) Convert2Model() model.Model {
	return model.Branding{
		Code:               c.Code,
		SiteName:           c.SiteName,
		AdminName:          c.AdminName,
		BrowserTitle:       c.BrowserTitle,
		BrowserDescription: c.BrowserDescription,
		FooterCopyright:    c.FooterCopyright,
		HomeLogo:           c.HomeLogo,
		HomeLogoAlt:        c.HomeLogoAlt,
		Home:               c.Home,
	}
}

type UpdateBrandingOptions struct {
	SiteName           *model.LocalizedText
	AdminName          *model.LocalizedText
	BrowserTitle       *model.LocalizedText
	BrowserDescription *model.LocalizedText
	FooterCopyright    *model.LocalizedText
	HomeLogo           *model.FileURL
	HomeLogoAlt        *model.LocalizedText
	Home               *model.BrandingHomeContent
}

func (u UpdateBrandingOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.SiteName != nil {
		options["site_name"] = *u.SiteName
	}
	if u.AdminName != nil {
		options["admin_name"] = *u.AdminName
	}
	if u.BrowserTitle != nil {
		options["browser_title"] = *u.BrowserTitle
	}
	if u.BrowserDescription != nil {
		options["browser_description"] = *u.BrowserDescription
	}
	if u.FooterCopyright != nil {
		options["footer_copyright"] = *u.FooterCopyright
	}
	if u.HomeLogo != nil {
		options["home_logo"] = *u.HomeLogo
	}
	if u.HomeLogoAlt != nil {
		options["home_logo_alt"] = *u.HomeLogoAlt
	}
	if u.Home != nil {
		options["home"] = *u.Home
	}
	return options
}

func InitBrandingRepo(tx *gorm.DB) *BrandingRepo {
	return &BrandingRepo{
		BaseRepo: BaseRepo[model.Branding]{
			DB: tx,
		},
	}
}

func (b *BrandingRepo) GetDefault(optionsL ...GetOptions) (model.Branding, model.RetVal) {
	return b.GetByUniqueField("code", model.DefaultBrandingCode, optionsL...)
}

func (b *BrandingRepo) InitDefault() model.RetVal {
	if _, ret := b.Create(CreateBrandingOptions{
		Code:               model.DefaultBranding().Code,
		SiteName:           model.DefaultBranding().SiteName,
		AdminName:          model.DefaultBranding().AdminName,
		BrowserTitle:       model.DefaultBranding().BrowserTitle,
		BrowserDescription: model.DefaultBranding().BrowserDescription,
		FooterCopyright:    model.DefaultBranding().FooterCopyright,
		HomeLogo:           model.DefaultBranding().HomeLogo,
		HomeLogoAlt:        model.DefaultBranding().HomeLogoAlt,
		Home:               model.DefaultBranding().Home,
	}); !ret.OK && ret.Msg != i18n.Model.DuplicateKeyValue {
		return ret
	}
	return model.SuccessRetVal()
}
