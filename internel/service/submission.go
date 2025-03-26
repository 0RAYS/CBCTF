package service

import (
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"gorm.io/gorm"
)

func CountFlagSolved(tx *gorm.DB, flag model.Flag) (int64, bool, string) {
	var (
		count                   int64
		submissionRepo          = db.InitSubmissionRepo(tx)
		submissions, _, ok, msg = submissionRepo.GetAllByKeyID("contest_id", flag.ContestID, -1, -1, true, 0, true)
	)
	if !ok {
		return count, false, msg
	}
	for _, submission := range submissions {
		if submission.FlagID == flag.ID {
			count++
		}
	}
	return count, true, "Success"
}
