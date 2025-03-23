package form

// DownloadChallengeForm for download challenge
type DownloadChallengeForm struct {
	File string `form:"file" json:"file" binding:"required"`
}

// GetCategoriesForm for get categories list
type GetCategoriesForm struct {
	Type string `form:"type" json:"type"`
}

type DeleteChallengeForm struct {
	Force bool `form:"force" json:"force"`
}

// UpdateChallengeForm for admin update challenge info
type UpdateChallengeForm struct {
	Name           *string `form:"name" json:"name"`
	Desc           *string `form:"desc" json:"desc"`
	Flag           *string `form:"flag" json:"flag"`
	Category       *string `form:"category" json:"category"`
	Type           *string `form:"type" json:"type"`
	GeneratorImage *string `form:"generator" json:"generator"`
	DockerImage    *string `form:"docker" json:"docker"`
	Port           *int32  `form:"port" json:"port"`
}

// GetChallengesForm for get challenges list
type GetChallengesForm struct {
	Offset   int    `form:"offset" json:"offset"`
	Limit    int    `form:"limit" json:"limit"`
	Type     string `form:"type" json:"type"`
	Category string `form:"category" json:"category"`
}

// CreateChallengeForm for create challenge
type CreateChallengeForm struct {
	Name           string `form:"name" json:"name" binding:"required"`
	Desc           string `form:"desc" json:"desc"`
	Flag           string `form:"flag" json:"flag"`
	Category       string `form:"category" json:"category"`
	Type           string `form:"type" json:"type"`
	GeneratorImage string `form:"generator" json:"generator"`
	DockerImage    string `form:"docker" json:"docker"`
	Port           int32  `form:"port" json:"port"`
}
