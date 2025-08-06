package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateContestForm for create contest
type CreateContestForm struct {
	Name      string           `form:"name" json:"name" binding:"required"`
	Desc      string           `form:"desc" json:"desc"`
	Captcha   string           `form:"captcha" json:"captcha"`
	Prefix    string           `form:"prefix" json:"prefix"`
	Blood     bool             `form:"blood" json:"blood"`
	Size      int              `form:"size" json:"size"`
	Start     time.Time        `form:"start" json:"start"`
	Duration  int64            `form:"duration" json:"duration"`
	Rules     model.StringList `form:"rules" json:"rules"`
	Prizes    model.Prizes     `form:"prizes" json:"prizes"`
	Timelines model.Timelines  `form:"timelines" json:"timelines"`
	Hidden    bool             `form:"hidden" json:"hidden"`
}

func (f *CreateContestForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return false, i18n.BadRequest
	}
	f.Name = strings.TrimSpace(f.Name)
	if f.Name == "" {
		return false, i18n.BadRequest
	}
	f.Prefix = strings.TrimSpace(f.Prefix)
	return true, i18n.Success
}

// UpdateContestForm for admin update contest info
type UpdateContestForm struct {
	Name      *string           `form:"name" json:"name"`
	Desc      *string           `form:"desc" json:"desc"`
	Captcha   *string           `form:"captcha" json:"captcha"`
	Blood     *bool             `form:"blood" json:"blood"`
	Prefix    *string           `form:"prefix" json:"prefix"`
	Size      *int              `form:"start" json:"size"`
	Start     *time.Time        `form:"start" json:"start"`
	Duration  *int64            `form:"duration" json:"duration"`
	Rules     *model.StringList `form:"rules" json:"rules"`
	Prizes    *model.Prizes     `form:"prizes" json:"prizes"`
	Timelines *model.Timelines  `form:"timelines" json:"timelines"`
	Hidden    *bool             `form:"hidden" json:"hidden"`
}

func (f *UpdateContestForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return false, i18n.BadRequest
	}
	if f.Name != nil {
		f.Name = utils.Ptr(strings.TrimSpace(*f.Name))
		if *f.Name == "" {
			return false, i18n.BadRequest
		}
	}
	if f.Prefix != nil {
		f.Prefix = utils.Ptr(strings.TrimSpace(*f.Prefix))
	}
	return true, i18n.Success
}
