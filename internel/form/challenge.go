package form

import "CBCTF/internel/model"

// GetChallengesForm for get challenges list
type GetChallengesForm struct {
	Offset   int    `form:"offset" json:"offset"`
	Limit    int    `form:"limit" json:"limit"`
	Type     string `form:"type" json:"type"`
	Category string `form:"category" json:"category"`
}

// GetCategoriesForm for get categories list
type GetCategoriesForm struct {
	Type string `form:"type" json:"type"`
}

// DownloadChallengeForm for download challenge
type DownloadChallengeForm struct {
	File string `form:"file" json:"file" binding:"required"`
}

type CreateChallengeForm struct {
	Name           string           `form:"name" json:"name" binding:"required"`
	Type           string           `form:"type" json:"type" binding:"required"`
	Desc           string           `form:"desc" json:"desc"`
	Category       string           `form:"category" json:"category"`
	Flags          model.StringList `form:"flags" json:"flags"`
	GeneratorImage string           `form:"generator_image" json:"generator_image"`
	DockerGroups   []struct {
		Yaml            string                `form:"yaml" json:"yaml"`
		NetworkPolicies model.NetworkPolicies `form:"yaml" json:"network_policies"`
	} `form:"docker_groups" json:"docker_groups"`
}

type UpdateChallengeForm struct {
	Name           *string `form:"name" json:"name"`
	Desc           *string `form:"desc" json:"desc"`
	Category       *string `form:"category" json:"category"`
	GeneratorImage *string `form:"generator_image" json:"generator_image"`
	Flags          []struct {
		ID    uint   `form:"id" json:"id"`
		Value string `form:"value" json:"value"`
	} `form:"flags" json:"flags"`
	DockerGroups []struct {
		ID              uint                  `form:"id" json:"id"`
		NetworkPolicies model.NetworkPolicies `json:"network_policies"`
	} `form:"docker_groups" json:"docker_groups"`
}
