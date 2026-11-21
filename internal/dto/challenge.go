package dto

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
	"CBCTF/internal/utils"
)

type GetChallengesForm struct {
	ListModelsForm
	Type        model.ChallengeType `form:"type" json:"type" binding:"omitempty,oneof=static dynamic pods"`
	Category    string              `form:"category" json:"category"`
	Name        string              `form:"name" json:"name"`
	Description string              `form:"description" json:"description"`
	// ContestID 为空表示不过滤比赛，0 表示只看不归属比赛的题目。
	ContestID *uint `form:"contest_id" json:"contest_id"`
}

func (f *GetChallengesForm) Validate(ctx *gin.Context) model.RetVal {
	if ret := f.ListModelsForm.Validate(ctx); !ret.OK {
		return ret
	}
	f.Category = utils.ToTitle(f.Category)
	return model.SuccessRetVal()
}

type GetCategoriesForm struct {
	Type      model.ChallengeType `form:"type" json:"type" binding:"omitempty,oneof=static dynamic pods"`
	ContestID *uint               `form:"contest_id" json:"contest_id"`
}

// GetChallengesInContestForm 描述比赛内选手可见的题目列表。
type GetChallengesInContestForm struct {
	ListModelsForm
	Unsolved bool   `form:"unsolved" json:"unsolved"`
	Category string `form:"category" json:"category"`
}

func (f *GetChallengesInContestForm) Validate(ctx *gin.Context) model.RetVal {
	if ret := f.ListModelsForm.Validate(ctx); !ret.OK {
		return ret
	}
	f.Category = utils.ToTitle(f.Category)
	return model.SuccessRetVal()
}

// GetAllChallengesInContestForm 描述比赛内管理员可见的题目列表。
type GetAllChallengesInContestForm struct {
	ListModelsForm
	Type     model.ChallengeType `form:"type" json:"type" binding:"omitempty,oneof=static dynamic pods"`
	Name     string              `form:"name" json:"name"`
	Category string              `form:"category" json:"category"`
}

func (f *GetAllChallengesInContestForm) Validate(ctx *gin.Context) model.RetVal {
	if ret := f.ListModelsForm.Validate(ctx); !ret.OK {
		return ret
	}
	f.Category = utils.ToTitle(f.Category)
	return model.SuccessRetVal()
}

// CopyChallengesToContestForm 描述把已有题目复制到某场比赛。
type CopyChallengesToContestForm struct {
	ChallengeIDs []string `form:"challenge_ids" json:"challenge_ids" binding:"required,dive,uuid"`
}

type CreateChallengeForm struct {
	Name            string                `form:"name" json:"name" binding:"required"`
	Type            model.ChallengeType   `form:"type" json:"type" binding:"required,oneof=static dynamic pods"`
	Description     string                `form:"description" json:"description"`
	Category        string                `form:"category" json:"category"`
	ContestID       uint                  `form:"contest_id" json:"contest_id"`
	Hidden          *bool                 `form:"hidden" json:"hidden"`
	Attempt         *int64                `form:"attempt" json:"attempt" binding:"omitempty,gte=0"`
	Hints           []string              `form:"hints" json:"hints"`
	Tags            []string              `form:"tags" json:"tags"`
	Flags           []string              `form:"flags" json:"flags"`
	GeneratorImage  string                `form:"generator_image" json:"generator_image" binding:"required_if=Type dynamic"`
	DockerCompose   string                `form:"docker_compose" json:"docker_compose" binding:"required_if=Type pods"`
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
	Hidden          *bool                  `form:"hidden" json:"hidden"`
	Attempt         *int64                 `form:"attempt" json:"attempt" binding:"omitempty,gte=0"`
	Hints           *[]string              `form:"hints" json:"hints"`
	Tags            *[]string              `form:"tags" json:"tags"`
	GeneratorImage  *string                `form:"generator_image" json:"generator_image"`
	NetworkPolicies *model.NetworkPolicies `form:"network_policies" json:"network_policies"`
	DockerCompose   *string                `form:"docker_compose" json:"docker_compose"`
	Flags           []struct {
		ID    uint   `form:"id" json:"id"`
		Value string `form:"value" json:"value"`
	} `form:"flags" json:"flags"`
}

func (f *UpdateChallengeForm) Validate(_ *gin.Context) model.RetVal {
	if f.Category != nil {
		f.Category = new(utils.ToTitle(*f.Category))
	}
	return model.SuccessRetVal()
}

// CopyChallengeForm 描述把题目复制到目标比赛。
// ContestID 为 0 表示复制出不归属比赛的题目。
type CopyChallengeForm struct {
	ContestID uint `form:"contest_id" json:"contest_id"`
}
