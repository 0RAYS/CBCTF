package service

import (
	"CBCTF/internel/model"
	db "CBCTF/internel/repo"
	"gorm.io/gorm"
)

func GetUsages(tx *gorm.DB, contest model.Contest, all bool) ([]model.Usage, bool, string) {
	usageRepo := db.InitUsageRepo(tx)
	usages, _, ok, msg := usageRepo.GetAll(contest.ID, -1, -1, true, 3, all)
	if !ok {
		return usages, false, msg
	}
	return usages, true, "Success"
}
