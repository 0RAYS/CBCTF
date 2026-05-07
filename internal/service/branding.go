package service

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

func GetDefaultBranding(tx *gorm.DB) (model.Branding, model.RetVal) {
	return db.InitBrandingRepo(tx).GetDefault()
}

func GetDefaultBrandingID(tx *gorm.DB) (uint, model.RetVal) {
	branding, ret := GetDefaultBranding(tx)
	if !ret.OK {
		return 0, ret
	}
	return branding.ID, model.SuccessRetVal()
}

func UpdateBranding(tx *gorm.DB, form dto.UpdateBrandingForm) (model.Branding, model.RetVal) {
	repo := db.InitBrandingRepo(tx)
	branding, ret := repo.GetDefault()
	if !ret.OK {
		return model.Branding{}, ret
	}
	ret = repo.Update(branding.ID, db.UpdateBrandingOptions{
		SiteName:           form.SiteName,
		AdminName:          form.AdminName,
		BrowserTitle:       form.BrowserTitle,
		BrowserDescription: form.BrowserDescription,
		FooterCopyright:    form.FooterCopyright,
		FooterICPNumber:    form.FooterICPNumber,
		FooterICPLink:      form.FooterICPLink,
		FooterContactEmail: form.FooterContactEmail,
		FooterGithubURL:    form.FooterGithubURL,
		HomeLogo:           form.HomeLogo,
		HomeLogoAlt:        form.HomeLogoAlt,
		Home:               form.Home,
	})
	if !ret.OK {
		return model.Branding{}, ret
	}
	return repo.GetDefault()
}
