package service

import (
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"gorm.io/gorm"
)

// IsGenerated model.Usage 需要递归预加载, depth = 3
func IsGenerated(tx *gorm.DB, usage model.Usage, team model.Team) (bool, bool, string) {
	repo := db.InitAnswerRepo(tx)
	for _, flag := range usage.Flags {
		answers, _, ok, msg := repo.GetAll(flag.ID, -1, -1, false, 0)
		if !ok {
			return false, false, msg
		}
		var count int
		for _, answer := range answers {
			if answer.TeamID == team.ID {
				count++
			}
		}
		if count < len(flag.Answers) {
			return false, false, "AnswerNotFound"
		}
	}
	return true, true, "Success"
}
