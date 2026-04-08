package router

import (
	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"

	"github.com/gin-gonic/gin"
)

func GetBranding(ctx *gin.Context) {
	branding, ret := db.InitBrandingRepo(db.DB).GetDefault()
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	resp.JSON(ctx, model.SuccessRetVal(resp.GetBrandingResp(branding)))
}

func GetAdminBranding(ctx *gin.Context) {
	GetBranding(ctx)
}

func UpdateBranding(ctx *gin.Context) {
	var form dto.UpdateBrandingForm
	if ret := dto.Bind(ctx, &form); !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventTypeKey, model.UpdateBrandingEventType)
	repo := db.InitBrandingRepo(db.DB)
	branding, ret := repo.GetDefault()
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ret = repo.Update(branding.ID, db.UpdateBrandingOptions{
		SiteName:           form.SiteName,
		AdminName:          form.AdminName,
		BrowserTitle:       form.BrowserTitle,
		BrowserDescription: form.BrowserDescription,
		FooterCopyright:    form.FooterCopyright,
		HomeLogo:           form.HomeLogo,
		HomeLogoAlt:        form.HomeLogoAlt,
		Home:               form.Home,
	})
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	branding, ret = repo.GetDefault()
	if !ret.OK {
		resp.JSON(ctx, ret)
		return
	}
	ctx.Set(middleware.CTXEventSuccessKey, true)
	resp.JSON(ctx, model.SuccessRetVal(resp.GetBrandingResp(branding)))
}
