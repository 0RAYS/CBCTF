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
	if form.Type == model.PodsChallenge {
		for i, docker := range form.Dockers {
			if len(docker.NetworkPolicies) == 0 {
				form.Dockers[i].NetworkPolicies = append(form.Dockers[i].NetworkPolicies, model.DefaultNetworkPolicy)
			}
		}
	}
	return repo.Create(db.CreateChallengeOptions{
		ID:        utils.UUID(),
		Name:      form.Name,
		Desc:      form.Desc,
		Category:  utils.ToTitle(form.Category),
		Type:      form.Type,
		Generator: form.Generator,
		Flags:     form.Flags,
		Dockers:   form.Dockers,
	})
}

func UpdateChallenge(tx *gorm.DB, challenge model.Challenge, form f.UpdateChallengeForm) (bool, string) {
	repo := db.InitChallengeRepo(tx)
	options := db.UpdateChallengeOptions{
		Name: form.Name,
		Desc: form.Desc,
		Category: func() *string {
			if form.Category != nil {
				tmp := utils.ToTitle(*form.Category)
				return &tmp
			}
			return nil
		}(),
		Type: form.Type,
	}
	targetType := challenge.Type
	if form.Type != nil {
		targetType = *form.Type
	}
	switch targetType {
	case model.StaticChallenge:
		options.Flags = form.Flags
	case model.DynamicChallenge:
		options.Flags = form.Flags
		options.Generator = form.Generator
	case model.PodsChallenge:
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
