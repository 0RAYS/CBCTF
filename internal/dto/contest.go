package dto

import (
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateContestForm for create contest
type CreateContestForm struct {
	Name        string           `form:"name" json:"name" binding:"required"`
	Description string           `form:"description" json:"description"`
	Captcha     string           `form:"captcha" json:"captcha"`
	Prefix      string           `form:"prefix" json:"prefix"`
	Blood       bool             `form:"blood" json:"blood"`
	Size        int              `form:"size" json:"size" binding:"omitempty,gte=1"`
	Start       time.Time        `form:"start" json:"start"`
	Duration    int64            `form:"duration" json:"duration" binding:"omitempty,gte=1"`
	Victims     int64            `form:"victims" json:"victims" binding:"omitempty,gte=1"`
	Rules       model.StringList `form:"rules" json:"rules"`
	Prizes      model.Prizes     `form:"prizes" json:"prizes"`
	Timelines   model.Timelines  `form:"timelines" json:"timelines"`
	Hidden      bool             `form:"hidden" json:"hidden"`
}

func (f *CreateContestForm) Validate(_ *gin.Context) model.RetVal {
	f.Prefix = strings.TrimSpace(f.Prefix)
	return model.SuccessRetVal()
}

// UpdateContestForm for admin update contest info
type UpdateContestForm struct {
	Name        *string           `form:"name" json:"name" binding:"omitempty,min=1"`
	Description *string           `form:"description" json:"description"`
	Captcha     *string           `form:"captcha" json:"captcha"`
	Blood       *bool             `form:"blood" json:"blood"`
	Prefix      *string           `form:"prefix" json:"prefix"`
	Size        *int              `form:"size" json:"size" binding:"omitempty,gte=1"`
	Start       *time.Time        `form:"start" json:"start"`
	Duration    *int64            `form:"duration" json:"duration" binding:"omitempty,gte=1"`
	Victims     *int64            `form:"victims" json:"victims" binding:"omitempty,gte=1"`
	Rules       *model.StringList `form:"rules" json:"rules"`
	Prizes      *model.Prizes     `form:"prizes" json:"prizes"`
	Timelines   *model.Timelines  `form:"timelines" json:"timelines"`
	Hidden      *bool             `form:"hidden" json:"hidden"`
}

func (f *UpdateContestForm) Validate(_ *gin.Context) model.RetVal {
	if f.Prefix != nil {
		f.Prefix = utils.Ptr(strings.TrimSpace(*f.Prefix))
	}
	return model.SuccessRetVal()
}
