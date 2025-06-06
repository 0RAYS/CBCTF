package model

import "CBCTF/internel/i18n"

const (
	StaticChallengeType  = "static"
	DynamicChallengeType = "dynamic"
	PodsChallengeType    = "pods"

	AttachmentFile = "attachment.zip"
	GeneratorFile  = "generator.zip"
)

// Challenge 题库中的挑战
// HasMany DockerGroup
// HasMany ChallengeFlag
// HasMany ContestChallenge
// HasMany Submission
type Challenge struct {
	DockerGroups      []DockerGroup      `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	ChallengeFlags    []ChallengeFlag    `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	ContestChallenges []ContestChallenge `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Submissions       []Submission       `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	RandID            string             `gorm:"uniqueIndex;not null" json:"rand_id"`
	Name              string             `json:"name"`
	Desc              string             `json:"desc"`
	Category          string             `json:"category"`
	Type              string             `json:"type"`
	GeneratorImage    string             `json:"generator_image"`
	Basic
}

func (c Challenge) GetModelName() string {
	return "Challenge"
}

func (c Challenge) GetID() uint {
	return c.ID
}

func (c Challenge) GetVersion() uint {
	return c.Version
}

func (c Challenge) CreateErrorString() string {
	return i18n.CreateChallengeError
}

func (c Challenge) DeleteErrorString() string {
	return i18n.DeleteChallengeError
}

func (c Challenge) GetErrorString() string {
	return i18n.GetChallengeError
}

func (c Challenge) NotFoundErrorString() string {
	return i18n.ChallengeNotFound
}

func (c Challenge) UpdateErrorString() string {
	return i18n.UpdateChallengeError
}

func (c Challenge) GetUniqueKey() []string {
	return []string{"id", "rand_id"}
}
