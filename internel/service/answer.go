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
func GenerateAnswer(tx *gorm.DB, team model.Team, usage model.Usage) ([]model.Answer, bool, string) {
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
			option.Value = result[0][1]
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
		if answer, ok, _ := repo.GetBy2ID(team.ID, option.FlagID); ok {
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

// ResetAnswer model.Usage 需要预加载
func ResetAnswer(tx *gorm.DB, team model.Team, usage model.Usage) ([]model.Answer, bool, string) {
	answers := make([]model.Answer, 0)
	submissionRepo, answerRepo := db.InitSubmissionRepo(tx), db.InitAnswerRepo(tx)
	submissionIDL, answerIDL := make([]uint, 0), make([]uint, 0)
	submissions, _, ok, msg := submissionRepo.GetByKeyID("team_id", team.ID, -1, -1, false)
	if !ok {
		return answers, false, msg
	}
	for _, submission := range submissions {
		submissionIDL = append(submissionIDL, submission.ID)
	}
	for _, flag := range usage.Flags {
		answer, ok, msg := answerRepo.GetBy2ID(team.ID, flag.ID)
		if !ok {
			return answers, false, msg
		}
		answerIDL = append(answerIDL, answer.ID)
	}
	if ok, msg = submissionRepo.Delete(submissionIDL...); !ok {
		return answers, false, msg
	}
	if ok, msg = answerRepo.Delete(answerIDL...); !ok {
		return answers, false, msg
	}
	return GenerateAnswer(tx, team, usage)
}
