package dto

import (
	"CBCTF/internal/model"
	"CBCTF/internal/utils"

	"github.com/gin-gonic/gin"
)

// GetChallengesForm for get challenges list
type GetChallengesForm struct {
	ListModelsForm
	Type     model.ChallengeType `form:"type" json:"type" binding:"omitempty,oneof=static question dynamic pods"`
	Category string              `form:"category" json:"category"`
}

func (f *GetChallengesForm) Validate(ctx *gin.Context) model.RetVal {
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 10
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
	f.Category = utils.ToTitle(f.Category)
	return model.SuccessRetVal()
}

// GetCategoriesForm for get categories list
type GetCategoriesForm struct {
	Type model.ChallengeType `form:"type" json:"type" binding:"omitempty,oneof=static question dynamic pods"`
}

type CreateChallengeForm struct {
	Name            string                `form:"name" json:"name" binding:"required"`
	Type            model.ChallengeType   `form:"type" json:"type" binding:"required,oneof=static question dynamic pods"`
	Description     string                `form:"description" json:"description"`
	Category        string                `form:"category" json:"category"`
	Flags           model.StringList      `form:"flags" json:"flags"`
	GeneratorImage  string                `form:"generator_image" json:"generator_image" binding:"required_if=Type dynamic"`
	DockerCompose   string                `form:"docker_compose" json:"docker_compose" binding:"required_if=Type pods"`
	Options         model.Options         `form:"options" json:"options"  binding:"required_if=Type question"`
	NetworkPolicies model.NetworkPolicies `form:"network_policies" json:"network_policies" `
}

func (f *CreateChallengeForm) Validate(_ *gin.Context) model.RetVal {
	f.Category = utils.ToTitle(f.Category)
	return model.SuccessRetVal()
}

type UpdateChallengeForm struct {
	Name            *string                `form:"name" json:"name" binding:"omitempty,min=1"`
	Description     *string                `form:"description" json:"description"`
	Category        *string                `form:"category" json:"category"`
	GeneratorImage  *string                `form:"generator_image" json:"generator_image"`
	NetworkPolicies *model.NetworkPolicies `form:"network_policies" json:"network_policies"`
	Options         *model.Options         `form:"options" json:"options"`
	DockerCompose   *string                `form:"docker_compose" json:"docker_compose"`
	Flags           []struct {
		ID    uint   `form:"id" json:"id"`
		Value string `form:"value" json:"value"`
	} `form:"flags" json:"flags"`
}

func (f *UpdateChallengeForm) Validate(_ *gin.Context) model.RetVal {
	if f.Category != nil {
		f.Category = utils.Ptr(utils.ToTitle(*f.Category))
	}
	return model.SuccessRetVal()
}
