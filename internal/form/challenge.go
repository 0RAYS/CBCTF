package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"github.com/gin-gonic/gin"
	"slices"
	"strings"
)

var (
	allowedChallengeType = []string{model.StaticChallengeType, model.QuestionChallengeType, model.DynamicChallengeType, model.PodsChallengeType}
	allowedFileName      = []string{model.AttachmentFile, model.GeneratorFile}
)

// GetChallengesForm for get challenges list
type GetChallengesForm struct {
	Offset   int    `form:"offset" json:"offset"`
	Limit    int    `form:"limit" json:"limit"`
	Type     string `form:"type" json:"type"`
	Category string `form:"category" json:"category"`
}

func (f *GetChallengesForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	if f.Limit > 100 {
		f.Limit = 15
	}
	if _, exists := ctx.GetQuery("limit"); !exists {
		f.Limit = 10
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		f.Offset = 0
	}
	if !slices.Contains(allowedChallengeType, f.Type) {
		f.Type = ""
	}
	f.Category = utils.ToTitle(f.Category)
	return true, i18n.Success
}

// GetCategoriesForm for get categories list
type GetCategoriesForm struct {
	Type string `form:"type" json:"type"`
}

func (f *GetCategoriesForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	if !slices.Contains(allowedChallengeType, f.Type) {
		f.Type = ""
	}
	return true, i18n.Success
}

// DownloadChallengeForm for download challenge
type DownloadChallengeForm struct {
	File string `form:"file" json:"file" binding:"required"`
}

func (f *DownloadChallengeForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	if !slices.Contains(allowedFileName, f.File) {
		return false, i18n.InvalidFileName
	}
	return true, i18n.Success
}

type CreateChallengeForm struct {
	Name            string                `form:"name" json:"name" binding:"required"`
	Type            string                `form:"type" json:"type" binding:"required"`
	Desc            string                `form:"desc" json:"desc"`
	Category        string                `form:"category" json:"category"`
	Flags           model.StringList      `form:"flags" json:"flags"`
	GeneratorImage  string                `form:"generator_image" json:"generator_image"`
	DockerCompose   string                `form:"docker_compose" json:"docker_compose"`
	Options         model.Options         `form:"options" json:"options"`
	NetworkPolicies model.NetworkPolicies `form:"network_policies" json:"network_policies"`
}

func (f *CreateChallengeForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	f.Name = strings.TrimSpace(f.Name)
	if f.Name == "" {
		return false, i18n.BadRequest
	}
	if !slices.Contains(allowedChallengeType, f.Type) {
		return false, i18n.InvalidChallengeType
	}
	f.Category = utils.ToTitle(f.Category)
	return true, i18n.Success
}

type UpdateChallengeForm struct {
	Name            *string                `form:"name" json:"name"`
	Desc            *string                `form:"desc" json:"desc"`
	Category        *string                `form:"category" json:"category"`
	GeneratorImage  *string                `form:"generator_image" json:"generator_image"`
	NetworkPolicies *model.NetworkPolicies `json:"network_policies"`
	Options         *model.Options         `form:"options" json:"options"`
	Flags           []struct {
		ID    uint   `form:"id" json:"id"`
		Value string `form:"value" json:"value"`
	} `form:"flags" json:"flags"`
}

func (f *UpdateChallengeForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	if f.Name != nil {
		f.Name = utils.Ptr(strings.TrimSpace(*f.Name))
		if *f.Name == "" {
			return false, i18n.BadRequest
		}
	}
	if f.Category != nil {
		f.Category = utils.Ptr(utils.ToTitle(*f.Category))
	}
	return true, i18n.Success
}
