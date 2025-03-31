package service

import (
	f "CBCTF/internel/form"
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"gorm.io/gorm"
)

func CreateChallenge(tx *gorm.DB, form f.CreateChallengeForm) (model.Challenge, bool, string) {
	repo := db.InitChallengeRepo(tx)
	return repo.Create(db.CreateChallengeOptions{
		ID:        utils.UUID(),
		Name:      form.Name,
		Desc:      form.Desc,
		Category:  utils.ToTitle(form.Category),
		Type:      form.Type,
		Generator: form.Generator,
		Flags:     form.Flags,
		Docker:    form.Docker,
		Dockers:   form.Dockers,
	})
}

func UpdateChallenge(tx *gorm.DB, challenge model.Challenge, form f.UpdateChallengeForm) (bool, string) {
	repo := db.InitChallengeRepo(tx)
	options := db.UpdateChallengeOptions{
		Name: form.Name,
		Desc: form.Desc,
		Category: func() *string {
			tmp := utils.ToTitle(*form.Category)
			return &tmp
		}(),
		Type: form.Type,
	}
	switch *form.Type {
	case model.StaticChallenge:
		options.Flags = form.Flags
	case model.DynamicChallenge:
		options.Flags = form.Flags
		options.Generator = form.Generator
	case model.DockerChallenge:
		options.Docker = form.Docker
	case model.DockersChallenge:
		options.Dockers = form.Dockers
	default:
		return false, "InvalidChallengeType"
	}
	return repo.Update(challenge.ID, options)
}

func DeleteChallenge(tx *gorm.DB, challenge model.Challenge) (bool, string) {
	repo := db.InitChallengeRepo(tx)
	return repo.Delete(challenge.ID)
}
