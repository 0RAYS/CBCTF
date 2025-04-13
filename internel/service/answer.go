package service

import (
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"CBCTF/internel/utils"
	"fmt"
	"gorm.io/gorm"
)

// IsGenerated model.Usage 需要递归预加载, Usage.Flags
func IsGenerated(tx *gorm.DB, team model.Team, usage model.Usage) bool {
	repo := db.InitAnswerRepo(tx)
	for _, flag := range usage.Flags {
		_, ok, _ := repo.GetBy2ID(team.ID, flag.ID)
		if !ok {
			return false
		}
	}
	return true
}

// GenerateAnswer model.Usage 需要预加载
func GenerateAnswer(tx *gorm.DB, usage model.Usage, team model.Team, reset bool) ([]model.Answer, bool, string) {
	repo := db.InitAnswerRepo(tx)
	answers := make([]model.Answer, 0)
	if len(usage.Flags) < 1 {
		return answers, false, "FlagNotFound"
	}
	options := make([]db.CreateAnswerOptions, 0)
	for _, flag := range usage.Flags {
		option := db.CreateAnswerOptions{
			TeamID: team.ID,
			FlagID: flag.ID,
			Solved: false,
		}
		if result := model.StaticFlag.FindAllStringSubmatch(flag.Value, 1); len(result) > 0 {
			option.Value = flag.Value
		} else if result := model.DynamicFlag.FindAllStringSubmatch(flag.Value, 1); len(result) > 0 {
			option.Value = utils.RandFlag(result[0][1])
		} else if result := model.UUIDFlag.FindAllStringSubmatch(flag.Value, 1); len(result) > 0 {
			option.Value = utils.UUID()
		} else {
			option.Value = flag.Value
		}
		options = append(options, option)
	}
	for _, option := range options {
		if answer, ok, _ := repo.GetBy2ID(team.ID, option.FlagID); !reset && ok {
			answers = append(answers, answer)
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
