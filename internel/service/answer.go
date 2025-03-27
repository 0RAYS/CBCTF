package service

import (
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"fmt"
	"gorm.io/gorm"
)

// IsGenerated model.Usage 需要递归预加载, depth = 3
func IsGenerated(tx *gorm.DB, usage model.Usage, team model.Team) bool {
	repo := db.InitAnswerRepo(tx)
	for _, flag := range usage.Flags {
		answers, _, ok, _ := repo.GetAll(flag.ID, -1, -1, false, 0)
		if !ok {
			return false
		}
		var count int
		for _, answer := range answers {
			if answer.TeamID == team.ID {
				count++
			}
		}
		switch usage.Challenge.Type {
		case model.StaticChallenge, model.DynamicChallenge, model.DockerChallenge:
			if count < 1 {
				return false
			}
		case model.DockersChallenge:
			if count < len(flag.Answers) {
				return false
			}
		default:
			return false
		}
	}
	return true
}

// InitAnswer model.Usage 需要预加载
func InitAnswer(tx *gorm.DB, usage model.Usage, team model.Team) ([]model.Answer, bool, string) {
	repo := db.InitAnswerRepo(tx)
	answers := make([]model.Answer, 0)
	if len(usage.Flags) < 1 {
		return answers, false, "FlagNotFound"
	}
	options := make([]db.CreateAnswerOptions, 0)
	switch usage.Challenge.Type {
	case model.StaticChallenge:
		option := db.CreateAnswerOptions{
			TeamID: team.ID,
			FlagID: usage.Flags[0].ID,
			Value:  usage.Flags[0].Value,
		}
		options = append(options, option)
	case model.DynamicChallenge:
		option := db.CreateAnswerOptions{
			TeamID: team.ID,
			FlagID: usage.Flags[0].ID,
		}
		if usage.Flags[0].Value == "uuid" {
			option.Value = utils.UUID()
		} else {
			option.Value = utils.RandFlag(usage.Flags[0].Value)
		}
		options = append(options, option)
	case model.DockerChallenge:
		option := db.CreateAnswerOptions{
			TeamID: team.ID,
			FlagID: usage.Flags[0].ID,
			Value:  utils.UUID(),
		}
		options = append(options, option)
	case model.DockersChallenge:
		for _, flag := range usage.Flags {
			option := db.CreateAnswerOptions{
				TeamID: team.ID,
				FlagID: flag.ID,
				Value:  utils.UUID(),
			}
			options = append(options, option)
		}
	default:
		return answers, false, "InvalidChallengeType"
	}
	for _, option := range options {
		if _, ok, _ := repo.GetBy2ID(team.ID, option.FlagID, false, 0); ok {
			continue
		}
		option.Value = fmt.Sprintf("%s{%s}", usage.Contest.Prefix, option.Value)
		answer, ok, msg := repo.Create(option)
		if !ok {
			return answers, false, msg
		}
		answers = append(answers, answer)
	}
	return answers, true, "Success"
}
