package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"slices"

	"github.com/gin-gonic/gin"
)

var (
	allowedChallengeType = []string{model.StaticChallengeType, model.QuestionChallengeType, model.DynamicChallengeType, model.PodsChallengeType}
)

// GetChallengesForm for get challenges list
type GetChallengesForm struct {
	Offset   int    `form:"offset" json:"offset"`
	Limit    int    `form:"limit" json:"limit"`
	Type     string `form:"type" json:"type"`
	Category string `form:"category" json:"category"`
}

func (f *GetChallengesForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	if f.Limit > 100 || f.Limit < 0 {
		f.Limit = 15
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 10
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
	if !slices.Contains(allowedChallengeType, f.Type) {
		f.Type = ""
	}
	f.Category = utils.ToTitle(f.Category)
	return model.SuccessRetVal()
}

// GetCategoriesForm for get categories list
type GetCategoriesForm struct {
	Type string `form:"type" json:"type"`
}

func (f *GetCategoriesForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	if !slices.Contains(allowedChallengeType, f.Type) {
		f.Type = ""
	}
	return model.SuccessRetVal()
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

func (f *CreateChallengeForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	if !slices.Contains(allowedChallengeType, f.Type) {
		return model.RetVal{Msg: i18n.Model.Challenge.InvalidType}
	}
	f.Category = utils.ToTitle(f.Category)
	return model.SuccessRetVal()
}

type UpdateChallengeForm struct {
	Name            *string                `form:"name" json:"name"`
	Desc            *string                `form:"desc" json:"desc"`
	Category        *string                `form:"category" json:"category"`
	GeneratorImage  *string                `form:"generator_image" json:"generator_image"`
	NetworkPolicies *model.NetworkPolicies `json:"network_policies"`
	Options         *model.Options         `form:"options" json:"options"`
	DockerCompose   *string                `form:"docker_compose" json:"docker_compose"`
	Flags           []struct {
		ID    uint   `form:"id" json:"id"`
		Value string `form:"value" json:"value"`
	} `form:"flags" json:"flags"`
}

func (f *UpdateChallengeForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	if f.Category != nil {
		f.Category = utils.Ptr(utils.ToTitle(*f.Category))
	}
	return model.SuccessRetVal()
}
